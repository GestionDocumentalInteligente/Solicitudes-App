import { CiCircleCheck } from "react-icons/ci";

import { Link, useLocation } from "react-router-dom";
import CustomButton from "../../../../components/Buttons/CustomButton";

const RequestSuccess = () => {
  const location = useLocation();

  const { updated } = location.state || null;

  return (
    <div className="mx-auto">
      <div className="flex items-center justify-center">
        <CiCircleCheck className="w-24 h-24 text-gray-600" />
      </div>

      <h2 className="text-2xl font-bold text-gray-900 text-center">
        ¡Listo! <br />
        {updated === null
          ? "Tu solicitud se realizó con éxito"
          : "Tu actualización se realizó con éxito"}
      </h2>
      <br />
      {updated === null ? (
        <p className="text-sm text-gray-700 text-center">
          En minutos recibirás tu número de solicitud por e-mail.
          <br />
          La validación puede demorar hasta 48 horas hábiles.
        </p>
      ) : (
        <p className="text-sm text-gray-700 text-center">
          Estamos verficando los datos que enviaste.
          <br />
          Este proceso puede demorar hasta 48 horas hábiles.
          <br />
          Recibiras un e-mail cuando la validación este completa.
        </p>
      )}
      <br />
      <div className="flex justify-center">
        <Link to="/admin/requests">
          <CustomButton type="button" className="bg-title w-full">
            Aceptar
          </CustomButton>
        </Link>
      </div>
    </div>
  );
};

export default RequestSuccess;
