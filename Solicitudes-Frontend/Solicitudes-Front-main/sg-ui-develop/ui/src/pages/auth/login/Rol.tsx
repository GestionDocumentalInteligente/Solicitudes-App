import React from "react";
import { useNavigate } from "react-router-dom";
import { AuthProvider } from "../provider/useAuth";
import CustomButton from "@/components/Buttons/CustomButton";

const Rol: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col">
      <div className="flex flex-col mx-auto md:py-0">
        <div className="md:space-y-4 sm:p-4">
          <div className="space-y-4 md:space-y-4 max-w-md">
            <p className="font-sans text-base">
              ¿Desde que rol queres ingresar?
            </p>
            <CustomButton
              onClick={() => navigate("/admin/requests")}
              className="bg-title w-full"
            >
              Vecino
            </CustomButton>
            <CustomButton
              onClick={() => navigate("/admin-panel/request-verification")}
              className="bg-title w-full"
            >
              Empleado Municipal
            </CustomButton>
            <CustomButton
              onClick={() => navigate("/admin-panel/request-validation")}
              className="bg-title w-full"
            >
              Funcionario
            </CustomButton>
          </div>
        </div>
        <p className="font-sans max-w-lg text-xs text-center w-full">
          Al usar este servicio, aceptás nuestros{" "}
          <a
            href="https://docs.google.com/document/d/1U33075GjuApU84zV5v5wTM5W1-dtD8y-ijj15i-ZQE4/preview"
            className="underline"
            target="_blank"
          >
            términos y condiciones.
          </a>
        </p>
      </div>
    </div>
  );
};

const RolPage = () => (
  <AuthProvider>
    <Rol />
  </AuthProvider>
);

export default RolPage;
