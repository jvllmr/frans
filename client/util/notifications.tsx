import { showNotification } from "@mantine/notifications";
import { IconCheck, IconX } from "@tabler/icons-react";

export function successNotification(message: string) {
  showNotification({ color: "teal", icon: <IconCheck />, message });
}

export function errorNotification(message: string) {
  showNotification({ color: "red", icon: <IconX />, message });
}
