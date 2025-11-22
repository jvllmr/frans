import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import axios, { AxiosError } from "axios";
import { useTranslation } from "react-i18next";
import { z } from "zod/v4";
import { errorNotification, successNotification } from "~/util/notifications";
import { baseFetchJSON, v1Url } from ".";

import { publicUserSchema } from "./user";

export const filesKey = ["FILES"];

function v1FileUrl(url: string) {
  return v1Url("/file" + url);
}

export const fileSchema = z.object({
  id: z.string(),
  sha512: z.string(),
  name: z.string(),
  size: z.int(),
  timesDownloaded: z.int(),
  createdAt: z.coerce.date(),
  lastDownloaded: z.coerce.date().nullable(),
  estimatedExpiry: z.coerce.date().nullable(),
  owner: publicUserSchema,
});

export function fetchReceivedFiles() {
  return baseFetchJSON(v1FileUrl("/received"), fileSchema.array());
}

export const receivedFilesQueryOptions = queryOptions({
  queryKey: [...filesKey, "RECEIVED"],
  queryFn: fetchReceivedFiles,
});

export function deleteFile(fileId: string) {
  return axios.delete(v1FileUrl(`/${fileId}`));
}

export function useDeleteFileMutation() {
  const { t } = useTranslation("notifications");
  const queryClient = useQueryClient();
  return useMutation<unknown, AxiosError, string>({
    mutationFn: deleteFile,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: filesKey });

      successNotification(t("file_delete_success"));
    },
    onError() {
      errorNotification(t("file_delete_failed"));
    },
  });
}
