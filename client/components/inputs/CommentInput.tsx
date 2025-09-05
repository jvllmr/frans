import { useTranslation } from "react-i18next";
import { NullTextarea, NullTextareaProps } from "./NullTextarea";

export function CommentInput(props: NullTextareaProps) {
  const { t } = useTranslation();
  return (
    <NullTextarea
      label={t("comment", { ns: "translation" })}
      resize="vertical"
      {...props}
    />
  );
}
