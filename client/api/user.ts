import { queryOptions } from "@tanstack/react-query";
import { z } from "zod/v4";
import { baseFetchJSON, v1Url } from ".";
export const usersKey = ["USER"];

export const meKey = [...usersKey, "ME"];

export const userSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  isAdmin: z.boolean(),
  email: z.email(),
});

export type User = z.infer<typeof userSchema>;

function v1UserUrl(url: string) {
  return v1Url("/user" + url);
}

export async function fetchMe() {
  return baseFetchJSON(v1UserUrl("/me"), userSchema);
}

export const meQueryOptions = queryOptions({
  queryKey: meKey,
  queryFn: fetchMe,
});
