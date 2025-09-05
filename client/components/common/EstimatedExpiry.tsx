import { differenceInDays } from "date-fns/differenceInDays";
import { useTranslation } from "react-i18next";
import { useRelativeDateFormatter } from "~/i18n";

interface EstimatedExpiryProps {
  estimatedExpiry: Date | null;
}

export function EstimatedExpiry({ estimatedExpiry }: EstimatedExpiryProps) {
  const now = new Date();
  const relativeDateFormatter = useRelativeDateFormatter();
  const { t } = useTranslation();

  return (
    <>
      {estimatedExpiry
        ? relativeDateFormatter.format(
            differenceInDays(estimatedExpiry, now),
            "days",
          )
        : t("expiration_type_none")}
    </>
  );
}
