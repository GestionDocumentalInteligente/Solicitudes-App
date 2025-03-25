import React, { useEffect, useState } from "react";
import { CiCircleCheck } from "react-icons/ci";
import { MdOutlineErrorOutline } from "react-icons/md";

import { useAuth } from "../provider/useAuth";
import { RequestError } from "../types";
import { HttpStatusCode } from "axios";
import { Link, useLocation } from "react-router-dom";
import CustomButton from "../../../components/Buttons/CustomButton";
import LoadingScreen from "@/components/LoadingScreen/LoadingScreen";

// Posibles estados del proceso de activación
type Status = "loading" | "success" | "error" | "resend";

const ActivateAccount: React.FC = () => {
  const auth = useAuth();
  const location = useLocation();
  const [status, setStatus] = useState<Status>("loading");
  const [message, setMessage] = useState<string>("");
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const urlParams = new URLSearchParams(location.search);
    const token = urlParams.get("token");
    setToken(token);
  }, [location.search]);

  useEffect(() => {
    if (token) {
      const validateToken = async () => {
        try {
          await auth.activateAccount(token);
          setStatus("success");
          setMessage("Verificación exitosa");
        } catch (error) {
          setStatus("error");
          if (error instanceof RequestError) {
            setMessage(error.message);
            return;
          }
          setMessage("Error en la solicitud de correo de activación.");
        }
      };
      validateToken();
    }
  }, [token, auth]);

  const handleResendEmail = async () => {
    if (token) {
      setStatus("loading");
      try {
        await auth.resendActivationEmail(token);
        setStatus("resend");
        setMessage(
          "Se envió un nuevo correo electrónico. Por favor verifique su casilla para activar la cuenta."
        );
      } catch (error) {
        setStatus("error");
        if (error instanceof RequestError) {
          if (error.getStatus() === HttpStatusCode.BadRequest) {
            setStatus("success");
          }
          setMessage(error.message);
          return;
        }
        setMessage("Error en la solicitud de correo de activación.");
      }
    }
  };

  if (status === "loading") {
    return <LoadingScreen title={["Validando datos de usuario"]} />;
  }

  return (
    <div className="mx-auto">
      <div className="flex items-center justify-center mt-4">
        {(status === "success" || status === "resend") && (
          <CiCircleCheck className="w-24 h-24 text-gray-600" />
        )}
        {status === "error" && (
          <MdOutlineErrorOutline className="w-24 h-24 text-gray-600" />
        )}
      </div>
      <h2 className="text-2xl font-bold text-gray-900 text-center">
        {message}
      </h2>
      <br />
      {status === "success" && (
        <>
          <p className="text-sm text-gray-700 text-center">
            Ahora podés comenzar a gestionar tus solicitudes
            <br />
            de forma rapida y simple.
          </p>
          <br />
          <div className="flex justify-center">
            <Link to="/login">
              <CustomButton type="button" className="bg-title w-full">
                Empezar
              </CustomButton>
            </Link>
          </div>
        </>
      )}
      {status === "error" && (
        <p className="text-sm text-gray-800 text-center font-bold">
          <a href="#" onClick={handleResendEmail}>
            <span className="text-blue-700">solicitar un nuevo envío.</span>
          </a>
        </p>
      )}
    </div>
  );
};

export default ActivateAccount;
