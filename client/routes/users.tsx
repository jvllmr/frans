import { Table, Text } from "@mantine/core";
import { IconCheck } from "@tabler/icons-react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { usersQueryOptions } from "~/api/user";
import { useFileSizeFormatter } from "~/i18n";

export const Route = createFileRoute("/users")({
  component: RouteComponent,
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(usersQueryOptions),
});

function RouteComponent() {
  const { data: users } = useSuspenseQuery(usersQueryOptions);
  const { t } = useTranslation("users");
  const fileSizeFormatter = useFileSizeFormatter();
  return (
    <Table withColumnBorders withRowBorders withTableBorder>
      <Table.Thead>
        <Table.Tr>
          <Table.Td>{t("user")}</Table.Td>
          <Table.Td>{t("admin")}</Table.Td>
          <Table.Td>{t("tickets_submitted")}</Table.Td>
          <Table.Td>{t("tickets_active")}</Table.Td>
          <Table.Td>{t("grants_submitted")}</Table.Td>
          <Table.Td>{t("grants_active")}</Table.Td>
          <Table.Td>{t("current_size")}</Table.Td>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>
        {users.map((user) => (
          <Table.Tr key={`user_${user.id}`}>
            <Table.Td>{user.name}</Table.Td>
            <Table.Td>
              {user.isAdmin ? (
                <Text c="teal">
                  <IconCheck />
                </Text>
              ) : null}
            </Table.Td>
            <Table.Td>{user.submittedTickets}</Table.Td>
            <Table.Td>{user.activeTickets}</Table.Td>
            <Table.Td>{user.submittedGrants}</Table.Td>
            <Table.Td>{user.activeGrants}</Table.Td>
            <Table.Td>{fileSizeFormatter(user.totalDataSize)}</Table.Td>
          </Table.Tr>
        ))}
      </Table.Tbody>
    </Table>
  );
}
