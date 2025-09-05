import { Fieldset, Group } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { useTranslation } from "react-i18next";
import { AvailableLanguage } from "~/i18n";
import { LangInput } from "../inputs/LangInput";
import { NullTextInput } from "../inputs/NullTextInput";

export interface HisHerEmailFormSectionProps {
  form: UseFormReturnType<{
    email: string | null;
    receiverLang: AvailableLanguage;
  }>;
}
export function HisHerEmailSection({ form }: HisHerEmailFormSectionProps) {
  const { t } = useTranslation("forms");
  return (
    <Fieldset>
      <Group w="100%" align="start">
        <NullTextInput
          {...form.getInputProps("email")}
          label={t("label_email")}
          w="50%"
        />
        <LangInput
          {...form.getInputProps("receiverLang")}
          label={t("label_receiver_lang")}
          required
          w="25%"
        />
      </Group>
    </Fieldset>
  );
}
