import { Button, Label } from "flowbite-react";
import Info from "../../../../components/Info/Info";
import { FaCloudDownloadAlt } from "react-icons/fa";

import StyledTitle from "../../../../components/Title/StyledTitle";
import { useEffect, useState } from "react";
import { RequestService } from "../requestsService";
import { Address, DocumentTypes, FilesType, User } from "../types";
import CustomButton from "@/components/Buttons/CustomButton";

interface StepResumeProps {
  address: Address | null;
  update: boolean;
  files: FilesType;
  selectedActivities: string[];
  estimatedTime: number;
  insurance: boolean | null;
  projectDescription: string;
  setDisabled: (value: boolean) => void;
  setActiveStep: (value: number) => void;
}

type Storage = {
  data: User;
};

export const StepResume = ({
  address,
  update,
  files,
  estimatedTime,
  insurance,
  projectDescription,
  setDisabled,
  setActiveStep,
}: StepResumeProps) => {
  const [terms, setTerms] = useState<boolean>(false);
  const [affidavit, setAffidavit] = useState<boolean>(false);
  const [user, setUser] = useState<User | null>(null);
  const [activity, setActivity] = useState<string>("");

  useEffect(() => {
    const getUser = async () => {
      const storedUser = localStorage.getItem("user");
      if (storedUser) {
        const storage = JSON.parse(storedUser) as Storage;
        const result = await RequestService.getUser(storage.data.cuil);
        setUser(result.data);
      }
    };
    const activityDescription = sessionStorage.getItem("activityDescription");
    setActivity(activityDescription ? activityDescription : "");
    getUser();
  }, []);

  const handleDownload = (fileName: string, fileContent: File) => {
    const url = URL.createObjectURL(fileContent);
    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    link.click();
    URL.revokeObjectURL(url);
  };

  useEffect(() => {
    if (
      address &&
      insurance !== null &&
      Object.keys(files).length > 0 &&
      estimatedTime > 0 &&
      projectDescription !== "" &&
      terms &&
      affidavit
    ) {
      setDisabled(false);
      return;
    }
    setDisabled(true);
  }, [
    files,
    address,
    estimatedTime,
    insurance,
    projectDescription,
    terms,
    affidavit,
  ]);

  return (
    <div>
      <StyledTitle title="Resumen" />
      <div className="w-full">
        <div className="mb-2 block">
          <Info text="Revisá que todos tus datos sean correctos y modificalos si ves algún error." />
        </div>
        <div className="mt-4 mb-4">
          <StyledTitle title="Datos personales solicitante" />
          <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="disabledCuil" value="CUIL" />
              <br />
              <span className="text-gray-700 text-sm">
                {user ? user.cuil : "xxxxxxxxxx"}
              </span>
            </div>
            <div>
              <Label htmlFor="disabledDni" value="DNI" />
              <br />
              <span className="text-gray-700 text-sm">
                {user ? user.dni : "xxxxxxxxxx"}
              </span>
            </div>
            <div>
              <Label htmlFor="disabledName" value="Nombre/s" />
              <br />
              <span className="text-gray-700 text-sm">
                {user ? user.first_name : "xxxxxxxxxx"}
              </span>
            </div>
            <div>
              <Label htmlFor="disabledSurname" value="Apellido/s" />
              <br />
              <span className="text-gray-700 text-sm">
                {user ? user.last_name : "xxxxxxxxxx"}
              </span>
            </div>
          </div>
        </div>
        <hr />
        <div className="mt-4 mb-4">
          <div className="flex items-center justify-between">
            <StyledTitle title="Domicilio" />
            <CustomButton
              disabled={update}
              onClick={() => setActiveStep(0)}
              className="mr-1"
              color="light"
            >
              Modificar datos
            </CustomButton>
          </div>
          <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="address" value="Dirección" />
              <br />
              <span className="text-gray-700 text-sm">
                {address
                  ? address.address_street + " " + address.address_number
                  : "xxxxxxxxxx"}
              </span>
            </div>
            <div>
              <Label value="Número ABL" />
              <br />
              <span className="text-gray-700 text-sm">
                {address ? address.abl_number : "xxxxxxxxxx"}
              </span>
            </div>
          </div>
        </div>
        <hr />
        <div className="mt-4 mb-4">
          <div className="flex items-center justify-between">
            <StyledTitle title="Potestad sobre el inmueble" />
            <CustomButton
              disabled={update}
              onClick={() => setActiveStep(1)}
              className="mr-1"
              color="light"
            >
              Modificar datos
            </CustomButton>
          </div>
          <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
            {Object.keys(files).length > 0 ? (
              Object.entries(files)
                .filter(([key]) => Number(key) !== DocumentTypes.Insurance.id)
                .map(([key, fileInfo]) => (
                  <div key={key} className="relative">
                    <input
                      key={key}
                      value={fileInfo.name}
                      type="text"
                      readOnly
                      className="col-span-6 block w-full rounded-lg border border-gray-300 bg-gray-50 px-2.5 py-4 text-sm text-gray-600"
                      disabled
                    />
                    <button
                      type="button"
                      onClick={() =>
                        handleDownload(fileInfo.name, fileInfo.content)
                      }
                      className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-600 hover:text-green-700"
                      aria-label="Download"
                    >
                      <FaCloudDownloadAlt size={20} />
                    </button>
                  </div>
                ))
            ) : (
              <>
                <div className="relative">
                  <span className="text-secondary-dark font-bold">
                    No existen documentos cargados
                  </span>
                </div>
                <div className="relative">
                  <Button
                    size="xs"
                    className="bg-primary-dark"
                    onClick={() => setActiveStep(1)}
                  >
                    Cargar documentos
                  </Button>
                </div>
              </>
            )}
          </div>
        </div>
        <hr />
        <div className="mt-4">
          <div className="flex items-center justify-between">
            <StyledTitle title="Tareas a realizar" />
            <CustomButton
              disabled={update}
              onClick={() => setActiveStep(2)}
              className="mr-1"
              color="light"
            >
              Modificar datos
            </CustomButton>
          </div>
          <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="task" value="Tarea a realizar" />
              <br />
              <span className="text-gray-700 text-sm">{activity}</span>
            </div>
            <div>
              <Label htmlFor="time" value="Tiempo estimado de ejecución" />
              <br />
              <span className="text-gray-700 text-sm">{estimatedTime}</span>
            </div>
          </div>
          <div className="mt-2 grid grid-cols-1">
            <div>
              <Label htmlFor="memory" value="Memoria descriptiva" />
              <br />
              <span className="text-gray-700 text-sm">
                {projectDescription}
              </span>
            </div>
          </div>
        </div>

        {insurance && (
          <div className="mt-4">
            <StyledTitle title="Seguro de accidente" />
            <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
              {DocumentTypes.Insurance.id in files ? (
                <div key={DocumentTypes.Insurance.id} className="relative">
                  <input
                    value={files[DocumentTypes.Insurance.id].name}
                    type="text"
                    readOnly
                    className="col-span-6 block w-full rounded-lg border border-gray-300 bg-gray-50 px-2.5 py-4 text-sm text-gray-600"
                    disabled
                  />
                  <button
                    type="button"
                    onClick={() =>
                      handleDownload(
                        files[DocumentTypes.Insurance.id].name,
                        files[DocumentTypes.Insurance.id].content
                      )
                    }
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-600 hover:text-green-700"
                    aria-label="Download"
                  >
                    <FaCloudDownloadAlt size={20} />
                  </button>
                </div>
              ) : (
                <>
                  <div className="relative">
                    <span className="text-secondary-dark font-bold">
                      No esta cargada la póliza de seguro
                    </span>
                  </div>
                  <div className="relative">
                    <Button
                      size="xs"
                      className="bg-primary-dark"
                      onClick={() => setActiveStep(2)}
                    >
                      Cargar documentos
                    </Button>
                  </div>
                </>
              )}
            </div>
          </div>
        )}
        <br />
        <hr />
        <Info
          text="Al firmar la solicitud, aceptás los Términos y Condiciones del
          servicio y manifestás, bajo declaración jurada, la veracidad de la
          información y documentación presentada, asumiendo total
          responsabilidad por su autenticidad y por el cumplimiento de las
          normativas vigentes."
        />
        <div className="flex items-center mb-2 mt-4">
          <input
            id="terms-checkbox"
            type="checkbox"
            checked={terms}
            onChange={(e) => setTerms(e.target.checked)}
            className="w-4 h-4 text-green-700 bg-gray-100 border-gray-300 rounded focus:ring-green-600 focus:ring-2"
          />
          <label
            form="terms-checkbox"
            className="ms-2 text-sm font-medium text-gray-900"
          >
            Acepto{" "}
            <a
              href="https://drive.google.com/file/d/1GEZXeTGKZD360HmgC3jsGhCQA9mcVyP8/preview"
              className="text-green-700 font-medium underline"
              target="_blank"
            >
              Términos y condiciones
            </a>
          </label>
        </div>
        <div>
          <input
            id="affidavit-checkbox"
            type="checkbox"
            checked={affidavit}
            onChange={(e) => setAffidavit(e.target.checked)}
            className="w-4 h-4 text-green-700 bg-gray-100 border-gray-300 rounded focus:ring-green-600 focus:ring-2"
          />
          <label
            form="affidavit-checkbox"
            className="ms-2 text-sm font-medium text-gray-900"
          >
            Acepto{" "}
            <a
              href="https://drive.google.com/file/d/1ipu9NrLWBl80El-G-IMHOUlQHP9QqoAG/preview"
              className="text-green-700 font-medium underline"
              target="_blank"
            >
              Declaración Jurada
            </a>
          </label>
        </div>
      </div>
    </div>
  );
};
