import { Table, Text } from "@mantine/core";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { receivedFilesQueryOptions } from "~/api/file";
import { EstimatedExpiry } from "~/components/common/EstimatedExpiry";
import { DownloadSuccessIndicator } from "~/components/file/DownloadSuccessIndicator";
import { useDateFormatter, useFileSizeFormatter } from "~/i18n";

export const Route = createFileRoute("/grants/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: receivedFiles } = useSuspenseQuery(receivedFilesQueryOptions);
  const { t } = useTranslation("files_received");
  const dateFormatter = useDateFormatter();
  const fileSizeFormatter = useFileSizeFormatter();
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
            <Table.Th>{t("table_title_size")}</Table.Th>
            <Table.Th>{t("table_title_date")}</Table.Th>
            <Table.Th>{t("table_title_expiration")}</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {receivedFiles.map((receivedFile) => (
            <Table.Tr key={receivedFile.id}>
              <Table.Td>
                <DownloadSuccessIndicator
                  lastDownloaded={receivedFile.lastDownloaded}
                  timesDownloaded={receivedFile.timesDownloaded}
                />
              </Table.Td>
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
