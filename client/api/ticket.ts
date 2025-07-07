import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod/v4";
import i18n from "~/i18n";
import { errorNotification, successNotification } from "~/util/notifications";
import { ProgressHandle } from "~/util/progress";
import { FetchError, v1Url } from ".";
import { fileSchema } from "./file";
import { userSchema } from "./user";

function v1TicketUrl(url: string) {
  return v1Url("/ticket" + url);
}

export const ticketExpiryType = z.enum(["auto", "single", "none", "custom"]);

export const createTicketSchema = z.object({
  comment: z.string().nullable(),
  email: z.email(i18n.t("email", { ns: "validation" })).nullable(),
  password: z
    .string()
    .min(12, i18n.t("min_length", { ns: "validation" }).replace("#", "12")),
  emailPassword: z.boolean(),
  expiryType: ticketExpiryType,
  expiryTotalDays: z.int(),
  expiryDaysSinceLastDownload: z.int(),
  expiryTotalDownloads: z.int(),
  emailOnDownload: z.email(i18n.t("email", { ns: "validation" })).nullable(),
  files: z.file().array().min(1),
});
export type CreateTicket = z.infer<typeof createTicketSchema>;

export const ticketSchema = z.object({
  id: z.uuid(),
  owner: userSchema,
  files: fileSchema.array(),
  createdAt: z.coerce.date(),
  estimatedExpiry: z.coerce.date(),
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
