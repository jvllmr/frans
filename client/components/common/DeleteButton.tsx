import { ActionIcon, ActionIconProps } from "@mantine/core";
import { IconX } from "@tabler/icons-react";
import { ButtonHTMLAttributes, DetailedHTMLProps } from "react";

export const DeleteButton: React.FC<
  DetailedHTMLProps<
    ButtonHTMLAttributes<HTMLButtonElement>,
    HTMLButtonElement
  > &
    ActionIconProps
> = (props) => {
  return (
    <ActionIcon variant="light" color="red" {...props}>
      <IconX />
    </ActionIcon>
  );
};
