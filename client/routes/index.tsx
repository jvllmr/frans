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
import passwordGenerator from "generate-password-browser";

import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { I18nextProvider, useTranslation } from "react-i18next";
import { queryClient } from "~/api";
import {
  CreateTicket,
  createTicketSchema,
  ticketsKey,
  useCreateTicketMutation,
} from "~/api/ticket";
import { meQueryOptions } from "~/api/user";
import { FormDebugInfo } from "~/components/FormDebugInfo";
import { FilesInput } from "~/components/inputs/FilesInput";
import { LangInput } from "~/components/inputs/LangInput";
import { NullTextarea } from "~/components/inputs/NullTextarea";
import { NullTextInput } from "~/components/inputs/NullTextInput";
import { ProgressBar } from "~/components/ProgressBar";
import i18n, { AvailableLanguage } from "~/i18n";
import { useProgressHandle } from "~/util/progress";
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
  const { t, i18n } = useTranslation();

  const form = useForm<CreateTicket>({
    initialValues: {
      comment: null,
      email: null,
      receiverLang: i18n.language as AvailableLanguage,
      password: "",
      emailPassword: false,
      expiryType: "auto",
      expiryTotalDays: window.fransDefaultExpiryTotalDays,
      expiryDaysSinceLastDownload:
        window.fransDefaultExpiryDaysSinceLastDownload,
      expiryTotalDownloads: window.fransDefaultExpiryTotalDownloads,
      emailOnDownload: null,
      creatorLang: i18n.language as AvailableLanguage,
      files: [],
    },
    validate: zod4Resolver(createTicketSchema),
  });
  const progressHandle = useProgressHandle();
  const createTicketMutation = useCreateTicketMutation(progressHandle);
  const { data: me } = useQuery(meQueryOptions);

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        createTicketMutation.mutate(values, {
          onSuccess() {
            form.reset();
            queryClient.invalidateQueries({ queryKey: ticketsKey });
          },
        });
      })}
    >
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
          <LangInput {...form.getInputProps("receiverLang")} />
          <PasswordInput
            {...form.getInputProps("password")}
            label={t("label_password")}
            withAsterisk
            required
          />
          <Button
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
          <Button
            onClick={() => {
              if (me) {
                form.setFieldValue("emailOnDownload", me.email);
              }
            }}
          >
            {t("label_own_email")}
          </Button>
          <LangInput {...form.getInputProps("receiverLang")} />
          <Flex justify="space-evenly">
            <Button type="submit" loading={createTicketMutation.isPending}>
              {t("upload", { ns: "translation" })}
            </Button>
            <Button
              onClick={() => {
                form.reset();
              }}
            >
              {t("reset", { ns: "translation" })}
            </Button>
          </Flex>
          <ProgressBar state={progressHandle.state} />
          <FilesInput {...form.getInputProps("files")} />
          <FormDebugInfo form={form} />
        </SimpleGrid>
      </Box>
    </form>
  );
}

function Index() {
  return (
    <I18nextProvider i18n={i18n} defaultNS="ticket_new">
      <NewTicketForm />
    </I18nextProvider>
  );
}
