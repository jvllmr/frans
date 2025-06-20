import { useQuery } from "@tanstack/react-query";
import { meQueryOptions } from "~/api/user";

export function useIsAdmin(): boolean {
  const { data: me } = useQuery(meQueryOptions);

  return !!me && me.isAdmin;
}
