import { Fieldset, Group } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { useTranslation } from "react-i18next";
import { AvailableLanguage } from "~/i18n";
import { LangInput } from "../inputs/LangInput";
import { NullTagsInput } from "../inputs/NullTagsInput";

interface HisHerEmailFormValues {
  email: string[] | null;
  receiverLang: AvailableLanguage;
}

export interface HisHerEmailFormSectionProps<
  TForm extends UseFormReturnType<HisHerEmailFormValues>,
> {
  form: TForm;
}
export function HisHerEmailSection<
  TForm extends UseFormReturnType<HisHerEmailFormValues>,
>({ form }: HisHerEmailFormSectionProps<TForm>) {
  const { t } = useTranslation("forms");
  return (
    <Fieldset>
      <Group w="100%" align="start">
        <NullTagsInput
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
