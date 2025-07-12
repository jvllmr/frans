import { Button, PasswordInput, Stack } from "@mantine/core";
import { useForm } from "@mantine/form";
import { IconLockOpen } from "@tabler/icons-react";
import { QueryKey, useQuery, useQueryClient } from "@tanstack/react-query";
import React, { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

interface TokenGeneratorProps {
  shareTokenGenerator: (password: string) => Promise<{ token: string }>;
  dataQueryKey: QueryKey;
  password: string;
}

export interface ShareAuthProps<TData>
  extends Omit<TokenGeneratorProps, "password"> {
  children?: React.ReactNode;
  DataContextProvider: React.Provider<TData | null>;
  dataFetcher: (password: string) => Promise<TData>;
  prompt: React.ReactNode;
  submitButtonLabel: React.ReactNode;
}

function TokenGenerator({
  shareTokenGenerator,
  dataQueryKey,
  password,
}: TokenGeneratorProps) {
  const tokenQueryKey = useMemo(
    () => [...dataQueryKey, "TOKEN"],
    [dataQueryKey],
  );
  const queryFn = useCallback(
    () => shareTokenGenerator(password),
    [password, shareTokenGenerator],
  );

  useQuery({ queryKey: tokenQueryKey, queryFn, refetchInterval: 9_500 });

  return null;
}

export function ShareAuth<TData>({
  children,
  DataContextProvider,
  dataFetcher,
  prompt,
  shareTokenGenerator,
  dataQueryKey,
  submitButtonLabel,
}: ShareAuthProps<TData>) {
  const form = useForm({ initialValues: { password: "" } });
  const [password, setPassword] = useState<string | null>(null);
  const { t } = useTranslation("share");
  const fetchData = useCallback(() => {
    if (password) {
      return dataFetcher(password);
    }
  }, [dataFetcher, password]);
  const queryClient = useQueryClient();
  const { data, error, isPending } = useQuery({
    queryKey: dataQueryKey,
    queryFn: fetchData,
    enabled: !!password,
  });

  if (password && data) {
    return (
      <>
        <TokenGenerator
          dataQueryKey={dataQueryKey}
          password={password}
          shareTokenGenerator={shareTokenGenerator}
        />
        <DataContextProvider value={data}>{children}</DataContextProvider>
      </>
    );
  }

  return (
    <form
      onSubmit={form.onSubmit(({ password }) => {
        setPassword(password);
        queryClient.invalidateQueries({ queryKey: dataQueryKey });
      })}
    >
      <Stack>
        {prompt}
        <PasswordInput
          {...form.getInputProps("password")}
          error={error ? t("password", { ns: "validation" }) : undefined}
        />
        <Button
          type="submit"
          loading={!!password && isPending}
          leftSection={<IconLockOpen />}
        >
          {submitButtonLabel}
        </Button>
      </Stack>
    </form>
  );
}
