package util

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
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
