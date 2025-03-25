import { useEffect, useState } from "react";
import { Alert, Toast } from "flowbite-react";
import { FaCloudDownloadAlt, FaRegFilePdf } from "react-icons/fa";
import { PDFDocument, PDFSignature } from "pdf-lib";

import { Label as CustomLabel } from "../../../../components/Label/Label";
import Info from "../../../../components/Info/Info";
import StyledTitle from "../../../../components/Title/StyledTitle";
import { Dropzone } from "../../../../components/Dropzone/Dropzone";
import { DocumentTypes, FilesType, UserDocumentRequirements } from "../types";
import { AlertDescription } from "@/components/ui/alert";

interface StepTwoProps {
  commonZone: boolean;
  observations: string;
  userType: string;
  setDisabled: (value: boolean) => void;
  setUserType: (value: string) => void;
  files: FilesType;
  setFiles: React.Dispatch<React.SetStateAction<FilesType>>;
}

export const StepTwo = ({
  commonZone,
  observations,
  userType,
  setUserType,
  setDisabled,
  files,
  setFiles,
}: StepTwoProps) => {
  const [title, _] = useState("Verificación de permisos sobre el inmueble");
  const [selectedOption, setSelectedOption] = useState("");
  const [observationsArray, setObservationsArray] = useState<string[]>([]);
  const [isInitialized, setIsInitialized] = useState(false);
  const [localFiles, setLocalFiles] = useState<FilesType>({});

  useEffect(() => {
    if (observations && !isInitialized && Object.keys(files).length > 0) {
      setLocalFiles(files);
      setIsInitialized(true);
    }
  }, [observations, files, isInitialized]);

  useEffect(() => {
    if (observations) {
      const processedObservations = observations
        .split("</br>")
        .map((obs) => obs.trim())
        .filter((obs) => obs !== "");

      setObservationsArray(processedObservations);
    } else {
      setObservationsArray([]);
    }
  }, [observations]);

  const handleSelectedOptionChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    if (e.target.value !== selectedOption) {
      setFiles({});
    }
    setSelectedOption(e.target.value);
    setUserType(e.target.value);
  };

  const handleDownload = (fileName: string, fileContent: File) => {
    const url = URL.createObjectURL(fileContent);
    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    link.click();
    URL.revokeObjectURL(url);
  };

  const checkForSignatureFields = async (file: File): Promise<boolean> => {
    const arrayBuffer = await file.arrayBuffer();
    const pdfDoc = await PDFDocument.load(arrayBuffer);
    const form = pdfDoc.getForm();

    const signatureFields = form
      .getFields()
      .filter((field) => field instanceof PDFSignature);

    return signatureFields.length > 0;
  };

  const handleFileUpload = async (
    file: File,
    fileName: string,
    key: number
  ) => {
    if (!file) return;

    if (await checkForSignatureFields(file)) {
      alert("Error en el documento. Contiene campos de firma.");
      return;
    }

    setFiles((prevFiles: FilesType) => ({
      ...prevFiles,
      [key]: {
        content: file,
        name: fileName,
      },
    }));
  };

  const handleToastClose = (key: number) => {
    setFiles((prevFiles: FilesType) => {
      const newFiles = { ...prevFiles };
      delete newFiles[key];
      return newFiles;
    });
  };

  useEffect(() => {
    if (commonZone) {
      //setTitle((prevTitle) => prevTitle + " en espacio común");
      setUserType("Admin");
    }
  }, [commonZone]);

  useEffect(() => {
    if (userType !== "") {
      setSelectedOption(userType);
      sessionStorage.setItem("userType", userType);
    }
  }, [userType]);

  useEffect(() => {
    if (userType !== "" && Object.keys(files).length > 0) {
      const requiredDocuments =
        UserDocumentRequirements[
          userType as keyof typeof UserDocumentRequirements
        ] || [];

      const allDocumentsUploaded = requiredDocuments.every(
        (docId) => files[docId]
      );
      setDisabled(!allDocumentsUploaded);
    } else {
      setDisabled(true);
    }
  }, [userType, files]);

  return (
    <div>
      <StyledTitle title={title} />
      {observationsArray.length > 0 && (
        <div className="space-y-2 my-4">
          <Info text="Por favor, actualizá la documentación según las observaciones indicadas." />
          <br />
          <span>Observaciones</span>
          <Alert className="bg-gray-200/80 border-none p-2">
            <AlertDescription className="text-gray-700 ml-2">
              <ul className="list-disc list-inside">
                {observationsArray.map((observation, index) => (
                  <li key={index}>{observation}</li>
                ))}
              </ul>
            </AlertDescription>
          </Alert>
          <br />
          <span>Archivos cargados previamente</span>
          <div className="mt-2 grid grid-cols-1 sm:grid-cols-2 gap-4">
            {Object.keys(localFiles).length > 0 ? (
              Object.entries(localFiles)
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
              </>
            )}
          </div>
        </div>
      )}
      {commonZone ? (
        <>
          {observationsArray.length === 0 && (
            <div>
              <Info text="Solo el Administrador del Consorcio puede solicitar la autorización." />
              <Info text="Cargá los documentos que acreditan tu carácter de Administrador de Consorcio designado." />
            </div>
          )}
          <div className="mt-4">
            <CustomLabel
              text="Reglamento de Co-Propiedad (en PDF)"
              mandatory={true}
            />
            <Dropzone
              id={DocumentTypes.CoOwnership.description}
              onDrop={(file: File, fileName: string) =>
                handleFileUpload(file, fileName, DocumentTypes.CoOwnership.id)
              }
            />
            <div className="mt-2">
              {files[DocumentTypes.CoOwnership.id]?.name && (
                <Toast>
                  <div className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg">
                    <FaRegFilePdf className="h-5 w-5" />
                  </div>
                  <div className="ml-3 text-sm font-normal">
                    Archivo subido: {files[DocumentTypes.CoOwnership.id]?.name}
                  </div>
                  <Toast.Toggle
                    onClick={() =>
                      handleToastClose(DocumentTypes.CoOwnership.id)
                    }
                  />
                </Toast>
              )}
            </div>
          </div>
          <div className="mt-4">
            <CustomLabel text="Acta de designación (en PDF)" mandatory={true} />
            {observationsArray.length === 0 && (
              <Info text="Si tu designación como administrador está contemplada en el Reglamento de Co-propiedad y tu mandato sigue vigente, podés adjuntar el Reglamento para acreditar tu rol y legitimación para avanzar con la solicitud." />
            )}
            <Dropzone
              id={DocumentTypes.AppointmentCertificate.description}
              onDrop={(file: File, fileName: string) =>
                handleFileUpload(
                  file,
                  fileName,
                  DocumentTypes.AppointmentCertificate.id
                )
              }
            />
            <div className="mt-2">
              {files[DocumentTypes.AppointmentCertificate.id]?.name && (
                <Toast>
                  <div className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg">
                    <FaRegFilePdf className="h-5 w-5" />
                  </div>
                  <div className="ml-3 text-sm font-normal">
                    Archivo subido:{" "}
                    {files[DocumentTypes.AppointmentCertificate.id]?.name}
                  </div>
                  <Toast.Toggle
                    onClick={() =>
                      handleToastClose(DocumentTypes.AppointmentCertificate.id)
                    }
                  />
                </Toast>
              )}
            </div>
          </div>
        </>
      ) : (
        <>
          <Info text="Por favor, indicá el tipo de relación que tenés con el inmueble:" />
          <div className="mt-4 mb-4">
            <label>
              <input
                type="radio"
                value="Owner"
                disabled={observationsArray.length > 0}
                className="text-green-700 focus:ring-green-600"
                checked={selectedOption === "Owner"}
                onChange={handleSelectedOptionChange}
              />
              <strong className="ml-2 text-sm text-gray-500">
                Soy titular del inmueble
              </strong>
            </label>
            <label className="ml-4">
              <input
                type="radio"
                value="Occupant"
                disabled={observationsArray.length > 0}
                className="text-green-700 focus:ring-green-600"
                checked={selectedOption === "Occupant"}
                onChange={handleSelectedOptionChange}
              />
              <strong className="ml-2 text-sm text-gray-500">
                Soy ocupante o autorizado legítimo
              </strong>
            </label>
          </div>
          {selectedOption === "Owner" ? (
            <div>
              <Info text="Necesitamos que acredites tu titularidad sobre el inmueble." />
              <Info text="Se requiere la carga de uno de los siguientes documentos que respalden tu potestad:" />
              <div className="mt-4">
                <CustomLabel
                  text="Título de propiedad o Informe de dominio (en PDF)"
                  mandatory={true}
                />
                <Dropzone
                  id={DocumentTypes.PropertyTitle.description}
                  onDrop={(file: File, fileName: string) =>
                    handleFileUpload(
                      file,
                      fileName,
                      DocumentTypes.PropertyTitle.id
                    )
                  }
                />
                <div className="mt-2">
                  {files[DocumentTypes.PropertyTitle.id]?.name && (
                    <Toast>
                      <div className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg">
                        <FaRegFilePdf className="h-5 w-5" />
                      </div>
                      <div className="ml-3 text-sm font-normal">
                        Archivo subido:{" "}
                        {files[DocumentTypes.PropertyTitle.id]?.name}
                      </div>
                      <Toast.Toggle
                        onClick={() =>
                          handleToastClose(DocumentTypes.PropertyTitle.id)
                        }
                      />
                    </Toast>
                  )}
                </div>
              </div>
            </div>
          ) : (
            selectedOption === "Occupant" && (
              <div>
                <div className="mt-4">
                  <CustomLabel
                    text="Título de propiedad o Informe de dominio (en PDF)"
                    mandatory={true}
                  />
                  <Info text="Se requiere la carga de uno de los siguientes documentos" />
                  <Dropzone
                    id={DocumentTypes.PropertyTitle.description}
                    onDrop={(file: File, fileName: string) =>
                      handleFileUpload(
                        file,
                        fileName,
                        DocumentTypes.PropertyTitle.id
                      )
                    }
                  />
                  <div className="mt-2">
                    {files[DocumentTypes.PropertyTitle.id]?.name && (
                      <Toast>
                        <div className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg">
                          <FaRegFilePdf className="h-5 w-5" />
                        </div>
                        <div className="ml-3 text-sm font-normal">
                          Archivo subido:{" "}
                          {files[DocumentTypes.PropertyTitle.id]?.name}
                        </div>
                        <Toast.Toggle
                          onClick={() =>
                            handleToastClose(DocumentTypes.PropertyTitle.id)
                          }
                        />
                      </Toast>
                    )}
                  </div>
                </div>
                <div className="mt-4">
                  <CustomLabel
                    text="Autorización del propietario (en PDF)"
                    mandatory={true}
                  />

                  <Info
                    text={
                      <>
                        Subí acá tu autorización firmada por el propietario. Si
                        todavía no la tenés podés{" "}
                        <a
                          href="https://docs.google.com/document/d/13B3jaCmxH3AdS3Da6SaN4Wnv0HgT4pXLBJGiDEWKvkU/export?format=doc"
                          className="text-green-700 font-medium"
                        >
                          descargar nuestra plantilla
                        </a>{" "}
                        modelo para completarla.
                      </>
                    }
                  />
                  <Dropzone
                    id={DocumentTypes.OwnerAuthorization.description}
                    onDrop={(file: File, fileName: string) =>
                      handleFileUpload(
                        file,
                        fileName,
                        DocumentTypes.OwnerAuthorization.id
                      )
                    }
                  />
                  <div className="mt-2">
                    {files[DocumentTypes.OwnerAuthorization.id]?.name && (
                      <Toast>
                        <div className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg">
                          <FaRegFilePdf className="h-5 w-5" />
                        </div>
                        <div className="ml-3 text-sm font-normal">
                          Archivo subido:{" "}
                          {files[DocumentTypes.OwnerAuthorization.id]?.name}
                        </div>
                        <Toast.Toggle
                          onClick={() =>
                            handleToastClose(
                              DocumentTypes.OwnerAuthorization.id
                            )
                          }
                        />
                      </Toast>
                    )}
                  </div>
                </div>
              </div>
            )
          )}
        </>
      )}
    </div>
  );
};
