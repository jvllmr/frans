import { ActionIcon, ActionIconProps } from "@mantine/core";
import { createLink } from "@tanstack/react-router";

export const ActionIconLink = createLink(
  (props: Omit<ActionIconProps, "component"> & { title?: string }) => (
    <ActionIcon component="a" {...props} />
  ),
);
