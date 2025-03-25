import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useMemo,
} from "react";
import { useNavigate } from "react-router-dom";
import { HttpStatusCode } from "axios";

import { JwtPayload, User, RequestError, Token } from "../types";
import { AuthService } from "../authService";
import {
  clearLocalStorage,
  getAutenticarInfo,
  getUser,
  setLocalStorage,
} from "./useLocalStorage";
import { jwtDecode } from "jwt-decode";

interface AuthContextType {
  isRegistrationRequired: () => boolean;
  handleExternalAuth: (
    authCode: string,
    provider: string,
    redirectUri: string
  ) => Promise<void>;
  completeRegistration: (data: User) => Promise<void>;
  validateRedirection: () => void;
  activateAccount: (token: string) => Promise<void>;
  resendActivationEmail: (token: string) => Promise<void>;
  decodeToken: () => JwtPayload | null;
  login: (url: string) => void;
  logout: () => void;
  validateToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const navigate = useNavigate();

  const isRegistrationRequired = useCallback((): boolean => {
    const ls = localStorage.getItem("isRegistrationRequired");
    if (ls) {
      return ls as unknown as boolean;
    }
    return false;
  }, []);

  const login = useCallback((url: string) => {
    clearLocalStorage();
    window.location.href = url;
  }, []);

  const decodeToken = useCallback((): JwtPayload | null => {
    try {
      const user = getUser();
      if (user) {
        const decoded = jwtDecode<JwtPayload>(user.access_token);
        return decoded;
      }
      throw "user is null";
    } catch (error) {
      console.error("Error al decodificar el token:", error);
      return null;
    }
  }, []);

  const validateToken = useCallback(async () => {
    try {
      await AuthService.validateToken();
    } catch (error) {
      if (error instanceof RequestError) {
        if (error.getStatus() === HttpStatusCode.Unauthorized) {
          clearLocalStorage();
          navigate("/login");
          return;
        }
      }
      throw error;
    }
  }, []);

  const checkExistingUser = useCallback(
    async (cuit: string, provider: string) => {
      const userExists = await AuthService.checkExistingUser(cuit, provider);
      if (userExists.success) {
        localStorage.setItem("isRegistrationRequired", "false");
        const { data: info } = userExists;
        if (!info.email_validated) {
          clearLocalStorage();
          navigate("/email-verification-sent", {
            state: { email: info.email, type: 2 },
          });
          return;
        }

        const result = await AuthService.loginWithAuthApi(cuit, provider);
        if (result.success) {
          const { data } = result;
          const loginInfo: Token = {
            access_token: data.access_token,
            data: info,
          };
          setLocalStorage(loginInfo);
        }

        if (userExists.admin) {
          navigate("/admin-panel");
          return;
        }
        navigate("/admin/requests");
      } else {
        localStorage.setItem("isRegistrationRequired", "true");
        navigate("/complete-registration");
      }
    },
    [navigate]
  );

  const handleExternalAuth = useCallback(
    async (authCode: string, provider: string, redirectUri: string) => {
      const executionKey = `auth-executed-${provider}-${authCode}`;
      if (sessionStorage.getItem(executionKey)) {
        return;
      }
      sessionStorage.setItem(executionKey, "true");

      try {
        if (getUser()) {
          return;
        }

        const authData = await AuthService.getExternalUserInfo(
          authCode,
          provider,
          redirectUri
        );
        await checkExistingUser(authData.cuit, provider);
      } catch (error) {
        clearLocalStorage();
        console.error("Error en autenticación:", error);
        throw error;
      } finally {
        sessionStorage.removeItem(executionKey);
      }
    },
    [checkExistingUser]
  );

  const completeRegistration = useCallback(
    async (registrationData: User) => {
      try {
        await AuthService.registerUser(registrationData);
        clearLocalStorage();
      } catch (error) {
        if (error instanceof RequestError) {
          if (error.getStatus() === HttpStatusCode.Unauthorized) {
            alert("Se perdió la sesion de usuario. Ingrese nuevamente");
            clearLocalStorage();
            navigate("/login");
          }
        }
        throw error;
      }
    },
    [navigate]
  );

  const activateAccount = useCallback(async (token: string) => {
    await AuthService.activateAccount(token);
  }, []);

  const resendActivationEmail = useCallback(async (token: string) => {
    await AuthService.resendActivationEmail(token);
  }, []);

  // call this function to sign out logged in user
  const logout = useCallback(async () => {
    try {
      const provider = localStorage.getItem("provider") || "";
      const tokenInfo = getAutenticarInfo();
      if (tokenInfo === null) {
        throw "invalid info";
      }
      await AuthService.externalLogout(provider, tokenInfo.refresh_token);
    } catch (error) {
      console.error("logout error:", error);
    } finally {
      clearLocalStorage();
      navigate("/login");
    }
  }, [navigate]);

  const validateRedirection = useCallback(() => {
    if (!getUser()) {
      navigate("/login");
      return;
    }

    const ls = localStorage.getItem("isRegistrationRequired");
    if (ls) {
      const requireRegistration = ls as unknown as boolean;
      if (requireRegistration === true) {
        navigate("/complete-registration");
      }
    }
  }, [navigate]);

  const value = useMemo(
    () => ({
      isRegistrationRequired,
      handleExternalAuth,
      completeRegistration,
      validateRedirection,
      activateAccount,
      resendActivationEmail,
      decodeToken,
      login,
      logout,
      validateToken,
    }),
    [
      isRegistrationRequired,
      handleExternalAuth,
      completeRegistration,
      validateRedirection,
      activateAccount,
      resendActivationEmail,
      decodeToken,
      login,
      logout,
      validateToken,
    ]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe ser usado dentro de un AuthProvider");
  }
  return context;
};
