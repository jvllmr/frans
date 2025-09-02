package util

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"time"
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

func estimatedExpiry(
	expiryType string,
	defaultExpiryTotalDays uint8,
	defaultExpiryDaysSinceLastDownload uint8,
	customExpiryTotalDays uint8,
	customExpiryDaysSinceLastDownload uint8,
	createdAt time.Time,
	lastDownload *time.Time,
) *time.Time {
	if expiryType == "none" {
		return nil
	}

	expiryTotalDays := defaultExpiryTotalDays
	expiryDaysSinceLastDownload := defaultExpiryDaysSinceLastDownload
	if expiryType == "custom" {
		expiryTotalDays = customExpiryTotalDays
		expiryDaysSinceLastDownload = customExpiryDaysSinceLastDownload
	}

	totalLimit := createdAt.Add(time.Hour * 24 * time.Duration(expiryTotalDays)).UTC()

	if lastDownload == nil {
		return &totalLimit
	}

	lastDownloadLimit := lastDownload.Add(
		time.Hour * 24 * time.Duration(expiryDaysSinceLastDownload),
	).UTC()

	if lastDownloadLimit.Before(totalLimit) {
		return &lastDownloadLimit
	}

	return &totalLimit
}
