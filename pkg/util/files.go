package util

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
)

func EnsureFilesTmpPath(configValue config.Config) {
	os.MkdirAll(GetFilesTmpPath(configValue), 0664)
}

func GetFilesTmpPath(configValue config.Config) string {
	return fmt.Sprintf("%s/%s", configValue.FilesDir, "tmp")
}

func GetFilesTmpFilePath(configValue config.Config) string {
	return fmt.Sprintf("%s/%s", GetFilesTmpPath(configValue), uuid.New())
}

func GetFilesFilePath(configValue config.Config, fileName string) string {
	return fmt.Sprintf("%s/%s", configValue.FilesDir, fileName)
}
