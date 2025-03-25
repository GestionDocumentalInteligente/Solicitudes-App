import { Label } from "flowbite-react";
import CustomButton from "../../../components/Buttons/CustomButton";
import { User } from "../types";

interface ProfileInfoProps {
  userInfo: User | null;
  error: string;
  setError: (value: string) => void;
  onNext: () => void;
}

export function ProfileInfo({
  userInfo,
  error,
  onNext,
  setError,
}: ProfileInfoProps) {
  const getDniFromCuil = (cuil: string): string => {
    if (cuil && cuil.length === 11) {
      return cuil.substring(2, 10);
    }
    return cuil;
  };

  const dataError = () => {
    setError(
      "Por favor, actualizá la información en AFIP, ANSES o Mi Argentina, según la plataforma que utilizaste para ingresar."
    );
  };

  return (
    <div className="h-full">
      <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <Label htmlFor="disabledCuil" value="CUIL" />
          <br />
          <span className="text-gray-700 text-sm">
            {userInfo ? userInfo.cuit : "xxxxxxxxxx"}
          </span>
        </div>
        <div>
          <Label htmlFor="disabledDni" value="DNI" />
          <br />
          <span className="text-gray-700 text-sm">
            {userInfo ? getDniFromCuil(userInfo.dni) : "xxxxxxxx"}
          </span>
        </div>
        <div>
          <Label htmlFor="disabledName" value="Nombre/s" />
          <br />
          <span className="text-gray-700 text-sm">
            {userInfo ? userInfo.first_name : "xxxxxxxxxx"}
          </span>
        </div>
        <div>
          <Label htmlFor="disabledSurname" value="Apellido/s" />
          <br />
          <span className="text-gray-700 text-sm">
            {userInfo ? userInfo.last_name : "xxxxxxxxxx"}
          </span>
        </div>
      </div>
      <div className="w-full flex flex-col items-center mt-10 justify-center">
        <CustomButton
          disabled={error !== ""}
          className="bg-title w-full"
          onClick={onNext}
        >
          Continuar
        </CustomButton>
        <a
          href="#"
          onClick={dataError}
          className="text-sm text-black underline mt-4"
        >
          Hay un error con mis datos
        </a>
      </div>
    </div>
  );
}
