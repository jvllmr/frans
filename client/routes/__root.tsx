import {
  Box,
  Container,
  Divider,
  Flex,
  MantineProvider,
  Paper,
  SegmentedControl,
  Title,
} from "@mantine/core";
import "@mantine/core/styles.css";
import {
  createRootRoute,
  Outlet,
  Register,
  useChildMatches,
  useNavigate,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import i18n, { availableLanguages, availableLanguagesLabels } from "~/i18n";

function LanguageControls() {
  const [language, setLanguage] = useState(i18n.language);
  return (
    <SegmentedControl
      data={availableLanguages.map((lang) => ({
        value: lang,
        label: availableLanguagesLabels[lang],
      }))}
      size="xs"
      radius="xs"
      value={language}
      onChange={async (lang) => {
        setLanguage(lang);
        await i18n.changeLanguage(lang);
      }}
    />
  );
}
type TabTitles = Record<
  Exclude<keyof Register["router"]["routesById"], "__root__">,
  string
>;
const tabTitles = {
  "/": "new",
  "/tickets/": "tickets",
  "/grants/new": "new_grant",
  "/grants/active": "active_grants",
  "/grants/": "grants",
} satisfies Partial<TabTitles>;

const hiddenTabTitles: TabTitles = { ...tabTitles };

function TabControls() {
  const { t } = useTranslation("tabs");
  const data = useMemo(
    () =>
      Object.entries(tabTitles).map(([value, label]) => ({
        value,
        label: t(label),
      })),
    [t],
  );
  const navigate = useNavigate();
  return (
    <SegmentedControl
      data={data}
      onChange={(value) => {
        navigate({ to: value });
      }}
      size="xs"
    />
  );
}

function TabTitle() {
  const { t } = useTranslation("tabs");
  const deepestMatch = useChildMatches({ select: (m) => m.at(-1) });

  return (
    <Title size="h3">
      {t(
        // @ts-expect-error we know that id is correct
        hiddenTabTitles[deepestMatch.id],
      )}
    </Title>
  );
}

function RootRoute() {
  return (
    <>
      <MantineProvider>
        <Container pt={100}>
          <Paper withBorder p="lg">
            <Flex justify="space-between" p={3}>
              <TabControls />
              <LanguageControls />
            </Flex>
            <Box py="sm">
              <TabTitle />
              <Divider />
            </Box>
            <Outlet />
          </Paper>
        </Container>
      </MantineProvider>
      <TanStackRouterDevtools />
    </>
  );
}

export const Route = createRootRoute({
  component: RootRoute,
});
