import { Request, Response, Router } from "express";
import { ApiClient } from "../clients/ApiClient";
import { configService } from "../configService";

const apiClient = new ApiClient(configService.baseUsersApi);
const mailingClient = new ApiClient(configService.baseMailingApi);

const router: Router = Router();

const cuilAdmins = [
  27929013722, 27384011077, 20316497633, 27378053221, 20350722220, 27357961152,
  20960093071, 20421461873, 27183233179, 20324781340, 27255607931, 20346692457,
  20348134664, 20378552797,
];

router.post("/users", async (req: Request, res: Response) => {
  try {
    const authHeader = req.headers.authorization;
    const headers = authHeader ? { Authorization: authHeader } : {};
    const provider = req.query.provider as string;

    const data = await apiClient.post<any>(
      "/users?provider=" + provider,
      req.body,
      headers
    );

    const name = req.body.first_name + " " + req.body.last_name;

    void sendWelcomeEmail(provider, req.body.email, name, headers);

    res.status(201).json({ success: true, data });
  } catch (error: any) {
    console.error(error);
    res
      .status(error.status || 500)
      .json(error.data || { message: "Error inesperado" });
  }
});

async function sendWelcomeEmail(
  provider: string,
  email: string,
  name: string,
  headers?: Record<string, string | undefined>
) {
  try {
    await mailingClient.post(
      "/protected/email?provider=" + provider,
      { email: email, name: name },
      headers
    );
  } catch (error) {
    console.error("Error sending welcome email:", error);
  }
}

router.get("/users", async (req: Request, res: Response) => {
  try {
    const cuit = req.query.cuit as string;
    const provider = req.query.provider as string;

    if (cuit && provider) {
      const authHeader = req.headers.authorization;
      const headers = authHeader ? { Authorization: authHeader } : {};

      const data = await apiClient.get<any>(
        "/users?cuil=" + cuit + "&provider=" + provider,
        headers
      );
      res.status(200).json({ success: true, admin: isAnAdminCuil(cuit), data });
    } else {
      res.status(400).json({ message: "cuit or provider not provided" });
    }
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

router.get("/users/activate", async (req: Request, res: Response) => {
  try {
    const token = req.query.token as string;

    const data = await mailingClient.get<any>(
      "/activate-account?token=" + token
    );
    res.status(200).json({ success: true, data });
  } catch (error: any) {
    console.error(error);
    res
      .status(error.status || 500)
      .json(error.data || { message: "Error inesperado" });
  }
});

router.post(
  "/users/resend-activation-email",
  async (req: Request, res: Response) => {
    try {
      const { token } = req.body;

      const data = await mailingClient.post<any>("/resend-activation-email", {
        token,
      });
      res.status(200).json({ success: true, data });
    } catch (error: any) {
      console.error(error);
      res
        .status(error.status || 500)
        .json(error.data || { message: "Error inesperado" });
    }
  }
);

router.post(
  "/users/request-new-activation",
  async (req: Request, res: Response) => {
    try {
      const { email } = req.body;
      if (email === "") {
        res.status(400).json({ message: "email not provided" });
        return;
      }

      const data = await mailingClient.post<any>("/request-new-activation", {
        email,
      });
      res.status(200).json({ success: true, data });
    } catch (error: any) {
      console.error(error);
      res
        .status(error.status || 500)
        .json(error.data || { message: "Error inesperado" });
    }
  }
);

function isAnAdminCuil(cuit: string): boolean {
  const cuitNumber = Number(cuit);
  if (isNaN(cuitNumber)) {
    console.error("El CUIT no es un número válido", cuit);
    return false;
  }

  return cuilAdmins.includes(cuitNumber);
}

export default router;
