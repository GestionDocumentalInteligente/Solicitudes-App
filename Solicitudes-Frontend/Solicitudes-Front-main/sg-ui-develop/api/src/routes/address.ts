import { Request, Response, Router } from "express";
import { ApiClient } from "../clients/ApiClient";
import { configService } from "../configService";

const apiClient = new ApiClient(configService.baseRequestApi);
const fileManagerClient = new ApiClient(configService.baseFileManagerApi);

const router: Router = Router();

// const mockAddresses = [
//   {
//     abl_number: 987654,
//     address_street: "HEROES DE LAS MALVINAS",
//     address_number: 0,
//     property_id: 1,
//   },
//   {
//     property_id: 2,
//     address_street: "HEROES DE LAS MALVINAS",
//     address_number: 357,
//     abl_number: 128663,
//   },
//   {
//     property_id: 4,
//     address_street: "SAN MARTIN",
//     address_number: 123,
//     abl_number: 456789,
//   },
//   {
//     property_id: 3,
//     address_street: "AV. BELGRANO",
//     address_number: 50,
//     abl_number: 987654,
//   },
// ];

router.get("/address", async (req: Request, res: Response) => {
  try {
    const query = (req.query.q as string)?.toLowerCase();

    if (!query) {
      res.status(400).json({ message: "Query no proporcionado" });
      return;
    }

    // const filteredAddresses = mockAddresses.filter((address) =>
    //   address.address_street.toLowerCase().includes(query)
    // );

    // res
    //   .status(200)
    //   .json({ success: true, data: { suggestions: filteredAddresses } });
    // return;

    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const data = await apiClient.get<any>(
      "/requests/address/autocomplete?q=" + query,
      headers
    );
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

router.get("/address/debt", async (req: Request, res: Response) => {
  try {
    const ablNumber = (req.query.abl as string)?.toLowerCase();

    if (!ablNumber) {
      res.status(400).json({ message: "ablNumber no proporcionado" });
      return;
    }

    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};

    const payload = {
      abl_number: ablNumber,
      type: 1,
    };

    const data = await fileManagerClient.post<any>(
      "/file-manager/address/debt",
      payload,
      headers
    );

    res.status(200).json({ success: true, data: data });
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

export default router;
