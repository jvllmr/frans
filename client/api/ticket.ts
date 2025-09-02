import { queryOptions, useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod/v4";
import i18n, { availableLanguages } from "~/i18n";
import { errorNotification, successNotification } from "~/util/notifications";
import { ProgressHandle } from "~/util/progress";
import { baseFetchJSON, expiryType, FetchError, v1Url } from ".";
import { fileSchema } from "./file";
import { publicUserSchema } from "./user";

export const ticketsKey = ["TICKET"];

function v1TicketUrl(url: string) {
  return v1Url("/ticket" + url);
}

export const createTicketSchemaFactory = (t: typeof i18n.t) =>
  z.object({
    comment: z.string().nullable(),
    email: z.email(t("email", { ns: "validation" })).nullable(),
    password: z
      .string()
      .min(12, i18n.t("min_length", { ns: "validation" }).replace("#", "12")),
    emailPassword: z.boolean(),
    expiryType: expiryType,
    expiryTotalDays: z.int(),
    expiryDaysSinceLastDownload: z.int(),
    expiryTotalDownloads: z.int(),
    emailOnDownload: z.email(i18n.t("email", { ns: "validation" })).nullable(),
    files: z.file().array().min(1),
    creatorLang: z.enum(availableLanguages),
    receiverLang: z.enum(availableLanguages),
  });

export const createTicketSchema = createTicketSchemaFactory(i18n.t);
export type CreateTicket = z.infer<typeof createTicketSchema>;

export const ticketSchema = z.object({
  id: z.uuid(),
  owner: publicUserSchema,
  files: fileSchema.array(),
  createdAt: z.coerce.date(),
  estimatedExpiry: z.coerce.date().nullable(),
  comment: z.string().nullable(),
});

export type Ticket = z.infer<typeof ticketSchema>;

export async function createTicket(
  data: CreateTicket,
  progressHandle?: ProgressHandle,
) {
  const resp = await axios.postForm(v1TicketUrl(""), data, {
    onUploadProgress(progressEvent) {
      progressHandle?.updateProgressState(progressEvent);
    },
  });

  return ticketSchema.parse(resp.data);
}

export function useCreateTicketMutation(progressHandle?: ProgressHandle) {
  const { t } = useTranslation("notifications");
  const partialCreateTicket = useCallback(
    (data: CreateTicket) => createTicket(data, progressHandle),
    [progressHandle],
  );
  return useMutation<Ticket, FetchError, CreateTicket>({
    mutationFn: partialCreateTicket,
    onSuccess() {
      successNotification(t("ticket_new_success"));
      progressHandle?.setFinished();
    },
    onError() {
      errorNotification(t("ticket_new_failed"));
    },
  });
}

export async function fetchTickets() {
  return baseFetchJSON(v1TicketUrl(""), ticketSchema.array());
}

export const ticketQueryOptions = queryOptions({
  queryKey: ticketsKey,
  queryFn: fetchTickets,
});

export async function fetchTicketShare({
  ticketId,
  password,
}: {
  ticketId: string;
  password: string;
}) {
  return baseFetchJSON(v1Url(`/share/ticket/${ticketId}`), ticketSchema, {
    auth: { username: ticketId, password: password },
  });
}

export async function fetchTicketShareAccessToken({
  ticketId,
  password,
}: {
  ticketId: string;
  password: string;
}) {
  return baseFetchJSON(
    v1Url(`/share/ticket/${ticketId}/token`),
    z.object({ token: z.string() }),
    {
      auth: { username: ticketId, password: password },
    },
  );
}

export function getTicketShareFileUrl({
  fileId,
  ticketId,
}: {
  ticketId: string;
  fileId: string;
}) {
  return v1Url(`/share/ticket/${ticketId}/file/${fileId}`);
}
