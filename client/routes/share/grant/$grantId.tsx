import { Box, Button, Flex, List, Stack, Text } from "@mantine/core";
import { FileWithPath } from "@mantine/dropzone";
import { useForm } from "@mantine/form";
import { QueryKey, useQueryClient } from "@tanstack/react-query";
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
import { FileRefText } from "~/components/file/FileRef";
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

function grantShareQueryKey(grantId: string): QueryKey {
  return ["SHARE", "GRANT", grantId];
}

function GrantShare() {
  const grant = useShareGrantContext();
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
        uploadFilesToGrantMutation.mutate(values, {
          onSuccess() {
            queryClient.invalidateQueries({ queryKey: grantsKey });
            queryClient.invalidateQueries({ queryKey: filesKey });
            queryClient.invalidateQueries({
              queryKey: grantShareQueryKey(grantId),
            });
            form.reset();
          },
        });
      })}
    >
      <Stack>
        <Text>
          <b>
            {grant.owner.name} ({grant.owner.email})
          </b>{" "}
          {t("grant_message")}
        </Text>
        <Text>
          <b>{t("comment")}: </b>
          {grant.comment}
        </Text>
        {grant.files.length ? (
          <Box>
            <b>{t("files")}:</b>
            <List ml="sm">
              {grant.files.map((file) => (
                <List.Item key={file.id}>
                  <FileRefText file={file} />
                </List.Item>
              ))}
            </List>
          </Box>
        ) : null}

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
  const queryKey = useMemo(() => grantShareQueryKey(grantId), [grantId]);
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
