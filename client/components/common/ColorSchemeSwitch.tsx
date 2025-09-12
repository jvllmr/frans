import { ActionIcon, useMantineColorScheme } from "@mantine/core";
import { IconMoon, IconSun } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";

export function ColorSchemeSwitch() {
  const { t } = useTranslation();
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();
  return (
    <ActionIcon
      variant="light"
      size="sm"
      title={t(
        colorScheme === "light"
          ? "switch_color_scheme_dark"
          : "switch_color_scheme_light",
      )}
      color="gray"
      onClick={() => {
        toggleColorScheme();
      }}
    >
      {colorScheme === "light" ? <IconSun /> : <IconMoon />}
    </ActionIcon>
  );
}
