import { Group, Table } from "@mantine/core";
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  Grant,
  grantQueryOptions,
  grantsKey,
  useDeleteGrantMutation,
} from "~/api/grant";
import { meQueryOptions } from "~/api/user";
import { DeleteButton } from "~/components/common/DeleteButton";
import { EstimatedExpiry } from "~/components/common/EstimatedExpiry";
import { DownloadSuccessIndicator } from "~/components/file/DownloadSuccessIndicator";
import { FileRef } from "~/components/file/FileRef";
import { ShareLinkButtons } from "~/components/share/ShareLink";
import { useDateFormatter } from "~/i18n";
import { getInternalFileLink } from "~/util/link";

export const Route = createFileRoute("/grants/active")({
  component: RouteComponent,
});

function DeleteGrantButton({ grantId }: { grantId: string }) {
  const mutation = useDeleteGrantMutation();
  const { t } = useTranslation("grant_active");
  return (
    <DeleteButton
      loading={mutation.isPending}
      onClick={() => {
        mutation.mutate(grantId);
      }}
      title={t("title_delete")}
    />
  );
}

function GrantButtons({ grant }: { grant: Grant }) {
  return (
    <Group>
      <DeleteGrantButton grantId={grant.id} />
      <ShareLinkButtons shareId={grant.id} />
    </Group>
  );
}

function RouteComponent() {
  const { t } = useTranslation("grant_active");
  const { data: grants } = useSuspenseQuery(grantQueryOptions);
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
        {grants.flatMap((grant) =>
          grant.files.length > 0
            ? grant.files.map((file, index) => (
                <Table.Tr key={file.id}>
                  {index === 0 ? (
                    <Table.Td rowSpan={grant.files.length}>
                      <GrantButtons grant={grant} />
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
                      link={getInternalFileLink(file.id, true)}
                      onClick={() => {
                        queryClient.invalidateQueries({ queryKey: grantsKey });
                      }}
                    />
                  </Table.Td>
                  {index === 0 && me.isAdmin ? (
                    <Table.Td rowSpan={grant.files.length}>
                      {grant.owner.name}
                    </Table.Td>
                  ) : null}
                  {index === 0 ? (
                    <Table.Td rowSpan={grant.files.length}>
                      {dateFormatter.format(grant.createdAt)}
                    </Table.Td>
                  ) : null}
                  {index === 0 ? (
                    <Table.Td rowSpan={grant.files.length}>
                      <EstimatedExpiry
                        estimatedExpiry={grant.estimatedExpiry}
                      />
                    </Table.Td>
                  ) : null}
                </Table.Tr>
              ))
            : [
                <Table.Tr key={grant.id}>
                  <Table.Td colSpan={2}>
                    <GrantButtons grant={grant} />
                  </Table.Td>
                  <Table.Td />
                  <Table.Td>{grant.id}</Table.Td>
                  {me.isAdmin ? <Table.Td>{grant.owner.name}</Table.Td> : null}
                  <Table.Td>{dateFormatter.format(grant.createdAt)}</Table.Td>
                  <Table.Td>
                    <EstimatedExpiry estimatedExpiry={grant.estimatedExpiry} />
                  </Table.Td>
                </Table.Tr>,
              ],
        )}
      </Table.Tbody>
    </Table>
  );
}
