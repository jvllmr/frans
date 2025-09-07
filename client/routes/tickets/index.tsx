import { Table } from "@mantine/core";
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { ticketQueryOptions, ticketsKey } from "~/api/ticket";
import { EstimatedExpiry } from "~/components/common/EstimatedExpiry";
import { DownloadSuccessIndicator } from "~/components/file/DownloadSuccessIndicator";
import { FileRef } from "~/components/file/FileRef";
import { ShareLinkButtons } from "~/components/share/ShareLink";
import { useDateFormatter } from "~/i18n";
import { getInternalFileLink } from "~/util/link";
export const Route = createFileRoute("/tickets/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation("ticket_active");
  const { data: tickets } = useSuspenseQuery(ticketQueryOptions);
  const dateFormatter = useDateFormatter();
  const queryClient = useQueryClient();
  return (
    <Table withColumnBorders withTableBorder withRowBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th />
          <Table.Th />
          <Table.Th>{t("table_title_file")}</Table.Th>
          <Table.Th>{t("table_title_created_at")}</Table.Th>
          <Table.Th>{t("table_title_expiration")}</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {tickets.flatMap((ticket) =>
          ticket.files.map((file, index) => (
            <Table.Tr key={file.id}>
              {index === 0 ? (
                <Table.Td rowSpan={ticket.files.length}>
                  <ShareLinkButtons shareId={ticket.id} />
                </Table.Td>
              ) : null}

              <Table.Td>
                <DownloadSuccessIndicator
                  lastDownloaded={file.lastDownloaded}
                  timesDownloaded={file.timesDownloaded}
                />
              </Table.Td>
              <Table.Td>
                <FileRef
                  file={file}
                  link={getInternalFileLink(file.id)}
                  withoutSize
                  onClick={() => {
                    queryClient.invalidateQueries({ queryKey: ticketsKey });
                  }}
                />
              </Table.Td>
              {index === 0 ? (
                <Table.Td rowSpan={ticket.files.length}>
                  {dateFormatter.format(ticket.createdAt)}
                </Table.Td>
              ) : null}
              {index === 0 ? (
                <Table.Td rowSpan={ticket.files.length}>
                  <EstimatedExpiry estimatedExpiry={ticket.estimatedExpiry} />
                </Table.Td>
              ) : null}
            </Table.Tr>
          )),
        )}
      </Table.Tbody>
    </Table>
  );
}
