import {
  Box,
  Button,
  Checkbox,
  Fieldset,
  Flex,
  Grid,
  Group,
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
import { useMemo } from "react";
import { I18nextProvider, useTranslation } from "react-i18next";
import { queryClient } from "~/api";
import {
  CreateTicket,
  createTicketSchemaFactory,
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

function NewTicketForm() {
  const { t, i18n } = useTranslation();

  const expiryChoices: { label: string; value: CreateTicket["expiryType"] }[] =
    useMemo(
      () => [
        {
          value: "auto",
          label: t("expiry_automatic", { ns: "ticket_new" }),
        },
        {
          value: "single",
          label: t("expiry_single_use", { ns: "ticket_new" }),
        },
        {
          value: "none",
          label: t("expiry_none", { ns: "ticket_new" }),
        },
        {
          value: "custom",
          label: t("expiry_custom", { ns: "ticket_new" }),
        },
      ],
      [t],
    );
  const translatedCreateTicketSchema = useMemo(
    () => createTicketSchemaFactory(t),
    [t],
  );

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
    validate: zod4Resolver(translatedCreateTicketSchema),
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
        <SimpleGrid spacing="xl">
          <Box mb="sm">
            <FilesInput {...form.getInputProps("files")} />
          </Box>
          <NullTextarea
            {...form.getInputProps("comment")}
            label={t("label_comment")}
            resize="vertical"
          />
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

          <Fieldset>
            <Select
              {...form.getInputProps("expiryType")}
              data={expiryChoices}
              label={t("label_expiry")}
            />
            {form.values.expiryType === "custom" ? (
              <Group mt="xs" align="end" grow>
                <NumberInput
                  {...form.getInputProps("expiryTotalDays")}
                  label={t("label_expiry_total_days")}
                />
                <NumberInput
                  {...form.getInputProps("expiryDaysSinceLastDownload")}
                  label={t("label_expiry_last_download")}
                />
                <NumberInput
                  {...form.getInputProps("expiryTotalDownloads")}
                  label={t("label_expiry_total_downloads")}
                />
              </Group>
            ) : null}
          </Fieldset>
          <Fieldset>
            <Grid>
              <Grid.Col span={12}>
                <NullTextInput
                  {...form.getInputProps("emailOnDownload")}
                  label={t("label_notify_email")}
                />
              </Grid.Col>
              <Grid.Col span={6}>
                <Button
                  fullWidth
                  onClick={() => {
                    if (me) {
                      form.setFieldValue("emailOnDownload", me.email);
                    }
                  }}
                  mt="lg"
                >
                  {t("label_own_email")}
                </Button>
              </Grid.Col>
              <Grid.Col span={6}>
                <LangInput
                  {...form.getInputProps("creatorLang")}
                  label={t("label_your_lang")}
                  required
                />
              </Grid.Col>
            </Grid>
          </Fieldset>
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
