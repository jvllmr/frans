import {
  Box,
  Button,
  Checkbox,
  Flex,
  Highlight,
  NumberInput,
  PasswordInput,
  Select,
  SimpleGrid,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { createFileRoute } from "@tanstack/react-router";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { I18nextProvider, useTranslation } from "react-i18next";
import { CreateTicket, createTicketSchema } from "~/api/ticket";
import { FormDebugInfo } from "~/components/FormDebugInfo";
import { FilesInput } from "~/components/inputs/FilesInput";
import { NullTextarea } from "~/components/inputs/NullTextarea";
import { NullTextInput } from "~/components/inputs/NullTextInput";
import i18n from "~/i18n";
export const Route = createFileRoute("/")({
  component: Index,
});

const selectData: { label: string; value: CreateTicket["expiryType"] }[] = [
  { value: "auto", label: i18n.t("expiry_automatic", { ns: "ticket_new" }) },
  {
    value: "single",
    label: i18n.t("expiry_single_use", { ns: "ticket_new" }),
  },
  {
    value: "none",
    label: i18n.t("expiry_none", { ns: "ticket_new" }),
  },
  {
    value: "custom",
    label: i18n.t("expiry_custom", { ns: "ticket_new" }),
  },
] as const;

function NewTicketForm() {
  const { t } = useTranslation();
  const form = useForm<CreateTicket>({
    initialValues: {
      comment: null,
      email: null,
      password: "",
      emailPassword: false,
      expiryType: "auto",
      expiryTotalDays: window.fransDefaultExpiryTotalDays,
      expiryDaysSinceLastDownload:
        window.fransDefaultExpiryDaysSinceLastDownload,
      expiryTotalDownloads: window.fransDefaultExpiryTotalDownloads,
      emailOnDownload: null,
      files: [],
    },
    validate: zod4Resolver(createTicketSchema),
  });
  return (
    <Box p="lg">
      <SimpleGrid>
        <NullTextarea
          {...form.getInputProps("comment")}
          label={t("label_comment")}
        />
        <NullTextInput
          {...form.getInputProps("email")}
          label={t("label_email")}
        />
        <PasswordInput
          {...form.getInputProps("password")}
          label={t("label_password")}
        />
        <Button>{t("generate", { ns: "translation" })}</Button>
        <Checkbox
          {...form.getInputProps("emailPassword", { type: "checkbox" })}
          label={
            <Highlight highlight={t("label_password_email_highlight")}>
              {t("label_password_email")}
            </Highlight>
          }
        />
        <Select
          {...form.getInputProps("expiryType")}
          data={selectData}
          label={t("label_expiry")}
        />
        {form.values.expiryType === "custom" ? (
          <>
            <NumberInput
              {...form.getInputProps("expiryTotalDays")}
              label={t("label_expiry_total_days")}
            />
            <NumberInput
              {...form.getInputProps("expiryTotalDays")}
              label={t("label_expiry_last_download")}
            />
            <NumberInput
              {...form.getInputProps("expiryTotalDownloads")}
              label={t("label_expiry_total_downloads")}
            />
          </>
        ) : null}
        <NullTextInput
          {...form.getInputProps("emailOnDownload")}
          label={t("label_notify_email")}
        />
        <Button>{t("label_own_email")}</Button>
        <Flex justify="space-evenly">
          <Button>{t("upload", { ns: "translation" })}</Button>
          <Button>{t("reset", { ns: "translation" })}</Button>
        </Flex>
        <FilesInput {...form.getInputProps("files")} />
        <FormDebugInfo form={form} />
      </SimpleGrid>
    </Box>
  );
}

function Index() {
  return (
    <I18nextProvider i18n={i18n} defaultNS="ticket_new">
      <NewTicketForm />
    </I18nextProvider>
  );
}
