import { Table } from "@mantine/core";
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Grant, grantQueryOptions, grantsKey } from "~/api/grant";
import { EstimatedExpiry } from "~/components/common/EstimatedExpiry";
import { DownloadSuccessIndicator } from "~/components/file/DownloadSuccessIndicator";
import { FileRef } from "~/components/file/FileRef";
import { ShareLinkButtons } from "~/components/share/ShareLink";
import { useDateFormatter } from "~/i18n";
import { getInternalFileLink } from "~/util/link";

export const Route = createFileRoute("/grants/active")({
  component: RouteComponent,
});

function GrantButtons({ grant }: { grant: Grant }) {
  return <ShareLinkButtons shareId={grant.id} />;
}

function RouteComponent() {
  const { t } = useTranslation("grant_active");
  const { data: grants } = useSuspenseQuery(grantQueryOptions);
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
                  <Table.Td>{grant.id}</Table.Td>
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
