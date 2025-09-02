import { Button, Flex, Stack, Text } from "@mantine/core";
import { FileWithPath } from "@mantine/dropzone";
import { useForm } from "@mantine/form";
import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import React, { useCallback, useContext, useMemo } from "react";
import { useTranslation } from "react-i18next";
import z from "zod";
import { filesKey } from "~/api/file";

import {
  fetchGrantShare,
  fetchGrantShareAccessToken,
  Grant,
  grantsKey,
  useGrantUploadMutation,
} from "~/api/grant";
import { ProgressBar } from "~/components/form/ProgressBar";

import { FilesInput } from "~/components/inputs/FilesInput";
import { ShareAuth } from "~/components/share/ShareAuth";
import { useShareAuthContext } from "~/components/share/shareAuthContext";
import { useProgressHandle } from "~/util/progress";

export const Route = createFileRoute("/share/grant/$grantId")({
  component: RouteComponent,
  params: { parse: z.object({ grantId: z.uuid() }).parse },
});

const shareGrantContext = React.createContext<Grant | null>(null);

function useShareGrantContext() {
  const grant = useContext(shareGrantContext);
  if (!grant) throw TypeError("Expected grant to be available");
  return grant;
}

function GrantShare() {
  const ticket = useShareGrantContext();
  const { t } = useTranslation("share");
  const queryClient = useQueryClient();
  const form = useForm<{ files: FileWithPath[] }>({
    initialValues: { files: [] },
  });
  const grantId = Route.useParams({ select: (p) => p.grantId });
  const progressHandle = useProgressHandle();
  const { password } = useShareAuthContext();
  const uploadFilesToGrantMutation = useGrantUploadMutation(
    {
      grantId,
      password,
    },
    progressHandle,
  );
  return (
    <form
      onSubmit={form.onSubmit((values) => {
        uploadFilesToGrantMutation.mutate(values);
        queryClient.invalidateQueries({ queryKey: grantsKey });
        queryClient.invalidateQueries({ queryKey: filesKey });
      })}
    >
      <Stack>
        <Text>
          <b>
            {ticket.owner.name} ({ticket.owner.email})
          </b>{" "}
          {t("grant_message")}
        </Text>
        <Text>
          <b>{t("comment")}: </b>
          {ticket.comment}
        </Text>
        <FilesInput {...form.getInputProps("files")} />
        <Flex justify="space-evenly">
          <Button type="submit">{t("upload", { ns: "translation" })}</Button>
          <Button
            onClick={() => {
              form.reset();
            }}
          >
            {t("reset", { ns: "translation" })}
          </Button>
        </Flex>
        <ProgressBar state={progressHandle.state} />
      </Stack>
    </form>
  );
}

function RouteComponent() {
  const grantId = Route.useParams({ select: (p) => p.grantId });
  const dataFetcher = useCallback(
    (password: string) => fetchGrantShare({ grantId: grantId, password }),
    [grantId],
  );
  const queryKey = useMemo(() => ["SHARE", "GRANT", grantId], [grantId]);
  const shareTokenGenerator = useCallback(
    (password: string) =>
      fetchGrantShareAccessToken({ grantId: grantId, password }),
    [grantId],
  );
  const { t } = useTranslation("share");
  return (
    <ShareAuth
      DataContextProvider={shareGrantContext.Provider}
      dataFetcher={dataFetcher}
      dataQueryKey={queryKey}
      shareTokenGenerator={shareTokenGenerator}
      prompt={t("grant_prompt")}
      submitButtonLabel={t("grant_submit")}
    >
      <GrantShare />
    </ShareAuth>
  );
}
