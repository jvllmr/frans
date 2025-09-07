import {
  Anchor,
  Box,
  Button,
  Center,
  Flex,
  Group,
  SimpleGrid,
  Stack,
  Text,
} from "@mantine/core";
import { useForm } from "@mantine/form";

import { createFileRoute } from "@tanstack/react-router";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useMemo, useState } from "react";
import { I18nextProvider, useTranslation } from "react-i18next";
import { queryClient } from "~/api";
import {
  CreateTicket,
  createTicketSchemaFactory,
  ticketsKey,
  useCreateTicketMutation,
} from "~/api/ticket";
import { FormDebugInfo } from "~/components/dev/FormDebugInfo";
import { ExpiryParamsDownloadSection } from "~/components/form/ExpiryParamsDownloadSection";
import { HisHerEmailSection } from "~/components/form/HisHerEmailSection";
import { MyEmailSection } from "~/components/form/MyEmailSection";
import { PasswordSection } from "~/components/form/PasswordSection";
import { ProgressBar } from "~/components/form/ProgressBar";
import { CommentInput } from "~/components/inputs/CommentInput";
import { FilesInput } from "~/components/inputs/FilesInput";
import { CopyShareLinkButton } from "~/components/share/ShareLink";
import i18n, { AvailableLanguage } from "~/i18n";
import { getShareLink } from "~/util/link";
import { useProgressHandle } from "~/util/progress";
export const Route = createFileRoute("/")({
  component: Index,
});

function NewTicketForm() {
  const { t, i18n } = useTranslation();

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
  const [ticketId, setTicketId] = useState<string | null>(null);
  const shareLink = getShareLink(ticketId ?? "");

  return ticketId ? (
    <Stack>
      <ProgressBar state={progressHandle.state} />
      <Center>
        <Text>{t("ticket_available")}</Text>
      </Center>
      <Center>
        <Box>
          <Group gap={2}>
            <b>
              {t("url", { ns: "translation" })}:{" "}
              <Anchor href={shareLink}>{shareLink}</Anchor>
            </b>
            <CopyShareLinkButton
              variant="subtle"
              color="gray"
              shareId={ticketId}
            />
          </Group>
          <Text>
            <b>{t("password", { ns: "translation" })}:</b>{" "}
            {form.values.password}
          </Text>
        </Box>
      </Center>

      <Button
        fullWidth
        onClick={() => {
          form.reset();
          setTicketId(null);
          progressHandle.reset();
        }}
      >
        {t("ticket_another")}
      </Button>
    </Stack>
  ) : (
    <form
      onSubmit={form.onSubmit((values) => {
        createTicketMutation.mutate(values, {
          onSuccess(data) {
            setTicketId(data.id);
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
          <CommentInput {...form.getInputProps("comment")} />
          <HisHerEmailSection
            // @ts-expect-error the type should match...
            form={form}
          />
          <PasswordSection
            // @ts-expect-error the type should match...
            form={form}
          />

          <ExpiryParamsDownloadSection
            // @ts-expect-error the type should match...
            form={form}
            variant="ticket"
            label={t("label_expiry")}
          />
          <MyEmailSection
            // @ts-expect-error the type should match...
            form={form}
            variant="download"
          />
          <Flex justify="space-evenly">
            <Button
              type="submit"
              loading={createTicketMutation.isPending}
              title={t("title_upload")}
            >
              {t("upload", { ns: "translation" })}
            </Button>
            <Button
              onClick={() => {
                form.reset();
              }}
              title={t("title_reset")}
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
