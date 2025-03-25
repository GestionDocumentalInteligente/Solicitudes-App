import { Request, Response, Router } from "express";
import { ApiClient } from "../clients/ApiClient";
import { configService } from "../configService";

const apiClient = new ApiClient(configService.baseRequestApi);

const router: Router = Router();

router.get("/requests/verifications", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const data = await apiClient.get<any>(
      "/requests/protected/verifications",
      headers
    );
    res.status(200).json({ success: true, data });
    // res.status(200).json({
    //     success: true, data: [
    //         {
    //             "id": 9,
    //             "recordNumber": "EX-2024-00038483-   -SI-MDSI",
    //             "requestType": "Aviso de obra",
    //             "documentType": "Postestad del inmueble",
    //             "deliveryDate": "2024-12-01 13:44:24",
    //             "status": "Pending",
    //             "requesterFullName": "LESLIE ANN CHRISTINE",
    //             "requesterCuil": "20348134664",
    //             "requesterAddress": " 0, ",
    //             "documents": null
    //         },
    //         {
    //             "id": 11,
    //             "recordNumber": "code",
    //             "requestType": "Aviso de obra",
    //             "documentType": "Postestad del inmueble",
    //             "deliveryDate": "2024-12-03 11:49:03",
    //             "status": "Pending",
    //             "requesterFullName": "LESLIE ANN CHRISTINE",
    //             "requesterCuil": "20348134664",
    //             "requesterAddress": " 0, ",
    //             "documents": null
    //         }
    //     ],
    // });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      success: false,
      error: "Error getting requests verification",
    });
  }
});

router.get(
  "/requests/verifications/:recordNumber",
  async (req: Request, res: Response) => {
    try {
      const authHeader = req.headers.authorization;
      const headers = authHeader ? { Authorization: authHeader } : {};

      const { recordNumber } = req.params;
      const data = await apiClient.get<any>(
        `requests/protected/documents?recordNumber=${recordNumber}`,
        headers
      );

      res.status(200).json({ success: true, data });
      // res.status(200).json({
      //     success: true, data: {
      //         "documents": [
      //             {
      //                 "id": "177",
      //                 "title": "Reglamento de Co-propiedad",
      //                 "gedoCode": "IF-2024-00039883-SI-MDSI"
      //             },
      //             {
      //                 "id": "178",
      //                 "title": "",
      //                 "gedoCode": "IF-2024-00039886-SI-MDSI"
      //             },
      //         ]
      //     }
      // });
    } catch (error) {
      console.error(error);

      res.status(500).json({
        success: false,
        error: "Error getting requests verification",
      });
    }
  }
);

router.put("/requests/verifications/:id", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const { id } = req.params;
    const body = req.body;

    const data = await apiClient.put<any>(`requests/protected/verification/${id}`, body, headers);
    res.status(200).json({ success: true, data });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      success: false,
      error: "Error getting requests verification",
    });
  }
});

export default router;
