import { QueryClient } from "@tanstack/react-query";
import axios, { AxiosRequestConfig } from "axios";
import { ZodType } from "zod/v4";

export function v1Url(url: string) {
  return `${window.fransRootPath}/api/v1${url}`;
}

export async function baseFetch<T>(url: string, opts?: AxiosRequestConfig<T>) {
  return axios.get(
    url.startsWith(window.fransRootPath)
      ? url
      : `${window.fransRootPath}/${url}`,
    opts,
  );
}

export class FetchError extends Error {
  statusCode: number;
  constructor(message: string, statusCode: number) {
    super(message);
    this.statusCode = statusCode;
  }
}

export async function baseFetchJSON<T>(
  url: string,
  schema: ZodType<T>,
  opts?: AxiosRequestConfig<T>,
): Promise<T> {
  const resp = await baseFetch(url, opts);

  if (resp.status >= 400) {
    throw new FetchError(resp.statusText, resp.status);
  }
  return schema.parseAsync(resp.data);
}

export const queryClient = new QueryClient();
