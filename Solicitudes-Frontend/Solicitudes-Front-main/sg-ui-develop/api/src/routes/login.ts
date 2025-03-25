import { Request, Response, Router } from "express";
import axios from "axios";
import { configService } from "../configService";
import { ApiClient } from "../clients/ApiClient";

const router: Router = Router();

const apiClient = new ApiClient(configService.baseLoginApi);

router.post("/login", async (req: Request, res: Response) => {
  const { cuit, provider } = req.body;
  if (!cuit) {
    res.status(400).json({ error: "cuit is required" });
    return;
  }

  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const data = await apiClient.get<any>(
      "/auth/protected/login?cuil=" + cuit + "&provider=" + provider,
      headers
    );
    res.status(200).json({ success: true, data });
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error("Error de Axios:", error);
    } else {
      console.error("Error inesperado:", error);
    }
    res.status(500).json({ error: "failed getting token" });
  }
});

export default router;
