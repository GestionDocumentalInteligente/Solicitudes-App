import { Request, Response, Router } from "express";
import axios from "axios";
import { configService } from "../configService";
import { ApiClient } from "../clients/ApiClient";

const router: Router = Router();

const apiClient = new ApiClient(configService.baseLoginApi);

export enum Provider {
  AFIP = "AFIP",
  ANSES = "ANSES",
  MI_ARGENTINA = "MI_ARGENTINA",
}

router.post("/login", async (req: Request, res: Response) => {
  const { code, provider, redirectUri } = req.body;
  if (!code || !isValidProvider(provider)) {
    res.status(400).json({ error: "Code is required" });
    return;
  }

  let realmId = "";
  let secret = "";

  switch (provider) {
    case Provider.AFIP:
      realmId = configService.realmAfip;
      secret = configService.clientSecretAfip;
      break;
    case Provider.ANSES:
      realmId = configService.realmAnses;
      secret = configService.clientSecretAnses;
      break;
    case Provider.MI_ARGENTINA:
      realmId = configService.realmMiArg;
      secret = configService.clientSecretMiArg;
      break;
    default:
      break;
  }

  try {
    const response = await axios.post(
      `${configService.authUrl}/auth/realms/${realmId}/protocol/openid-connect/token`,
      new URLSearchParams({
        grant_type: "authorization_code",
        client_id: configService.clientId,
        code: code,
        redirect_uri: `${redirectUri}?provider=${provider}`,
        client_secret: secret,
      }),
      {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      }
    );

    res.status(200).json(response.data);
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Error de Axios:", error);
    } else {
      console.error("Error inesperado:", error);
    }
    res.status(500).json({ error: "Failed to refresh token" });
  }
});

router.post("/logout", async (req: Request, res: Response) => {
  const { provider } = req.body;
  if (!isValidProvider(provider)) {
    res.status(400).json({ error: "Code is required" });
    return;
  }

  let realmId = "";
  let secret = "";

  switch (provider) {
    case Provider.AFIP:
      realmId = configService.realmAfip;
      secret = configService.clientSecretAfip;
      break;
    case Provider.ANSES:
      realmId = configService.realmAnses;
      secret = configService.clientSecretAnses;
      break;
    case Provider.MI_ARGENTINA:
      realmId = configService.realmMiArg;
      secret = configService.clientSecretMiArg;
      break;
    default:
      break;
  }

  const data = {
    url: configService.authUrl,
    realmId: realmId,
    clientId: configService.clientId,
    clientSecret: secret,
  };

  res.status(200).json(data);
});

router.get("/validate", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const data = await apiClient.get<any>("/auth/protected/hi", headers);
    res.status(200).json({ success: true, data });
  } catch (error: any) {
    console.error(error);
    res
      .status(error.status || 500)
      .json(
        error.data
          ? { success: false, error }
          : { success: false, message: "Error inesperado" }
      );
  }
});

function isValidProvider(provider: any): provider is Provider {
  return Object.values(Provider).includes(provider);
}

export default router;
