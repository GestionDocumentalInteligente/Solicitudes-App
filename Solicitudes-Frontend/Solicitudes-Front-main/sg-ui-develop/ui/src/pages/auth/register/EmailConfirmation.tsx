import { useLocation, useNavigate } from "react-router-dom";
import { useEffect } from "react";
import { getUser } from "../provider/useLocalStorage";
import EmailSuccessMessage from "../../../components/EmailSuccessMessage/EmailSuccessMessage";

const EmailConfirmation = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const { email, type } = location.state || {};
  var title =
    "¡Listo! Te enviamos un e-mail a la dirección que proporcionaste.";
  var content =
    "Chequeá tu bandeja de entrada y hacé clic en el enlace de verificación. Si no lo ves en los próximos minutos, revisá en spam.";

  if (type === 2) {
    title = "El usuario se enuentra registrado pero aun no activo la cuenta";
    content =
      "Chequeá tu bandeja de entrada y hacé clic en el enlace de verificación.";
  }

  useEffect(() => {
    const token = getUser();
    if (token) {
      navigate("/admin/requests");
    }
  }, [navigate]);

  return <EmailSuccessMessage email={email} title={title} content={content} />;
};

export default EmailConfirmation;
