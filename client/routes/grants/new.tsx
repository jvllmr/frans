import { Box, Button, Flex, SimpleGrid } from "@mantine/core";
import { useForm } from "@mantine/form";
import { createFileRoute } from "@tanstack/react-router";
import { zod4Resolver } from "mantine-form-zod-resolver";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import {
  CreateGrant,
  createGrantSchemaFactory,
  useCreateGrantMutation,
} from "~/api/grant";
import { FormDebugInfo } from "~/components/dev/FormDebugInfo";
import { ExpiryParamsDownloadSection } from "~/components/form/ExpiryParamsDownloadSection";
import { ExpiryParamsUploadSection } from "~/components/form/ExpiryParamsUploadSection";
import { HisHerEmailSection } from "~/components/form/HisHerEmailSection";
import { MyEmailSection } from "~/components/form/MyEmailSection";
import { PasswordSection } from "~/components/form/PasswordSection";
import { CommentInput } from "~/components/inputs/CommentInput";
import { AvailableLanguage } from "~/i18n";

export const Route = createFileRoute("/grants/new")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t, i18n } = useTranslation("grant_new");
  const translatedCreateGrantSchema = useMemo(
    () => createGrantSchemaFactory(t),
    [t],
  );
  const form = useForm<CreateGrant>({
    initialValues: {
      comment: "",
      email: null,
      emailOnUpload: null,
      emailPassword: false,
      expiryDaysSinceLastUpload:
        window.fransGrantDefaultExpiryDaysSinceLastUpload,
      expiryTotalDays: window.fransGrantDefaultExpiryTotalDays,
      expiryTotalUploads: window.fransGrantDefaultExpiryTotalUploads,
      expiryType: "auto",
      fileExpiryDaysSinceLastDownload:
        window.fransDefaultExpiryDaysSinceLastDownload,
      fileExpiryTotalDays: window.fransDefaultExpiryTotalDays,
      fileExpiryTotalDownloads: window.fransDefaultExpiryTotalDownloads,
      fileExpiryType: "auto",
      password: "",
      receiverLang: i18n.language as AvailableLanguage,
      creatorLang: i18n.language as AvailableLanguage,
    },
    validate: zod4Resolver(translatedCreateGrantSchema),
  });
  const createGrantMutation = useCreateGrantMutation();
  return (
    <form
      onSubmit={form.onSubmit((values) => {
        createGrantMutation.mutate(values);
      })}
    >
      <Box p="lg">
        <SimpleGrid spacing="xl">
          <MyEmailSection
            // @ts-expect-error the type should match...
            form={form}
            variant="upload"
          />
          <CommentInput {...form.getInputProps("comment")} />
          <HisHerEmailSection
            // @ts-expect-error the type should match...
            form={form}
          />
          <PasswordSection
            // @ts-expect-error the type should match...
            form={form}
          />
          <ExpiryParamsUploadSection
            // @ts-expect-error the type should match...
            form={form}
            label={t("label_expiry")}
          />
          <ExpiryParamsDownloadSection
            // @ts-expect-error the type should match...
            form={form}
            variant="grant"
            label={t("label_download_expiry")}
          />
          <Flex justify="space-evenly">
            <Button
              type="submit"
              loading={createGrantMutation.isPending}
              title={t("title_create")}
            >
              {t("label_create")}
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
          <FormDebugInfo form={form} />
        </SimpleGrid>
      </Box>
    </form>
  );
}
