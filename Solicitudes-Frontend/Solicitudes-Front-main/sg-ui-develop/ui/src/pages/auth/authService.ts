import axios, { AxiosError, HttpStatusCode } from "axios";
import { api } from "../apiInstance";
import {
  RequestError,
  JwtPayload,
  Token,
  TokenResponse,
  User,
  UserResponse,
  ErrorResponse,
  AutenticarInfo,
} from "./types";
import { jwtDecode } from "jwt-decode";
import {
  setAutenticarStorage,
  setLocalStorage,
} from "./provider/useLocalStorage";

export class AuthService {
  static async getExternalUserInfo(
    authCode: string,
    provider: string,
    redirectUri: string
  ): Promise<JwtPayload> {
    try {
      const token = await api.post<Token>("/auth/login", {
        code: authCode,
        provider,
        redirectUri,
      });

      if (token.access_token) {
        setLocalStorage(token);
        setAutenticarStorage(token);
        window.localStorage.setItem("provider", provider);
        const decoded = jwtDecode<JwtPayload>(token.access_token);
        return decoded;
      } else {
        throw new Error("Error en la autenticación");
      }
    } catch (error) {
      const axiosError = error as AxiosError;

      throw new RequestError(
        axiosError.status,
        "Error en la obtencion del token del proveedor externo."
      );
    }
  }

  static async externalLogout(provider: string, refreshToken: string) {
    try {
      const userInfo = await api.post<AutenticarInfo>("/auth/logout", {
        provider,
      });

      await axios.post(
        `${userInfo.url}/auth/realms/${userInfo.realmId}/protocol/openid-connect/logout`,
        new URLSearchParams({
          client_id: userInfo.clientId,
          client_secret: userInfo.clientSecret,
          refresh_token: refreshToken,
        }),
        {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
        }
      );
    } catch (error) {
      const axiosError = error as AxiosError;

      throw new RequestError(axiosError.status, "logout error");
    }
  }

  static async checkExistingUser(
    cuit: string,
    provider: string
  ): Promise<UserResponse> {
    try {
      return await api.get<UserResponse>(
        "/admin/users?cuit=" + cuit + "&provider=" + provider
      );
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.status === HttpStatusCode.NotFound) {
          return {
            success: false,
            admin: false,
            data: {
              cuit: "",
              dni: "",
              first_name: "",
              last_name: "",
              email: "",
              phone: "",
            },
            error: "",
          };
        }
      }

      throw new RequestError(
        HttpStatusCode.InternalServerError,
        "Ocurrió un error en la busqueda de datos, por favor intente mas tarde."
      );
    }
  }

  static async loginWithAuthApi(
    cuit: string,
    provider: string
  ): Promise<TokenResponse> {
    try {
      const response = await api.post<TokenResponse>("/admin/login", {
        cuit: cuit,
        provider,
      });
      return response;
    } catch (error) {
      const axiosError = error as AxiosError;

      throw new RequestError(
        axiosError.status,
        "Ocurrió un error en la busqueda de datos, por favor intente mas tarde."
      );
    }
  }

  static async validateToken() {
    try {
      await api.get("/auth/validate");
    } catch (error) {
      const axiosError = error as AxiosError;

      throw new RequestError(
        axiosError.status,
        "Ocurrió un error en la busqueda de datos, por favor intente mas tarde."
      );
    }
  }

  static async registerUser(userData: User): Promise<User> {
    const provider = window.localStorage.getItem("provider");
    try {
      const response = await api.post<UserResponse>(
        "/admin/users?provider=" + provider,
        userData
      );
      return response.data;
    } catch (error) {
      const axiosError = error as AxiosError;

      let msg = "Error inesperado";
      const status = axiosError.response?.status;

      if (status && status === 400) {
        const data = axiosError.response?.data;
        if (isErrorResponse(data)) {
          switch (data.code) {
            case "USER_ALREADY_EXISTS":
              msg = "El correo ya está registrado. Usa otro email.";
              break;
            case "INVALID_REQUEST_BODY":
              msg = "Información incompleta. Llena todos los campos.";
              break;
            default:
              msg = "Error desconocido. Inténtalo nuevamente.";
          }
        }
      }

      throw new RequestError(status, msg);
    }
  }

  static async getUser(cuit: string): Promise<UserResponse> {
    try {
      return await api.get<UserResponse>(
        "/admin/users?cuit=" + cuit + "&provider=jwt"
      );
    } catch (error) {
      throw new RequestError(400, "Error en la busqueda del usuario.");
    }
  }

  static async activateAccount(token: string) {
    try {
      const response = await api.get("/admin/users/activate?token=" + token);
      return response;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.status === HttpStatusCode.Unauthorized) {
          throw new RequestError(
            error.status,
            "El enlace de activación ha expirado. Por favor, solicita un nuevo correo de activación para completar el registro."
          );
        }
        if (error.status === HttpStatusCode.NotFound) {
          throw new RequestError(error.status, "Usuario inexistente.");
        }
        if (error.status === HttpStatusCode.BadRequest) {
          throw new RequestError(
            error.status,
            "Tu cuenta ya ha sido activada anteriormente. Puedes iniciar sesión con tus credenciales."
          );
        }
      }

      throw new RequestError(
        HttpStatusCode.InternalServerError,
        "La solicitud no tuvo éxito."
      );
    }
  }

  static async resendActivationEmail(token: string) {
    try {
      const response = await api.post("/admin/users/resend-activation-email", {
        token,
      });
      return response;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.status === HttpStatusCode.Unauthorized) {
          throw new RequestError(
            error.status,
            "Expiró el tiempo de activacón de la cuenta, por favor solicite un nuevo correo de activación."
          );
        }
        if (error.status === HttpStatusCode.NotFound) {
          throw new RequestError(error.status, "Usuario inexistente.");
        }
        if (error.status === HttpStatusCode.BadRequest) {
          throw new RequestError(
            error.status,
            "El usuario ya se encuentra activo, inicie sesion para continuar."
          );
        }
      }

      throw new RequestError(
        HttpStatusCode.InternalServerError,
        "La solicitud no tuvo éxito."
      );
    }
  }
}

function isErrorResponse(data: unknown): data is ErrorResponse {
  return (data as ErrorResponse).code !== undefined;
}
