import { Text } from "@mantine/core";
import { IconCheck } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import { useDateFormatter } from "~/i18n";
export interface DownloadSuccessIndicatorProps {
  timesDownloaded: number;
  lastDownloaded: Date | null;
}

export function DownloadSuccessIndicator({
  lastDownloaded,
  timesDownloaded,
}: DownloadSuccessIndicatorProps) {
  const { t } = useTranslation();
  const dateFormatter = useDateFormatter();
  return (
    <>
      {timesDownloaded > 0 && lastDownloaded ? (
        <Text
          span
          c="teal"
          title={`${t("file_downloaded_success")} (${dateFormatter.format(lastDownloaded)})`}
        >
          <IconCheck />
        </Text>
      ) : null}
    </>
  );
}
