package services

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/file"
	"github.com/jvllmr/frans/internal/ent/filedata"
	"github.com/jvllmr/frans/internal/otel"
)

type ErrFileTooBig struct {
	size    int64
	maxSize int64
}

func (e *ErrFileTooBig) Error() string {
	return fmt.Sprintf(
		"file too big: tried to create file with %d bytes in size, but only %d bytes are allowed",
		e.size,
		e.maxSize,
	)
}

var _ error = (*ErrFileTooBig)(nil)

type FileService struct {
	config config.Config
	db     *ent.Client
}

func (fs FileService) EnsureFilesTmpPath() {
	err := os.MkdirAll(fs.FilesTmpPath(), 0775)
	if err != nil {
		panic(err)
	}
}

func (fs FileService) FilesTmpPath() string {
	return fmt.Sprintf("%s/%s", fs.config.FilesDir, "tmp")
}

func (fs FileService) FilesTmpFilePath() string {
	return fmt.Sprintf("%s/%s", fs.FilesTmpPath(), uuid.New())
}

func (fs FileService) FilesFilePath(fileName string) string {
	return fmt.Sprintf("%s/%s", fs.config.FilesDir, fileName)
}

func (fs FileService) ShouldDeleteFile(
	fileValue *ent.File,
) bool {
	if fileValue.ExpiryType == config.TicketExpiryTypeNone {
		return false
	}
	if fileValue.ExpiryType == config.TicketExpiryTypeSingle {
		return fileValue.TimesDownloaded > 0
	}
	estimatedExpiry := *fs.FileEstimatedExpiry(fileValue)
	now := time.Now()

	if fileValue.ExpiryType == config.TicketExpiryTypeCustom {
		return fileValue.TimesDownloaded >= uint64(fileValue.ExpiryTotalDownloads) ||
			estimatedExpiry.Before(now)
	}
	return fileValue.TimesDownloaded >= uint64(fs.config.DefaultExpiryTotalDownloads) ||
		estimatedExpiry.Before(now)
}

func (fs FileService) CreateFile(
	ctx context.Context,
	tx *ent.Tx,
	fileHeader *multipart.FileHeader,
	user *ent.User,
	expiryType string,
	expiryDaysSinceLastDownload uint8,
	expiryTotalDays uint8,
	expiryTotalDownloads uint8,
) (*ent.File, error) {
	ctx, span := otel.NewSpan(ctx, "createFile")
	defer span.End()
	if fileHeader.Size > fs.config.MaxSizes {
		return nil, &ErrFileTooBig{
			size:    fileHeader.Size,
			maxSize: fs.config.MaxSizes,
		}
	}

	incomingFileHandle, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	hasher := sha512.New()
	tmpFilePath := fs.FilesTmpFilePath()
	tmpFileHandle, err := os.Create(tmpFilePath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if err := tmpFileHandle.Close(); err != nil {
			slog.Warn("could not close temporary file", "path", tmpFilePath, "err", err)
		}
		if err := os.Remove(tmpFilePath); err != nil {
			slog.Warn("could not remove temporary file", "path", tmpFilePath, "err", err)
		}
	}()
	writer := io.MultiWriter(hasher, tmpFileHandle)
	_, err = io.Copy(writer, incomingFileHandle)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}

	hash := hasher.Sum(nil)
	sha512sum := hex.EncodeToString(hash)

	fileData, err := tx.FileData.Get(ctx, sha512sum)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}
		fileData, err = tx.FileData.Create().
			SetID(sha512sum).
			SetSize(uint64(fileHeader.Size)).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	dbFile, err := tx.File.Create().
		SetID(uuid.New()).
		SetName(fileHeader.Filename).
		SetExpiryType(expiryType).
		SetExpiryDaysSinceLastDownload(expiryDaysSinceLastDownload).
		SetExpiryTotalDays(expiryTotalDays).
		SetExpiryTotalDownloads(expiryTotalDownloads).
		SetData(fileData).
		SetOwner(user).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	targetFilePath := fs.FilesFilePath(sha512sum)
	if _, err = os.Stat(targetFilePath); err != nil {
		if err = os.Rename(tmpFilePath, targetFilePath); err != nil {
			return nil, fmt.Errorf("create file: %w", err)
		}
	}

	return dbFile, nil
}

func (fs FileService) DeleteFile(ctx context.Context, fileValue *ent.File) error {
	ctx, span := otel.NewSpan(ctx, "DeleteFile")
	defer span.End()
	fileDataFilesCount, err := fs.db.FileData.Query().
		Where(filedata.HasFilesWith(file.ID(fileValue.ID))).
		QueryFiles().
		Count(ctx)
	if err != nil {
		return err
	}
	err = fs.db.File.DeleteOne(fileValue).Exec(ctx)
	if err != nil {
		return err
	}
	deleteFromFS := fileDataFilesCount <= 1
	if deleteFromFS {
		filePath := fs.FilesFilePath(fileValue.Edges.Data.ID)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
		err = fs.db.FileData.DeleteOne(fileValue.Edges.Data).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return err
}

func (fs FileService) FileEstimatedExpiry(fileValue *ent.File) *time.Time {
	return estimatedExpiry(
		fileValue.ExpiryType,
		fs.config.DefaultExpiryTotalDays,
		fs.config.DefaultExpiryDaysSinceLastDownload,
		fileValue.ExpiryTotalDays,
		fileValue.ExpiryDaysSinceLastDownload,
		fileValue.CreatedAt,
		fileValue.LastDownload,
	)
}

type PublicFile struct {
	Id              uuid.UUID  `json:"id"`
	Sha512          string     `json:"sha512"`
	Size            uint64     `json:"size"`
	Name            string     `json:"name"`
	CreatedAt       string     `json:"createdAt"`
	TimesDownloaded uint64     `json:"timesDownloaded"`
	LastDownloaded  *string    `json:"lastDownloaded"`
	EstimatedExpiry *string    `json:"estimatedExpiry"`
	User            PublicUser `json:"owner"`
}

func (fs FileService) ToPublicFile(file *ent.File) PublicFile {
	var lastDownloadedValue *string = nil
	if file.LastDownload != nil {
		formattedValue := file.LastDownload.UTC().Format(http.TimeFormat)
		lastDownloadedValue = &formattedValue
	}

	var estimatedExpiryValue *string = nil

	if estimatedExpiryResult := fs.FileEstimatedExpiry(file); estimatedExpiryResult != nil {
		estimatedExpiry := estimatedExpiryResult.Format(http.TimeFormat)
		estimatedExpiryValue = &estimatedExpiry
	}

	return PublicFile{
		Id:              file.ID,
		Sha512:          file.Edges.Data.ID,
		Size:            file.Edges.Data.Size,
		Name:            file.Name,
		CreatedAt:       file.CreatedAt.UTC().Format(http.TimeFormat),
		TimesDownloaded: file.TimesDownloaded,
		LastDownloaded:  lastDownloadedValue,
		EstimatedExpiry: estimatedExpiryValue,
		User:            ToPublicUser(file.Edges.Owner),
	}

}

func NewFileService(c config.Config, db *ent.Client) FileService {
	fs := FileService{config: c, db: db}
	fs.EnsureFilesTmpPath()
	return fs
}
