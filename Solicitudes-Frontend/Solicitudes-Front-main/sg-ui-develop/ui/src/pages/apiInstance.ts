import axios, { AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { getUser } from "./auth/provider/useLocalStorage";
import { Token } from "./auth/types";

declare module "axios" {
  export interface InternalAxiosRequestConfig {
    _retry?: boolean;
  }
}

class APIClient {
  private user: Token | null;
  private client: AxiosInstance;

  constructor() {
    this.user = null;
    this.client = axios.create({
      baseURL: "/api",
    });

    this.client.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        this.user = getUser();
        const token = this.user?.access_token;
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        if (!config.headers["Content-Type"]) {
          config.headers["Content-Type"] = "application/json";
        }

        return config;
      }
    );
  }

  async get<T>(endpoint: string, params?: object): Promise<T> {
    const response = await this.client.get(endpoint, { params });
    return response.data;
  }

  async post<T>(endpoint: string, data: unknown): Promise<T> {
    let headers = {};

    if (data instanceof FormData) {
      headers = { "Content-Type": "multipart/form-data" };
    }

    const response = await this.client.post(endpoint, data, { headers });
    return response.data;
  }

  async put<T>(endpoint: string, data: unknown): Promise<T> {
    let headers = {};

    if (data instanceof FormData) {
      headers = { "Content-Type": "multipart/form-data" };
    }

    const response = await this.client.put(endpoint, data, { headers });
    return response.data;
  }
}

export const api = new APIClient();
