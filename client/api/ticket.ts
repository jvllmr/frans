import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { z } from "zod/v4";
import i18n from "~/i18n";
import { FetchError, v1Url } from ".";

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
  userId: z.uuid(),
  createdAt: z.date(),
  expiryType: ticketExpiryType,
  expiryTotalDays: z.int(),
  expiryDaysSinceLastDownload: z.int(),
  expiryTotalDownloads: z.int(),
  lastDownload: z.date().nullable(),
});

export type Ticket = z.infer<typeof ticketSchema>;

export async function createTicket(data: CreateTicket): Promise<Ticket> {
  return (await axios.postForm(v1TicketUrl(""), data)).data;
}

export function useCreateTicketMutation() {
  return useMutation<Ticket, FetchError, CreateTicket>({
    mutationFn: createTicket,
  });
}
