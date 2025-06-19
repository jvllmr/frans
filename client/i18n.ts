import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import LanguageDetector from "i18next-browser-languagedetector";

import resources from "virtual:i18next-loader";

export const availableLanguages = [
  "pt-BR",
  "cz",
  "de",
  "en",
  "es",
  "fr",
  "it",
  "ja",
  "nl",
  "ru",
  "zh",
] as const;

export type AvailableLanguage = (typeof availableLanguages)[number];

export const availableLanguagesLabels: Record<AvailableLanguage, string> = {
  "pt-BR": "BR",
  cz: "CZ",
  de: "DE",
  en: "EN",
  es: "ES",
  fr: "FR",
  it: "IT",
  ja: "JA",
  nl: "NL",
  ru: "RU",
  zh: "ZH",
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: "en",
    resources,
    detection: {
      convertDetectedLanguage(lng) {
        if (!lng.startsWith("pt") && lng.includes("-")) {
          return lng.substring(0, 2);
        }
        return lng;
      },
    },
    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },
  });

export default i18n;
