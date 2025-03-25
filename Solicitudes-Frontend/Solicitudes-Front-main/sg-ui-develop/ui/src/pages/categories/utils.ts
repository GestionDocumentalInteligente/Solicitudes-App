import { getUser } from "../auth/provider/useLocalStorage";
import { Category, CategoryResponse, RequestError } from "./types";

interface RequestOptions {
  method: string;
  headers: {
    "Content-Type": string;
    Authorization: string;
  };
  body?: string;
}

const getRequestOptions = (method: string, body: string) => {
  const user = getUser();
  const options: RequestOptions = {
    method: method,
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
};

export const getCategoriesInfo = async (
  params?: string
): Promise<Category[]> => {
  params = params || "";

  const response = await fetch(
    "/api/admin/categories" + params,
    getRequestOptions("GET", "")
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  const result: CategoryResponse = await response.json();
  if (response.status >= 200 && response.status < 300) {
    return result.data;
  } else {
    const { error } = result;
    throw new RequestError(error);
  }
};
