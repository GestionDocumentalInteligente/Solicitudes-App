import { Alert, TextInput } from "flowbite-react";
import { HiInformationCircle } from "react-icons/hi";
import { MdCheckCircle } from "react-icons/md";

import Info from "../../../../components/Info/Info";
import { Label } from "../../../../components/Label/Label";
import StyledTitle from "../../../../components/Title/StyledTitle";
import CustomButton from "../../../../components/Buttons/CustomButton";

import { RequestService } from "../requestsService";

import { useEffect, useState } from "react";
import { Address } from "../types";

interface StepOneProps {
  address: Address | null;
  commonZone: boolean;
  ablDebt: boolean | null;
  setAblDebt: (value: boolean | null) => void;
  setDisabled: (value: boolean) => void;
  onAddressChange: (value: Address | null) => void;
  setCommonZone: (value: boolean) => void;
}

export const StepOne = ({
  address,
  commonZone,
  onAddressChange,
  ablDebt,
  setAblDebt,
  setCommonZone,
  setDisabled,
}: StepOneProps) => {
  const [query, setQuery] = useState<string>("");
  const [suggestions, setSuggestions] = useState<Address[]>([]);
  const [showSuggestions, setShowSuggestions] = useState<boolean>(false);
  const [ablChecking, setAblChecking] = useState(false);

  const handleCheckABLDebt = async () => {
    if (!address) {
      alert("debe completar la direccion");
      return;
    }

    setAblChecking(true);

    try {
      const result = await RequestService.checkABL(address.abl_number);
      if (result.success) {
        setAblDebt(result.data.dbt);
        setAblChecking(false);
      }
    } catch (error) {
      setAblDebt(null);
      setAblChecking(false);
      alert("Error en la consulta de ABL. Por favor intente mas tarde.");
    }
  };

  useEffect(() => {
    sessionStorage.removeItem("ablDebt");
    if (ablDebt !== null) {
      sessionStorage.setItem("ablDebt", ablDebt.toString());
    }
  }, [ablDebt]);

  const handleInputChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.toLowerCase();
    setQuery(value);
    onAddressChange(null);
    setAblDebt(null);

    if (value.length >= 3) {
      try {
        const result = await RequestService.getAddressInfo(value);
        const suggestionsArray = result.data.suggestions || [];
        setSuggestions(suggestionsArray);
        setShowSuggestions(true);
      } catch (error) {
        setSuggestions([]);
        setShowSuggestions(false);
      }
    } else {
      setSuggestions([]);
      setShowSuggestions(false);
    }
  };

  useEffect(() => {
    if (address) {
      setQuery(
        `ABL: ${address.abl_number} - ${address.address_street} ${address.address_number}`
      );
    }
  }, [address]);

  const handleSuggestionClick = (address: Address) => {
    onAddressChange(address);
    setShowSuggestions(false);
  };

  useEffect(() => {
    if (address) {
      if (commonZone || (ablDebt !== null && !ablDebt)) {
        sessionStorage.setItem("address", JSON.stringify(address));
        setDisabled(false);
        return;
      }
      setDisabled(true);
    }
  }, [address, commonZone, ablDebt]);

  const handleSetCommonZone = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setAblDebt(null);
    }

    setCommonZone(e.target.checked);
    sessionStorage.setItem("commonZone", String(e.target.checked));
  };

  return (
    <>
      <StyledTitle title="Domicilio (donde se realizará la obra)" />
      <Info
        text="Recordá que solo podés realizar esta trámite si sos propietario,
        inquilino o administrador de los espacios comunes."
      />
      <Info
        text="Si el domicilio es un edificio o PH,
          seleccioná el número de ABL asociado a la Unidad Funcional."
      />
      <div className="space-y-2 mt-4">
        <Label text="Domicilio y n° de ABL*" mandatory={true} />
        <TextInput
          id="address"
          type="address"
          value={query}
          onChange={handleInputChange}
          placeholder="Ingresá dirección"
        />
        {showSuggestions && (
          <ul className="w-full mt-1 bg-white border rounded-lg shadow-md">
            {suggestions.length > 0 ? (
              suggestions.map((address, index) => (
                <li
                  key={index}
                  onClick={() => handleSuggestionClick(address)}
                  className="px-4 py-2 cursor-pointer hover:bg-gray-200"
                >
                  {`ABL:${address.abl_number} - ${address.address_street} ${address.address_number}`}
                </li>
              ))
            ) : (
              <li className="px-4 py-2 text-gray-500">No hay sugerencias</li>
            )}
          </ul>
        )}
        <div className="flex items-center mb-4 mt-6">
          <input
            id="default-checkbox"
            type="checkbox"
            checked={commonZone}
            onChange={handleSetCommonZone}
            className="w-4 h-4 text-green-700 bg-gray-100 border-gray-300 rounded focus:ring-green-600 focus:ring-2"
          />
          <label
            form="default-checkbox"
            className="ms-2 text-sm font-medium text-gray-900"
          >
            Marcá si la obra se realizará en un área de uso compartido (zona
            común).
          </label>
        </div>
        <br />
        <div className="w-full mt-4">
          <CustomButton
            onClick={handleCheckABLDebt}
            className="w-full mt-4 bg-green-700 text-white ml-1"
            isLoading={ablChecking}
            disabled={ablChecking || commonZone || address === null}
          >
            Consultar deuda ABL
          </CustomButton>
        </div>
        {ablDebt !== null && (
          <Alert
            color={ablDebt === true ? "failure" : "success"}
            icon={ablDebt ? HiInformationCircle : MdCheckCircle}
          >
            {ablDebt
              ? "La unidad tiene deuda de ABL. Para continuar, es necesario abonarla conforme Art. 45 de la Ordenanza Fiscal de S.I."
              : "¡Genial! La unidad no tiene deuda de ABL. ¡Continuemos!"}
          </Alert>
        )}
      </div>
    </>
  );
};
