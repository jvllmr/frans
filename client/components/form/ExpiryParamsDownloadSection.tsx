import { Fieldset, Group, NumberInput, Select } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { ExpiryType } from "~/api";

interface BaseProps {
  label: NonNullable<React.ReactNode>;
}

interface TicketExpiryParamsDownloadSectionProps extends BaseProps {
  form: UseFormReturnType<{
    expiryType: ExpiryType;
    expiryTotalDays: number;
    expiryDaysSinceLastDownload: number;
    expiryTotalDownloads: number;
  }>;
  variant: "ticket";
}

interface GrantExpiryParamsDownloadSectionProps extends BaseProps {
  form: UseFormReturnType<{
    fileExpiryType: ExpiryType;
    fileExpiryTotalDays: number;
    fileExpiryDaysSinceLastDownload: number;
    fileExpiryTotalDownloads: number;
  }>;
  variant: "grant";
}

export type ExpiryParamsDownloadSectionProps =
  | TicketExpiryParamsDownloadSectionProps
  | GrantExpiryParamsDownloadSectionProps;

export function ExpiryParamsDownloadSection({
  form,
  variant,
  label,
}: ExpiryParamsDownloadSectionProps) {
  const { t } = useTranslation("forms");
  const expiryChoices: { label: string; value: ExpiryType }[] = useMemo(
    () => [
      {
        value: "auto",
        label: t("expiry_automatic"),
      },
      {
        value: "single",
        label: t("expiry_single_use"),
      },
      {
        value: "none",
        label: t("expiry_none"),
      },
      {
        value: "custom",
        label: t("expiry_custom"),
      },
    ],
    [t],
  );
  const expiryTypeProps =
    variant === "grant"
      ? form.getInputProps("fileExpiryType")
      : form.getInputProps("expiryType");
  const expiryTotalDaysProps =
    variant === "grant"
      ? form.getInputProps("fileExpiryTotalDays")
      : form.getInputProps("expiryTotalDays");
  const expiryTotalDownloadsProps =
    variant === "grant"
      ? form.getInputProps("fileExpiryTotalDownloads")
      : form.getInputProps("expiryTotalDownloads");
  const expiryDaysSinceLastDownloadProps =
    variant === "grant"
      ? form.getInputProps("fileExpiryDaysSinceLastDownload")
      : form.getInputProps("expiryDaysSinceLastDownload");

  const isExpiryTypeCustom =
    variant === "grant"
      ? form.values.fileExpiryType === "custom"
      : form.values.expiryType === "custom";
  return (
    <Fieldset>
      <Select {...expiryTypeProps} data={expiryChoices} label={label} />
      {isExpiryTypeCustom ? (
        <Group mt="xs" align="end" grow>
          <NumberInput
            {...expiryTotalDaysProps}
            label={t("label_expiry_total_days")}
          />
          <NumberInput
            {...expiryDaysSinceLastDownloadProps}
            label={t("label_expiry_last_download")}
          />
          <NumberInput
            {...expiryTotalDownloadsProps}
            label={t("label_expiry_total_downloads")}
          />
        </Group>
      ) : null}
    </Fieldset>
  );
}
