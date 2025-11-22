package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/util"
)

type TicketService struct {
	cfg config.Config
	fs  FileService
}

func (ts TicketService) TicketEstimatedExpiry(ticketValue *ent.Ticket) *time.Time {
	var latestDownload *time.Time = nil
	for _, file := range ticketValue.Edges.Files {
		if latestDownload == nil ||
			(file.LastDownload != nil && latestDownload.Before(*file.LastDownload)) {
			latestDownload = file.LastDownload
		}
	}

	return estimatedExpiry(
		ticketValue.ExpiryType,
		ts.cfg.DefaultExpiryTotalDays,
		ts.cfg.DefaultExpiryDaysSinceLastDownload,
		ticketValue.ExpiryTotalDays,
		ticketValue.ExpiryDaysSinceLastDownload,
		ticketValue.CreatedAt,
		latestDownload,
	)
}

func (ts TicketService) TicketShareLink(ctx *gin.Context, ticket *ent.Ticket) string {
	return fmt.Sprintf("%s/s/%s", ts.cfg.GetBaseURL(ctx.Request), ticket.ID.String())
}

func (ts TicketService) ShouldDeleteTicket(ticketValue *ent.Ticket) bool {
	estimatedExpiry := ts.TicketEstimatedExpiry(ticketValue)
	now := time.Now()

	return len(ticketValue.Edges.Files) == 0 ||
		(estimatedExpiry != nil && estimatedExpiry.Before(now))
}

type TicketFormParams struct {
	Comment                     *string `form:"comment"`
	Email                       *string `form:"email"`
	Password                    string  `form:"password"                    binding:"required"`
	EmailPassword               bool    `form:"emailPassword"`
	ExpiryType                  string  `form:"expiryType"                  binding:"required"`
	ExpiryTotalDays             uint8   `form:"expiryTotalDays"             binding:"required"`
	ExpiryDaysSinceLastDownload uint8   `form:"expiryDaysSinceLastDownload" binding:"required"`
	ExpiryTotalDownloads        uint8   `form:"expiryTotalDownloads"        binding:"required"`
	EmailOnDownload             *string `form:"emailOnDownload"`
	CreatorLang                 string  `form:"creatorLang"                 binding:"required"`
	ReceiverLang                string  `form:"receiverLang"                binding:"required"`
}

func (ts TicketService) CreateTicket(
	ctx context.Context,
	tx *ent.Tx,
	user *ent.User,
	form *TicketFormParams,
	files []*multipart.FileHeader,
) (*ent.Ticket, error) {
	salt := util.GenerateSalt()

	hashedPassword := util.HashPassword(form.Password, salt)
	ticketBuilder := tx.Ticket.Create().
		SetID(uuid.New()).
		SetExpiryType(form.ExpiryType).
		SetExpiryDaysSinceLastDownload(form.ExpiryDaysSinceLastDownload).
		SetExpiryTotalDays(form.ExpiryTotalDays).
		SetExpiryTotalDownloads(form.ExpiryTotalDownloads).
		SetHashedPassword(hashedPassword).
		SetSalt(hex.EncodeToString(salt)).
		SetOwner(user).
		SetCreatorLang(form.CreatorLang)

	if form.Comment != nil {
		ticketBuilder = ticketBuilder.SetComment(*form.Comment)
	}

	if form.EmailOnDownload != nil {
		ticketBuilder = ticketBuilder.SetEmailOnDownload(*form.EmailOnDownload)
	}

	ticketValue, err := ticketBuilder.Save(ctx)
	if err != nil {
		return nil, err
	}

	ts.fs.EnsureFilesTmpPath()

	for _, fileHeader := range files {
		dbFile, err := ts.fs.CreateFile(
			ctx,
			tx,
			fileHeader,
			user,
			ticketValue.ExpiryType,
			ticketValue.ExpiryDaysSinceLastDownload,
			ticketValue.ExpiryTotalDays,
			ticketValue.ExpiryTotalDownloads,
		)
		if err != nil {
			return nil, err
		}

		ticketValue, err = tx.Ticket.UpdateOne(ticketValue).
			AddFiles(dbFile).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	ticketValue, err = tx.Ticket.Query().
		Where(ticket.ID(ticketValue.ID)).
		WithFiles(func(fq *ent.FileQuery) { fq.WithData().WithOwner() }).
		WithOwner().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		err = util.RefreshUserTotalDataSize(ctx, user, tx)
		if err != nil {
			slog.Error(
				"Could not refresh total data size of user",
				"err",
				err,
				"username",
				user.Username,
			)
			return nil, err
		}
	}
	err = tx.User.UpdateOne(user).AddSubmittedTickets(1).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return ticketValue, nil
}

func (ts TicketService) DeleteTicket(ctx context.Context, tx *ent.Tx, t *ent.Ticket) error {
	for _, f := range t.Edges.Files {
		if err := ts.fs.DeleteFile(ctx, f); err != nil {
			return err
		}
	}
	return tx.Ticket.DeleteOne(t).Exec(ctx)
}

type PublicTicket struct {
	ID              uuid.UUID    `json:"id"`
	Comment         *string      `json:"comment"`
	EstimatedExpiry *string      `json:"estimatedExpiry"`
	User            PublicUser   `json:"owner"`
	Files           []PublicFile `json:"files"`
	CreatedAt       string       `json:"createdAt"`
}

func (ts TicketService) ToPublicTicket(ticket *ent.Ticket) PublicTicket {
	files := make([]PublicFile, len(ticket.Edges.Files))
	for i, file := range ticket.Edges.Files {
		files[i] = ts.fs.ToPublicFile(file)
	}
	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := ts.TicketEstimatedExpiry(ticket); estimatedExpiryResult != nil {
		estimatedExpiry := estimatedExpiryResult.Format(http.TimeFormat)
		estimatedExpiryValue = &estimatedExpiry
	}

	return PublicTicket{
		ID:              ticket.ID,
		Comment:         ticket.Comment,
		User:            ToPublicUser(ticket.Edges.Owner),
		EstimatedExpiry: estimatedExpiryValue,
		Files:           files,
		CreatedAt:       ticket.CreatedAt.UTC().Format(http.TimeFormat),
	}
}

func NewTicketService(cfg config.Config, db *ent.Client) TicketService {
	return TicketService{cfg: cfg, fs: NewFileService(cfg, db)}
}
