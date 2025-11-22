import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import axios, { AxiosError } from "axios";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod/v4";
import i18n, { availableLanguages } from "~/i18n";
import { errorNotification, successNotification } from "~/util/notifications";
import { ProgressHandle } from "~/util/progress";
import { baseFetchJSON, expiryType, FetchError, v1Url } from ".";
import { fileSchema } from "./file";
import { publicUserSchema } from "./user";

export const grantsKey = ["GRANT"];

function v1GrantUrl(url: string) {
  return v1Url("/grant" + url);
}

export const createGrantSchemaFactory = (t: typeof i18n.t) =>
  z.object({
    comment: z.string().nullable(),
    email: z.email(t("email", { ns: "validation" })).nullable(),
    password: z
      .string()
      .min(12, i18n.t("min_length", { ns: "validation" }).replace("#", "12")),
    emailPassword: z.boolean(),
    expiryType: expiryType,
    expiryTotalDays: z.int(),
    expiryDaysSinceLastUpload: z.int(),
    expiryTotalUploads: z.int(),
    fileExpiryType: expiryType,
    fileExpiryTotalDays: z.int(),
    fileExpiryDaysSinceLastDownload: z.int(),
    fileExpiryTotalDownloads: z.int(),
    emailOnUpload: z.email(i18n.t("email", { ns: "validation" })).nullable(),
    creatorLang: z.enum(availableLanguages),
    receiverLang: z.enum(availableLanguages),
  });

export const createGrantSchema = createGrantSchemaFactory(i18n.t);
export type CreateGrant = z.infer<typeof createGrantSchema>;

export const grantSchema = z.object({
  id: z.uuid(),
  owner: publicUserSchema,
  files: fileSchema.array(),
  createdAt: z.coerce.date(),
  estimatedExpiry: z.coerce.date().nullable(),
  comment: z.string().nullable(),
});

export type Grant = z.infer<typeof grantSchema>;

export async function createGrant(data: CreateGrant) {
  const resp = await axios.postForm(v1GrantUrl(""), data);

  return grantSchema.parse(resp.data);
}

export function useCreateGrantMutation() {
  const { t } = useTranslation("notifications");

  return useMutation<Grant, FetchError, CreateGrant>({
    mutationFn: createGrant,
    onSuccess() {
      successNotification(t("grant_new_success"));
    },
    onError() {
      errorNotification(t("grant_new_failed"));
    },
  });
}

export function deleteGrant(grantId: string) {
  return axios.delete(v1GrantUrl(`/${grantId}`));
}

export function useDeleteGrantMutation() {
  const { t } = useTranslation("notifications");
  const queryClient = useQueryClient();
  return useMutation<unknown, AxiosError, string>({
    mutationFn: deleteGrant,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: grantsKey });
      successNotification(t("grant_delete_success"));
    },
    onError() {
      errorNotification(t("grant_delete_failed"));
    },
  });
}

export async function fetchGrants() {
  return baseFetchJSON(v1GrantUrl(""), grantSchema.array());
}

export const grantQueryOptions = queryOptions({
  queryKey: grantsKey,
  queryFn: fetchGrants,
});

export async function fetchGrantShare({
  grantId,
  password,
}: {
  grantId: string;
  password: string;
}) {
  return baseFetchJSON(v1Url(`/share/grant/${grantId}`), grantSchema, {
    auth: { username: grantId, password: password },
  });
}

export async function fetchGrantShareAccessToken({
  grantId,
  password,
}: {
  grantId: string;
  password: string;
}) {
  return baseFetchJSON(
    v1Url(`/share/grant/${grantId}/token`),
    z.object({ token: z.string() }),
    {
      auth: { username: grantId, password: password },
    },
  );
}

export const grantUploadSchema = z.object({ files: z.file().array().min(1) });
export type GrantUpload = z.infer<typeof grantUploadSchema>;

export async function uploadToGrant(
  {
    grantId,
    password,
  }: {
    grantId: string;
    password: string;
  },
  data: GrantUpload,
  progressHandle?: ProgressHandle,
) {
  const resp = await axios.postForm(v1Url(`/share/grant/${grantId}`), data, {
    onUploadProgress(progressEvent) {
      progressHandle?.updateProgressState(progressEvent);
    },
    auth: { username: grantId, password: password },
  });

  return grantSchema.parse(resp.data);
}

export function useGrantUploadMutation(
  auth: {
    grantId: string;
    password: string;
  },
  progressHandle?: ProgressHandle,
) {
  const { t } = useTranslation("notifications");
  const partialUploadToGrant = useCallback(
    (data: GrantUpload) => uploadToGrant(auth, data, progressHandle),
    [auth, progressHandle],
  );
  return useMutation<Grant, FetchError, GrantUpload>({
    mutationFn: partialUploadToGrant,
    onSuccess() {
      successNotification(t("grant_upload_success"));
      progressHandle?.setFinished();
    },
    onError() {
      errorNotification(t("grant_upload_failed"));
    },
  });
}
