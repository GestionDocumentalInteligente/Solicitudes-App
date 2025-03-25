import { Router } from "express";

import auth from "./auth";
import categories from "./categories";
import users from "./users";
import address from "./address";
import login from "./login";
import verifications from "./requestVerification";
import validations from "./requestValidation";
import activities from "./activities";
import requests from "./requests";

const router: Router = Router();

router.get("/ping", (req, res) => {
  res.status(200).json({ message: "UI says Pong!" });
});

router.use("/auth", auth);

router.use(
  "/admin",
  categories,
  users,
  login,
  address,
  activities,
  verifications,
  validations,
  requests
);

export default router;
