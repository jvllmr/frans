import { Anchor, List, Stack, Text } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import React, { useCallback, useContext, useMemo } from "react";
import { useTranslation } from "react-i18next";
import z from "zod/v4";
import {
  fetchTicketShare,
  fetchTicketShareAccessToken,
  getTicketShareFileUrl,
  Ticket,
} from "~/api/ticket";
import { ShareAuth } from "~/components/ShareAuth";
import { useFileSizeFormatter } from "~/i18n";

export const Route = createFileRoute("/share/ticket/$ticketId")({
  component: RouteComponent,
  params: { parse: z.object({ ticketId: z.uuid() }).parse },
});

const shareTicketContext = React.createContext<Ticket | null>(null);

function useShareTicketContext() {
  const ticket = useContext(shareTicketContext);
  if (!ticket) throw TypeError("Expected ticket to be available");
  return ticket;
}

function TicketShare() {
  const ticket = useShareTicketContext();
  const { t } = useTranslation("share");
  const fileSizeFormatter = useFileSizeFormatter();
  return (
    <Stack>
      <Text>
        <b>
          {ticket.owner.name} ({ticket.owner.email})
        </b>{" "}
        {t("ticket_message")}:{" "}
      </Text>
      <List px="xl">
        {ticket.files.map((file) => (
          <List.Item key={file.id}>
            <Anchor
              href={getTicketShareFileUrl({
                ticketId: ticket.id,
                fileId: file.id,
              })}
            >
              {file.name} ({fileSizeFormatter(file.size)})
            </Anchor>
          </List.Item>
        ))}
      </List>
      <Text>
        <b>{t("comment")}: </b>
        {ticket.comment}
      </Text>
    </Stack>
  );
}

function RouteComponent() {
  const ticketId = Route.useParams({ select: (p) => p.ticketId });
  const dataFetcher = useCallback(
    (password: string) => fetchTicketShare({ ticketId, password }),
    [ticketId],
  );
  const queryKey = useMemo(() => ["SHARE", "TICKET", ticketId], [ticketId]);
  const shareTokenGenerator = useCallback(
    (password: string) => fetchTicketShareAccessToken({ ticketId, password }),
    [ticketId],
  );
  const { t } = useTranslation("share");
  return (
    <ShareAuth
      DataContextProvider={shareTicketContext.Provider}
      dataFetcher={dataFetcher}
      dataQueryKey={queryKey}
      shareTokenGenerator={shareTokenGenerator}
      prompt={t("ticket_prompt")}
      submitButtonLabel={t("ticket_submit")}
    >
      <TicketShare />
    </ShareAuth>
  );
}
