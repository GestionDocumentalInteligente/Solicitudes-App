import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { getUser } from "@/pages/auth/provider/useLocalStorage.ts";

interface RequestOptions {
  method: string;
  headers: {
    "Content-Type": string;
    Authorization: string;
  };
  body?: string;
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getRequestOptions(
  method: string,
  body: string
): RequestOptions {
  const user = getUser();
  const options: RequestOptions = {
    method,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${user?.access_token}`,
    },
  };

  if (method === "GET" || method === "DELETE") {
    return options;
  }

  options.body = body;

  return options;
}

export function base64ToFile(base64: string, fileName: string): File {
  const binary = atob(base64);
  const array = Uint8Array.from(binary, (c) => c.charCodeAt(0));
  return new File([array], fileName, { type: "application/pdf" });
}

export function replacePlaceholders(
  content: string,
  replaceObject: { [key: string]: string }
) {
  return content.replace(/\[([^\]]+)]/g, (match, key) => {
    return replaceObject[key.trim()] || match;
  });
}
