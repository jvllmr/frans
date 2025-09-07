import {
  ActionIcon,
  ActionIconProps,
  CopyButton,
  Group,
  GroupProps,
} from "@mantine/core";
import { IconCopy, IconCopyCheck, IconFolderOpen } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import { getShareLink } from "~/util/link";
import { successNotification } from "~/util/notifications";
import { ActionIconLink } from "../routing/Link";

interface ShareLinkProps {
  shareId: string;
}

export function ShareLinkButton({
  shareId,
  ...props
}: ShareLinkProps & Omit<ActionIconProps, "title">) {
  const { t } = useTranslation("comps");

  return (
    <ActionIconLink
      {...props}
      to="/s/$shareId"
      params={{ shareId }}
      target="_blank"
      title={t("title_open_share")}
    >
      <IconFolderOpen />
    </ActionIconLink>
  );
}

export function CopyShareLinkButton({
  shareId,
  ...props
}: ShareLinkProps & Omit<ActionIconProps, "onClick" | "title">) {
  const { t } = useTranslation("comps");
  return (
    <CopyButton value={getShareLink(shareId)}>
      {({ copied, copy }) => (
        <ActionIcon
          {...props}
          onClick={() => {
            copy();
            successNotification(t("link_copied"));
          }}
          title={t("title_copy_share_link")}
        >
          {copied ? <IconCopyCheck /> : <IconCopy />}
        </ActionIcon>
      )}
    </CopyButton>
  );
}

export function ShareLinkButtons({
  shareId,
  ...props
}: ShareLinkProps & GroupProps) {
  return (
    <Group {...props}>
      <ShareLinkButton shareId={shareId} />
      <CopyShareLinkButton shareId={shareId} />
    </Group>
  );
}
