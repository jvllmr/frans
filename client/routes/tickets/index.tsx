import {
  ActionIcon,
  Anchor,
  CopyButton,
  Group,
  Table,
  Text,
} from "@mantine/core";
import {
  IconCheck,
  IconCopy,
  IconCopyCheck,
  IconFolderOpen,
} from "@tabler/icons-react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { differenceInDays } from "date-fns";
import { useTranslation } from "react-i18next";
import { ticketQueryOptions } from "~/api/ticket";
import { ActionIconLink } from "~/components/Link";
import { useDateFormatter, useRelativeDateFormatter } from "~/i18n";
import { getShareLink } from "~/util/link";
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
                  <Group>
                    <ActionIconLink
                      to="/s/$shareId"
                      params={{ shareId: ticket.id }}
                      target="_blank"
                      title={t("title_open_share")}
                    >
                      <IconFolderOpen />
                    </ActionIconLink>
                    <CopyButton value={getShareLink(ticket.id)}>
                      {({ copied, copy }) => (
                        <ActionIcon
                          onClick={() => {
                            copy();
                          }}
                          title={t("title_copy_link")}
                        >
                          {copied ? <IconCopyCheck /> : <IconCopy />}
                        </ActionIcon>
                      )}
                    </CopyButton>
                  </Group>
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
