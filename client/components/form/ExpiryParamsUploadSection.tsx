import { Fieldset, Group, NumberInput, Select } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { ExpiryType } from "~/api";

export interface ExpiryParamsUploadSectionProps {
  form: UseFormReturnType<{
    expiryType: ExpiryType;
    expiryTotalDays: number;
    expiryDaysSinceLastUpload: number;
    expiryTotalUploads: number;
  }>;
  label: NonNullable<React.ReactNode>;
}

export function ExpiryParamsUploadSection({
  form,

  label,
}: ExpiryParamsUploadSectionProps) {
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

  return (
    <Fieldset>
      <Select
        {...form.getInputProps("expiryType")}
        data={expiryChoices}
        label={label}
      />
      {form.values.expiryType === "custom" ? (
        <Group mt="xs" align="end" grow>
          <NumberInput
            {...form.getInputProps("expiryTotalDays")}
            label={t("label_expiry_total_days")}
          />
          <NumberInput
            {...form.getInputProps("expiryDaysSinceLastUpload")}
            label={t("label_expiry_last_upload")}
          />
          <NumberInput
            {...form.getInputProps("expiryTotalUploads")}
            label={t("label_expiry_total_uploads")}
          />
        </Group>
      ) : null}
    </Fieldset>
  );
}
