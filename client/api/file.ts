import { queryOptions } from "@tanstack/react-query";
import { z } from "zod/v4";
import { baseFetchJSON, v1Url } from ".";

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
});

export function fetchReceivedFiles() {
  return baseFetchJSON(v1FileUrl("/received"), fileSchema.array());
}

export const receivedFilesQueryOptions = queryOptions({
  queryKey: [...filesKey, "RECEIVED"],
  queryFn: fetchReceivedFiles,
});
