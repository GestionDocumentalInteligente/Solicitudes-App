import axios, {
  AxiosInstance,
  AxiosResponse,
  AxiosRequestConfig,
  AxiosError,
} from "axios";

export class ApiClient {
  private axiosInstance: AxiosInstance;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({
      baseURL,
      timeout: 600000,
    });

    this.axiosInstance.interceptors.response.use(
      this.handleSuccessResponse,
      this.handleErrorResponse
    );
  }

  private handleSuccessResponse(response: AxiosResponse): AxiosResponse {
    return response;
  }

  private handleErrorResponse(error: AxiosError): Promise<never> {
    if (error.response) {
      return Promise.reject({
        status: error.response.status,
        data: error.response.data,
      });
    }
    return Promise.reject(error);
  }

  public async get<T>(
    url: string,
    headers?: Record<string, string | undefined>
  ): Promise<T> {
    const config: AxiosRequestConfig = headers ? { headers } : {};
    const response = await this.axiosInstance.get<T>(url, config);
    return response.data;
  }

  public async post<T>(
    url: string,
    data: any,
    headers?: Record<string, string | undefined>
  ): Promise<T> {
    const config: AxiosRequestConfig = headers ? { headers } : {};
    const response = await this.axiosInstance.post<T>(url, data, config);
    return response.data;
  }

  public async put<T>(
      url: string,
      data: any,
      headers?: Record<string, string | undefined>
  ): Promise<T> {
    const config: AxiosRequestConfig = headers ? { headers } : {};
    const response = await this.axiosInstance.put<T>(url, data, config);
    return response.data;
  }
}
