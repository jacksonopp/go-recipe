import { type H3Event, appendResponseHeader } from "h3";
import { type NitroFetchOptions } from "nitropack";

export const fetchWithCookie = async (
  event: H3Event,
  url: string,
  opts?: NitroFetchOptions<string>,
) => {
  const res = await $fetch.raw(url, opts);
  const cookies = (res.headers.get("set-cookie") || "").split(",");

  for (const cookie of cookies) {
    appendResponseHeader(event, "set-cookie", cookie);
  }

  return res._data;
};
