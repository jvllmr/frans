import { Button, Fieldset, Grid } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { useQuery } from "@tanstack/react-query";

import { useTranslation } from "react-i18next";
import { meQueryOptions } from "~/api/user";
import { AvailableLanguage } from "~/i18n";
import { LangInput } from "../inputs/LangInput";
import { NullTextInput } from "../inputs/NullTextInput";

interface UploadMyEmailFormSectionProps {
  form: UseFormReturnType<{
    emailOnUpload: string | null;
    creatorLang: AvailableLanguage;
  }>;
  variant: "upload";
}

interface DownloadMyEmailFormSectionProps {
  form: UseFormReturnType<{
    emailOnDownload: string | null;
    creatorLang: AvailableLanguage;
  }>;
  variant: "download";
}

export type MyEmailFormSectionProps =
  | UploadMyEmailFormSectionProps
  | DownloadMyEmailFormSectionProps;

export function MyEmailSection({ form, variant }: MyEmailFormSectionProps) {
  const { t } = useTranslation("forms");
  const { data: me } = useQuery(meQueryOptions);

  const inputPropsEmail =
    variant === "upload"
      ? form.getInputProps("emailOnUpload")
      : form.getInputProps("emailOnDownload");
  const setEmail =
    variant === "upload"
      ? (email: string) => form.setFieldValue("emailOnUpload", email)
      : (email: string) => form.setFieldValue("emailOnDownload", email);

  const inputPropsLang =
    variant === "upload"
      ? form.getInputProps("creatorLang")
      : form.getInputProps("creatorLang");

  return (
    <Fieldset>
      <Grid>
        <Grid.Col span={12}>
          <NullTextInput
            {...inputPropsEmail}
            label={
              variant === "upload"
                ? t("label_notify_email_upload")
                : t("label_notify_email_download")
            }
          />
        </Grid.Col>
        <Grid.Col span={6}>
          <Button
            fullWidth
            onClick={() => {
              if (me) {
                setEmail(me.email);
              }
            }}
            mt="lg"
            title={t("title_own_email")}
          >
            {variant === "upload"
              ? t("label_own_email_upload")
              : t("label_own_email_download")}
          </Button>
        </Grid.Col>
        <Grid.Col span={6}>
          <LangInput
            {...inputPropsLang}
            label={t("label_your_lang")}
            required
          />
        </Grid.Col>
      </Grid>
    </Fieldset>
  );
}
