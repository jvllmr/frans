import { createContext, useContext } from "react";

export const shareAuthContext = createContext<{
  password: string;
} | null>(null);

export function useShareAuthContext() {
  const value = useContext(shareAuthContext);
  if (!value) throw TypeError("Expected share auth context.");
  return value;
}
