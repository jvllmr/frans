package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func InterfaceSliceToStringSlice(in []interface{}) []string {
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

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string, salt []byte) string {
	h1 := sha256.Sum256([]byte(password))
	combined := append(salt, h1[:]...)
	h2 := sha256.Sum256(combined)
	return hex.EncodeToString(h2[:])
}

func GetFileHash(file multipart.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil

}
