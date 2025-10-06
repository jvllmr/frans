package testutil

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/services"
)

func SetupTestFileService(cfg config.Config, db *ent.Client) services.FileService {
	return services.NewFileService(cfg, db)
}

func createFileHeader(filename string, content string) *multipart.FileHeader {
	resultBytes := &bytes.Buffer{}

	writer := multipart.NewWriter(resultBytes)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		panic(err)
	}
	stringReader := strings.NewReader(content)
	_, err = io.Copy(part, stringReader)
	if err != nil {
		panic(err)
	}
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	reader := multipart.NewReader(resultBytes, writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	if err != nil {
		panic(err)
	}
	fileHeaders := form.File["file"]
	if len(fileHeaders) != 1 {
		panic("Expected exactly one file header")
	}
	return fileHeaders[0]
}

func SetupTestFile(
	t *testing.T,
	cfg config.Config,
	db *ent.Client,
	filename string,
	content string,
	user *ent.User,
	expiryType string,
	expirySinceLastDownload uint8,
	expiryTotalDays uint8,
	expiryTotalDownloads uint8,
) *ent.File {
	fs := SetupTestFileService(cfg, db)
	tx, err := db.Tx(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	fileHeader := createFileHeader(filename, content)
	dbFile, err := fs.CreateFile(
		t.Context(),
		tx,
		fileHeader,
		user,
		expiryType,
		expirySinceLastDownload,
		expiryTotalDays,
		expiryTotalDownloads,
	)
	if err != nil {
		t.Fatal(err)
	}
	tx.Commit()
	return dbFile
}
