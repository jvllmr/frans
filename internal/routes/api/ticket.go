package apiRoutes

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/ent/user"
	"github.com/jvllmr/frans/internal/mail"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/otel"
	apiTypes "github.com/jvllmr/frans/internal/routes/api/types"
	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/util"
)

type ticketController struct {
	config        config.Config
	db            *ent.Client
	ticketService services.TicketService
	mailer        mail.Mailer
}

func (tc *ticketController) createTicketHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "createTicket")
	defer span.End()

	currentUser := middleware.GetCurrentUser(c)
	var form services.TicketFormParams
	tx, err := tc.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
		return
	}
	if err := c.ShouldBind(&form); err == nil {
		multipartForm, _ := c.MultipartForm()
		files := multipartForm.File["files[]"]
		if len(files) > int(tc.config.MaxFiles) {
			util.GinAbortWithError(
				ctx,
				c,
				http.StatusBadRequest,
				fmt.Errorf(
					"maximum of %d files allowed per upload. %d uploaded",
					tc.config.MaxFiles,
					len(files),
				),
			)
			return
		}
		ticketValue, err := tc.ticketService.CreateTicket(ctx, tx, currentUser, &form, files)
		if err != nil {
			var errFileTooBig *services.ErrFileTooBig
			if errors.As(err, &errFileTooBig) {
				util.GinAbortWithError(ctx, c, http.StatusBadRequest, err)
			} else {
				util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
			}
		}
		c.JSON(http.StatusCreated, tc.ticketService.ToPublicTicket(ticketValue))
		if form.Email != nil {
			var toBeEmailedPassword *string = nil
			if form.EmailPassword {
				toBeEmailedPassword = &form.Password
			}
			if err := tc.mailer.SendTicketSharedNotification(
				c,
				tc.ticketService,
				*form.Email,
				form.ReceiverLang,
				ticketValue,
				toBeEmailedPassword,
			); err != nil {
				util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
			return
		}
	} else {
		util.GinAbortWithError(ctx, c, http.StatusUnprocessableEntity, err)
	}

}

func (tc *ticketController) fetchTicketsHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "fetchTickets")
	defer span.End()
	currentUser := middleware.GetCurrentUser(c)
	query := tc.db.Ticket.Query().
		WithFiles(func(fq *ent.FileQuery) { fq.WithData().WithOwner() }).
		WithOwner()

	if !currentUser.IsAdmin {
		query = query.Where(ticket.HasOwnerWith(user.ID(currentUser.ID)))
	}

	tickets, err := query.All(ctx)
	if err != nil {
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
	}
	publicTickets := make([]services.PublicTicket, len(tickets))
	for i, ticketValue := range tickets {
		publicTickets[i] = tc.ticketService.ToPublicTicket(ticketValue)
	}
	c.JSON(http.StatusOK, publicTickets)
}

func (tc *ticketController) deleteTicketHandler(c *gin.Context) {
	ctx, span := otel.NewSpan(c.Request.Context(), "deleteTicketManual")
	defer span.End()
	var requestedTicket apiTypes.RequestedTicketParam
	if err := c.ShouldBindUri(&requestedTicket); err != nil {
		util.GinAbortWithError(ctx, c, http.StatusBadRequest, err)
		return
	}
	t, err := tc.db.Ticket.Query().
		Where(ticket.ID(uuid.MustParse(requestedTicket.ID))).
		WithOwner().
		WithFiles(func(fq *ent.FileQuery) { fq.WithData() }).
		Only(ctx)
	if err != nil {
		util.GinAbortWithError(ctx, c, http.StatusNotFound, err)
		return
	}
	currentUser := middleware.GetCurrentUser(c)
	isUserOwner := t.Edges.Owner.ID == currentUser.ID
	if !currentUser.IsAdmin && !isUserOwner {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	tx, err := tc.db.Tx(ctx)
	if err != nil {
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
		return
	}

	if err := tc.ticketService.DeleteTicket(ctx, tx, t); err != nil {
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
		return
	}

	if !isUserOwner {
		if err := tc.mailer.SendTicketDeletionNotification(t, tc.config.GetBaseURL(c.Request)); err != nil {
			util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		util.GinAbortWithError(ctx, c, http.StatusInternalServerError, err)
		return
	}
	slog.InfoContext(
		ctx,
		"Manual ticket deletion",
		"username",
		currentUser.Username,
		"owner",
		t.Edges.Owner.Username,
		"ticketId",
		t.ID.String(),
	)
	c.Status(http.StatusOK)
}

func setupTicketGroup(r *gin.RouterGroup, configValue config.Config, db *ent.Client) {
	controller := ticketController{
		config:        configValue,
		db:            db,
		ticketService: services.NewTicketService(configValue, db),
		mailer:        mail.NewMailer(configValue),
	}
	r.POST("", controller.createTicketHandler)
	r.GET("", controller.fetchTicketsHandler)
	r.DELETE("/:ticketId", controller.deleteTicketHandler)
}
