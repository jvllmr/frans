import {
  Box,
  Button,
  Container,
  Divider,
  Flex,
  Loader,
  MantineProvider,
  Paper,
  SegmentedControl,
  Title,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import {
  createRootRoute,
  Outlet,
  Register,
  useChildMatches,
  useNavigate,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { Suspense, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { queryClient } from "~/api";
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

function LogoutButton() {
  const { t } = useTranslation();

  return (
    <Button
      size="xs"
      component="a"
      href={`${window.fransRootPath}/api/auth/logout?redirect_uri=${window.location.href}`}
    >
      {t("logout")}
    </Button>
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

function DevTools() {
  return (
    <>
      <TanStackRouterDevtools />
      <ReactQueryDevtools buttonPosition="bottom-right" position="top" />
    </>
  );
}

function RootRoute() {
  return (
    <>
      <MantineProvider>
        <QueryClientProvider client={queryClient}>
          <Container pt={100}>
            <Paper withBorder p="lg" mb="xs">
              <Flex justify="space-between" p={3}>
                <LanguageControls />
                <LogoutButton />
              </Flex>
              <Flex p={3}>
                <TabControls />
              </Flex>
              <Box py="sm">
                <TabTitle />
                <Divider />
              </Box>
              <Suspense
                fallback={
                  <Flex h="60vh" w="100%" justify="center" align="center">
                    <Loader />
                  </Flex>
                }
              >
                <Outlet />
              </Suspense>
            </Paper>
          </Container>
          <DevTools />
        </QueryClientProvider>
        <Notifications position="bottom-center" />
      </MantineProvider>
    </>
  );
}

export const Route = createRootRoute({
  component: RootRoute,
});
