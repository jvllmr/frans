import { createTheme } from "@mantine/core";

export const BASE_THEME = createTheme({
  primaryColor: window.fransColor,
  colors: {
    custom: JSON.parse(window.fransCustomColor),
  },
});
