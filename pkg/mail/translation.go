package mail

import (
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
)

//go:embed locales/**/mail.json
var localeFiles embed.FS

func getTranslationFactory(locale string) func(string) string {
	localeFileContent, err := localeFiles.ReadFile(fmt.Sprintf("locales/%s/mail.json", locale))
	if err != nil {
		slog.Error("Could not find locale", "err", err, "locale", locale)
		panic(err)
	}

	var localeJson map[string]string
	if err = json.Unmarshal(localeFileContent, &localeJson); err != nil {
		slog.Error("Could not unmarshal locale json", "err", err, "locale", locale)
		panic(err)
	}

	enFileContent, err := localeFiles.ReadFile("locales/en/mail.json")
	if err != nil {
		slog.Error("Could not find en locale", "err", err)
		panic(err)
	}

	var enJson map[string]string
	if err = json.Unmarshal(enFileContent, &enJson); err != nil {
		slog.Error("Could not unmarshal en locale json", "err", err)
		panic(err)
	}

	return func(translationKey string) string {
		if text, ok := localeJson[translationKey]; ok {
			return text
		} else if text, ok = enJson[translationKey]; ok {
			return text
		} else {
			panic(fmt.Sprintf("Could not get translation for key %s with locale %s", translationKey, locale))
		}
	}
}
