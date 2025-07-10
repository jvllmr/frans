import { ActionIcon, Anchor, Table, Text } from "@mantine/core";
import { IconCheck } from "@tabler/icons-react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { differenceInDays } from "date-fns";
import { useTranslation } from "react-i18next";
import { ticketQueryOptions } from "~/api/ticket";
import { useDateFormatter, useRelativeDateFormatter } from "~/i18n";
export const Route = createFileRoute("/tickets/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation("ticket_active");
  const { data: tickets } = useSuspenseQuery(ticketQueryOptions);
  const dateFormatter = useDateFormatter();
  const relativeDateFormatter = useRelativeDateFormatter();
  const now = new Date();
  return (
    <Table withColumnBorders>
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
                  <ActionIcon>
                    {
                      // TODO: Link somewhere
                    }
                  </ActionIcon>
                </Table.Td>
              ) : null}

              <Table.Td>
                {file.timesDownloaded > 0 ? (
                  <Text span c="teal">
                    <IconCheck />
                  </Text>
                ) : null}
              </Table.Td>
              <Table.Td>
                <Anchor href={`${window.fransRootPath}/api/v1/file/${file.id}`}>
                  {file.name}
                </Anchor>
              </Table.Td>
              {index === 0 ? (
                <Table.Td rowSpan={ticket.files.length}>
                  {dateFormatter.format(ticket.createdAt)}
                </Table.Td>
              ) : null}
              {index === 0 ? (
                <Table.Td rowSpan={ticket.files.length}>
                  {ticket.estimatedExpiry
                    ? relativeDateFormatter.format(
                        differenceInDays(ticket.estimatedExpiry, now),
                        "days",
                      )
                    : t("ticket_expiration_type_none")}
                </Table.Td>
              ) : null}
            </Table.Tr>
          )),
        )}
      </Table.Tbody>
    </Table>
  );
}
