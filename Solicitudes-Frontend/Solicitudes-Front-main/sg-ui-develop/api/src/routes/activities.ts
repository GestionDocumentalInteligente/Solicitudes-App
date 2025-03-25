import { Request, Response, Router } from "express";

import { ApiClient } from "../clients/ApiClient";
import { configService } from "../configService";

const apiClient = new ApiClient(configService.baseRequestApi);

const router: Router = Router();

const mockActivities = [
  {
    id: 1,
    description: "Ejecutar solados (piso)",
  },
  {
    id: 2,
    description: "Cambiar revestimientos",
  },
  {
    id: 3,
    description: "Terraplenar y rellenar terrenos",
  },
  {
    id: 4,
    description: "Cambiar el material de cubierta de techos",
  },
  {
    id: 5,
    description: "Ejecutar cielorrasos",
  },
  {
    id: 6,
    description: "Revocar cercas al frente",
  },
  {
    id: 7,
    description: "Ejecutar revoques exteriores o trabajos similares",
  },
  {
    id: 8,
    description: "Limpiar o pintar las fachadas principales",
  },
];

router.get("/site/activities", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    //const data = await apiClient.get<any>("/site/activities", headers);
    res.status(200).json({ success: true, data: mockActivities });
  } catch (error: any) {
    console.error(error);
    res
      .status(error.status || 500)
      .json(error.data || { message: "Error inesperado" });
  }
});

export default router;
