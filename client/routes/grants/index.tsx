import { Table, Text } from "@mantine/core";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { queryClient } from "~/api";
import {
  filesKey,
  receivedFilesQueryOptions,
  useDeleteFileMutation,
} from "~/api/file";
import { meQueryOptions } from "~/api/user";
import { DeleteButton } from "~/components/common/DeleteButton";
import { EstimatedExpiry } from "~/components/common/EstimatedExpiry";
import { DownloadSuccessIndicator } from "~/components/file/DownloadSuccessIndicator";
import { FileRef } from "~/components/file/FileRef";
import { useDateFormatter, useFileSizeFormatter } from "~/i18n";
import { getInternalFileLink } from "~/util/link";

export const Route = createFileRoute("/grants/")({
  component: RouteComponent,
});

function DeleteFileButton({ fileId }: { fileId: string }) {
  const mutation = useDeleteFileMutation();
  const { t } = useTranslation("files_received");
  return (
    <DeleteButton
      loading={mutation.isPending}
      onClick={() => {
        mutation.mutate(fileId);
      }}
      title={t("title_delete")}
    />
  );
}

function RouteComponent() {
  const { data: receivedFiles } = useSuspenseQuery(receivedFilesQueryOptions);
  const { t } = useTranslation("files_received");
  const dateFormatter = useDateFormatter();
  const fileSizeFormatter = useFileSizeFormatter();
  const { data: me } = useSuspenseQuery(meQueryOptions);
  const totalSize = useMemo(
    () =>
      receivedFiles.reduce((prev, receivedFile) => prev + receivedFile.size, 0),
    [receivedFiles],
  );
  return (
    <>
      <Table withColumnBorders withTableBorder withRowBorders>
        <Table.Thead>
          <Table.Tr>
            <Table.Th />
            <Table.Th>{t("table_title_file")}</Table.Th>
            {me.isAdmin ? (
              <Table.Th>{t("owner", { ns: "translation" })}</Table.Th>
            ) : null}
            <Table.Th>{t("table_title_size")}</Table.Th>
            <Table.Th>{t("table_title_date")}</Table.Th>
            <Table.Th>{t("table_title_expiration")}</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {receivedFiles.map((receivedFile) => (
            <Table.Tr key={receivedFile.id}>
              <Table.Td>
                <DeleteFileButton fileId={receivedFile.id} />
                <DownloadSuccessIndicator
                  lastDownloaded={receivedFile.lastDownloaded}
                  timesDownloaded={receivedFile.timesDownloaded}
                />
              </Table.Td>
              <Table.Td>
                <FileRef
                  file={receivedFile}
                  link={getInternalFileLink(receivedFile.id, true)}
                  withoutSize
                  onClick={() => {
                    queryClient.invalidateQueries({ queryKey: filesKey });
                  }}
                />
              </Table.Td>
              {me.isAdmin ? (
                <Table.Td>{receivedFile.owner.name}</Table.Td>
              ) : null}
              <Table.Td>{fileSizeFormatter(receivedFile.size)}</Table.Td>

              <Table.Td>
                {dateFormatter.format(receivedFile.createdAt)}
              </Table.Td>
              <Table.Td>
                <EstimatedExpiry
                  estimatedExpiry={receivedFile.estimatedExpiry}
                />
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
      <Text pt="xl">
        {t("total_size")}: <b>{fileSizeFormatter(totalSize)}</b>
      </Text>
    </>
  );
}
