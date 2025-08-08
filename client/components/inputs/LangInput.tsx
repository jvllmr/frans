import { ComboboxItem, Select, SelectProps } from "@mantine/core";
import { useTranslation } from "react-i18next";
import {
  AvailableLanguage,
  availableLanguages,
  availableLanguagesLabels,
} from "~/i18n";

export interface LangInputProps extends Omit<SelectProps, "data" | "value"> {
  value?: AvailableLanguage;
}

const langInputData: ComboboxItem[] = availableLanguages.map((lang) => ({
  value: lang,
  label: availableLanguagesLabels[lang],
}));

export function LangInput(props: LangInputProps) {
  const { t } = useTranslation("lang_input");
  return (
    <Select placeholder={t("placeholder")} {...props} data={langInputData} />
  );
}
