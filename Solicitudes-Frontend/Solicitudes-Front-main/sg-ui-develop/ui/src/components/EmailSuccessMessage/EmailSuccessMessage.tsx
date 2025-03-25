import { HiMail } from "react-icons/hi";

import CustomButton from "../Buttons/CustomButton";
import { Link, useNavigate } from "react-router-dom";

interface EmailMessageProps {
  email: string;
  title: string;
  content: string;
}

const EmailSuccessMessage = ({ email, title, content }: EmailMessageProps) => {
  const navigate = useNavigate();

  const handleResendEmail = async () => {
    try {
      const response = await fetch("/api/admin/users/request-new-activation", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: email }),
      });
      if (response.ok) {
        alert(
          "Se envió nuevamente el correo de activación. Por favor verifique su bandeja de entrada."
        );
        navigate("/email-verification-sent", {
          state: { email: email, type: 1 },
        });
      } else {
        alert("Hubo un problema al reenviar el correo. Intenta nuevamente.");
      }
    } catch (error) {
      console.error("Error al reenviar el correo:", error);
      alert("Error al procesar la solicitud.");
    }
  };

  return (
    <div className="mx-auto">
      <div className="flex items-center justify-center">
        <HiMail className="w-24 h-24 text-gray-600" />
      </div>
      <h2 className="text-2xl font-bold text-gray-900 text-center">{title}</h2>
      <br />
      <p className="text-sm text-gray-700 text-center">{content}</p>
      <p className="text-sm text-gray-800 text-center font-bold">
        <a
          href="#"
          onClick={(e) => {
            e.preventDefault();
            handleResendEmail();
          }}
        >
          ¿No te llegó?{" "}
          <span className="text-blue-700">solicitar un nuevo envío.</span>
        </a>
      </p>
      <br />
      <div className="flex justify-center">
        <Link to="/login">
          <CustomButton type="button" className="bg-title w-full">
            Volver al inicio de Sesión
          </CustomButton>
        </Link>
      </div>
    </div>
  );
};

export default EmailSuccessMessage;
