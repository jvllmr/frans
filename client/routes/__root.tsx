import {
  Box,
  Button,
  Center,
  Container,
  Divider,
  Flex,
  MantineProvider,
  Paper,
  SegmentedControl,
  Text,
  Title,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import {
  QueryClient,
  QueryClientProvider,
  queryOptions,
  useQuery,
} from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import {
  createRootRouteWithContext,
  Outlet,
  Register,
  useChildMatches,
  useNavigate,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { isAxiosError } from "axios";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { queryClient } from "~/api";
import { meQueryOptions } from "~/api/user";
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
  { translationKey: string; needsAdmin: boolean }
>;
const tabTitles = {
  "/": { translationKey: "new", needsAdmin: false },
  "/tickets/": { translationKey: "tickets", needsAdmin: false },
  "/grants/new": { translationKey: "new_grant", needsAdmin: false },
  "/grants/active": { translationKey: "active_grants", needsAdmin: false },
  "/grants/": { translationKey: "grants", needsAdmin: false },
  "/users": { translationKey: "users", needsAdmin: true },
} satisfies Partial<TabTitles>;

const hiddenTabTitles: TabTitles = {
  ...tabTitles,
  "/s/$shareId": { translationKey: "share", needsAdmin: false },
  "/share/ticket/$ticketId": {
    translationKey: "ticket_share",
    needsAdmin: false,
  },
};

const cyclicMeQueryOptions = queryOptions({
  ...meQueryOptions,
  refetchInterval: 30_000, // 30 seconds
  refetchIntervalInBackground: true,
});

function useDeepestMatch() {
  return useChildMatches({ select: (m) => m.at(-1) });
}

function useAuthRequired() {
  const deepestMatch = useDeepestMatch();
  return (
    !!deepestMatch && Object.keys(tabTitles).includes(deepestMatch.routeId)
  );
}

function AuthGuard({ children }: { children: React.ReactNode }) {
  const authRequired = useAuthRequired();
  const { t } = useTranslation();
  const { error } = useQuery(cyclicMeQueryOptions);

  if (authRequired && error && isAxiosError(error) && error.status === 401) {
    return (
      <>
        <Center>
          <Text fw="bold">{t("session_expired")}</Text>
        </Center>
        <Center>
          <Button
            onClick={() => {
              window.location.reload();
            }}
          >
            {t("session_renew")}
          </Button>
        </Center>
      </>
    );
  }

  return children;
}

function TabControls() {
  const { t } = useTranslation("tabs");
  const authRequired = useAuthRequired();
  const { data: me } = useQuery(meQueryOptions);
  const data = useMemo(
    () =>
      Object.entries(tabTitles)
        .filter(
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
          ([_, { needsAdmin }]) => !needsAdmin || (needsAdmin && me?.isAdmin),
        )
        .map(([value, { translationKey }]) => ({
          value,
          label: t(translationKey),
        })),
    [me?.isAdmin, t],
  );
  const navigate = useNavigate();
  const deepestMatch = useDeepestMatch();

  if (!me && !authRequired) {
    return null;
  }

  return (
    <SegmentedControl
      data={data}
      value={deepestMatch!.routeId}
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
      title={t("title_logout")}
    >
      {t("logout")}
    </Button>
  );
}

function TabTitle() {
  const { t } = useTranslation("tabs");
  const deepestMatch = useDeepestMatch();

  return (
    <Title size="h3">
      {t(
        // @ts-expect-error we know that id is correct
        hiddenTabTitles[deepestMatch!.routeId].translationKey,
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
          <Container pt={50}>
            <Paper withBorder p="lg" mb="xs">
              <AuthGuard>
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
                <Outlet />
              </AuthGuard>
            </Paper>
          </Container>
          <DevTools />
        </QueryClientProvider>
        <Notifications position="bottom-center" />
      </MantineProvider>
    </>
  );
}

interface RoutingContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RoutingContext>()({
  component: RootRoute,
});
