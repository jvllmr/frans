import axios, { AxiosRequestConfig } from "axios";
import { ZodSchema } from "zod";

export async function baseFetch<T>(url: string, opts?: AxiosRequestConfig<T>) {
  return axios.get(`${window.fransRootPath}/${url}`, opts);
}

class FetchError extends Error {}

export async function baseFetchJSON<T>(
  url: string,
  schema: ZodSchema<T>,
  opts?: AxiosRequestConfig<T>,
): Promise<T> {
  const resp = await baseFetch(url, opts);

  if (resp.status >= 300) {
    throw new FetchError();
  }
  return schema.parseAsync(resp.data);
}
