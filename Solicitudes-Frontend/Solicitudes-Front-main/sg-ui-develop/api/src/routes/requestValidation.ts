import { Request, Response, Router } from "express";
import { ApiClient } from "../clients/ApiClient";
import { configService } from "../configService";

const apiClient = new ApiClient(configService.baseRequestApi);

const router: Router = Router();

router.get("/requests/validations", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const data = await apiClient.get<any>(
      "/requests/protected/validations",
      headers
    );
    res.status(200).json({ success: true, data });
    // res.status(200).json({
    //     success: true, data: [
    //         {
    //             "recordNumber": "2024-000030",
    //             "requestType": "Aviso de obra",
    //             "documentType": "Potestad sobre el inmueble",
    //             "deliveryDate": "01/10/2024",
    //             "status": "pending",
    //             "requesterFullName": "Pablo Perez",
    //             "requesterCuil": "20354352617",
    //             "requesterAddress": "Ramón Falcón 2021, UF4, Beccar",
    //         },
    //         {
    //             "recordNumber": "2024-000031",
    //             "requestType": "Aviso de obra",
    //             "documentType": "Potestad sobre el inmueble",
    //             "deliveryDate": "02/10/2024",
    //             "status": "validated",
    //             "requesterFullName": "Pablo Perez",
    //             "requesterCuil": "20354352617",
    //             "requesterAddress": "Ramón Falcón 2021, UF4, Beccar",
    //         },
    //         {
    //             "recordNumber": "2024-000031",
    //             "requestType": "Aviso de obra",
    //             "documentType": "Potestad sobre el inmueble",
    //             "deliveryDate": "02/10/2024",
    //             "status": "rejected",
    //             "requesterFullName": "Pablo Perez",
    //             "requesterCuil": "20354352617",
    //             "requesterAddress": "Ramón Falcón 2021, UF4, Beccar",
    //         }
    //     ],
    // });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      success: false,
      error: "Error getting requests validation",
    });
  }
});

router.get(
  "/requests/validations/:recordNumber",
  async (req: Request, res: Response) => {
    try {
      const authHeader = req.headers.authorization;
      const headers = authHeader ? { Authorization: authHeader } : {};

      const { recordNumber } = req.params;
      const data = await apiClient.get<any>(
        `requests/protected/validations/documents?recordNumber=${recordNumber}`,
        headers
      );

      res.status(200).json({ success: true, data });
      //   res.status(200).json({
      //     success: true,
      //     data: {
      //       documents: [
      //         {
      //           id: "1",
      //           title: "Potestad sobre el inmueble - Expediente",
      //           gedoCode: "PV-2024-00015011-SI-TESTSADE",
      //         },
      //         {
      //           id: "2",
      //           title: "Potestad sobre el inmueble - Verificación",
      //           gedoCode: "PV-2024-00015012-SI-TESTSADE",
      //           verifiedBy: "Marta Lopez",
      //           verifiedDate: "01/10/2024",
      //         },
      //         {
      //           id: "3",
      //           title: "Tareas a realizar - Expediente",
      //           gedoCode: "PV-2024-00015013-SI-TESTSADE",
      //         },
      //         {
      //           id: "4",
      //           title: "Tareas a realizar - Verificación",
      //           gedoCode: "PV-2024-00015014-SI-TESTSADE",
      //           verifiedBy: "Marta Lopez",
      //           verifiedDate: "01/10/2024",
      //         },
      //       ],
      //     },
      //   });
    } catch (error) {
      console.error(error);

      res.status(500).json({
        success: false,
        error: "Error getting requests validation",
      });
    }
  }
);

router.put("/requests/validations/:id", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const { id } = req.params;
    const body = req.body;

    const data = await apiClient.put<any>(
      `requests/protected/validation/${id}`,
      body,
      headers
    );
    res.status(200).json({ success: true, data });
  } catch (error) {
    res.status(500).json({
      success: false,
      error: "Error getting requests validation",
    });
  }
});

export default router;
