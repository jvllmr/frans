package util

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func UnpackFSToPath(fsys fs.FS, targetPath string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(".", path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetPath, relPath)

		if d.IsDir() {

			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}

func InterfaceSliceToStringSlice(in []any) []string {
	out := make([]string, len(in))
	for i, v := range in {
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[i] = s
	}
	return out
}

func GenerateRandomString(byteLength int) []byte {
	value := make([]byte, byteLength)
	_, err := rand.Read(value)
	if err != nil {
		panic(err)
	}
	return value
}

func GenerateSalt() []byte {
	return GenerateRandomString(16)
}

func HashPassword(password string, salt []byte) string {
	h1 := sha256.Sum256([]byte(password))
	combined := append(salt, h1[:]...)
	h2 := sha256.Sum256(combined)
	return hex.EncodeToString(h2[:])
}

func compareStringsTimingSafe(s1, s2 string) bool {
	return subtle.ConstantTimeCompare([]byte(s1), []byte(s2)) == 1
}

func VerifyPassword(password string, hashedPassword string, salt string) bool {
	decodedSalt, err := hex.DecodeString(salt)
	if err != nil {
		panic(err)
	}
	return compareStringsTimingSafe(HashPassword(password, decodedSalt), hashedPassword)
}

func GinAbortWithError(ctx context.Context, c *gin.Context, code int, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, "route logic expected error")
	c.AbortWithError(code, err)
	slog.ErrorContext(ctx, "route resulted in error", "err", err)
}
