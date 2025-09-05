import {
  Button,
  Checkbox,
  Fieldset,
  Group,
  Highlight,
  PasswordInput,
} from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import passwordGenerator from "generate-password-browser";
import { useTranslation } from "react-i18next";

export interface PasswordSectionProps<
  TForm extends UseFormReturnType<{ password: string }>,
> {
  form: TForm;
}

export function PasswordSection<
  TForm extends UseFormReturnType<{ password: string }>,
>({ form }: PasswordSectionProps<TForm>) {
  const { t } = useTranslation("forms");

  return (
    <Fieldset>
      <Group w="100%" mb="xs" align="start">
        <PasswordInput
          {...form.getInputProps("password")}
          label={t("label_password")}
          withAsterisk
          required
          w="50%"
        />
        <Button
          mt="lg"
          title={t("title_generate_password")}
          onClick={() => {
            form.setFieldValue(
              "password",
              passwordGenerator.generate({
                length: 12,
                strict: true,
                numbers: true,
                excludeSimilarCharacters: true,
              }),
            );
          }}
        >
          {t("generate", { ns: "translation" })}
        </Button>
      </Group>
      <Checkbox
        {...form.getInputProps("emailPassword", { type: "checkbox" })}
        label={
          <Highlight highlight={t("label_password_email_highlight")}>
            {t("label_password_email")}
          </Highlight>
        }
      />
    </Fieldset>
  );
}
