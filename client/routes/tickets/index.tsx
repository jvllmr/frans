import { Group, Table } from "@mantine/core";
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  ticketQueryOptions,
  ticketsKey,
  useDeleteTicketMutation,
} from "~/api/ticket";
import { meQueryOptions } from "~/api/user";
import { DeleteButton } from "~/components/common/DeleteButton";
import { EstimatedExpiry } from "~/components/common/EstimatedExpiry";
import { DownloadSuccessIndicator } from "~/components/file/DownloadSuccessIndicator";
import { FileRef } from "~/components/file/FileRef";
import { ShareLinkButtons } from "~/components/share/ShareLink";
import { useDateFormatter } from "~/i18n";
import { getInternalFileLink } from "~/util/link";
export const Route = createFileRoute("/tickets/")({
  component: RouteComponent,
});

function DeleteTicketButton({ ticketId }: { ticketId: string }) {
  const mutation = useDeleteTicketMutation();
  const { t } = useTranslation("ticket_active");
  return (
    <DeleteButton
      loading={mutation.isPending}
      onClick={() => {
        mutation.mutate(ticketId);
      }}
      title={t("title_delete")}
    />
  );
}

function RouteComponent() {
  const { t } = useTranslation("ticket_active");
  const { data: tickets } = useSuspenseQuery(ticketQueryOptions);
  const dateFormatter = useDateFormatter();
  const queryClient = useQueryClient();
  const { data: me } = useSuspenseQuery(meQueryOptions);
  return (
    <Table withColumnBorders withTableBorder withRowBorders>
      <Table.Thead>
        <Table.Tr>
          <Table.Th />
          <Table.Th />
          <Table.Th>{t("table_title_file")}</Table.Th>
          {me.isAdmin ? (
            <Table.Th>{t("owner", { ns: "translation" })}</Table.Th>
          ) : null}
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
                  <Group>
                    <DeleteTicketButton ticketId={ticket.id} />
                    <ShareLinkButtons shareId={ticket.id} />
                  </Group>
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
              {index === 0 && me.isAdmin ? (
                <Table.Td rowSpan={ticket.files.length}>
                  {ticket.owner.name}
                </Table.Td>
              ) : null}
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
