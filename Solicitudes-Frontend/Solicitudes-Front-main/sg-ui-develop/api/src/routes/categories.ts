import { Request, Response, Router } from "express";
import { ApiClient } from "../clients/ApiClient";
import { configService } from "../configService";

const apiClient = new ApiClient(configService.baseRequestApi);

const router: Router = Router();

router.get("/categories", async (req: Request, res: Response) => {
  try {
    const data = await apiClient.get<any>("/categories");
    res.status(200).json({ success: true, data });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      success: false,
      error: "Error getting categories",
    });
  }
});

router.get("/categories/:id/requests", async (req: Request, res: Response) => {
  const categoryId = req.params.id;

  try {
    const data = await apiClient.get<any>(
      "/categories/" + categoryId + "/requests"
    );

    res.status(200).json({ success: true, data });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      success: false,
      error: "Error getting requests",
    });
  }
});

export default router;
