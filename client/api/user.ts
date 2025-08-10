import { queryOptions } from "@tanstack/react-query";
import { z } from "zod/v4";
import { baseFetchJSON, v1Url } from ".";
export const usersKey = ["USER"];

export const meKey = [...usersKey, "ME"];

export const publicUserSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  isAdmin: z.boolean(),
  email: z.email(),
});

export type PublicUser = z.infer<typeof publicUserSchema>;

export const userSchema = publicUserSchema.extend({
  submittedTickets: z.int(),
  activeTickets: z.int(),
  submittedGrants: z.int(),
  activeGrants: z.int(),
  totalDataSize: z.int(),
});

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

export async function fetchUsers() {
  return baseFetchJSON(v1UserUrl(""), userSchema.array());
}

export const usersQueryOptions = queryOptions({
  queryKey: usersKey,
  queryFn: fetchUsers,
});
