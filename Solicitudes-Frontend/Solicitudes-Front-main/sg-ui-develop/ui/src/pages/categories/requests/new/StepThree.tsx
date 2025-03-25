import { Alert, Label, Textarea, Toast } from "flowbite-react";
import { FaRegFilePdf } from "react-icons/fa";
import { PDFDocument, PDFSignature } from "pdf-lib";

import { Label as CustomLabel } from "../../../../components/Label/Label";
import StyledTitle from "../../../../components/Title/StyledTitle";
import { Dropzone } from "../../../../components/Dropzone/Dropzone";
import { Activities, DocumentTypes, FilesType } from "../types";
import { useEffect, useRef, useState } from "react";
import { RequestService } from "../requestsService";
import { Button } from "@/components/ui/button.tsx";
import { cn } from "@/lib/utils.ts";
import { ChevronUp } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox.tsx";
import { AlertDescription } from "@/components/ui/alert";
import { FormInfo } from "@/components/Tooltip/Tooltip";

interface StepThreeProps {
  observations: string;
  files: FilesType;
  setFiles: React.Dispatch<React.SetStateAction<FilesType>>;
  selectedActivities: string[];
  setSelectedActivity: (value: string[]) => void;
  estimatedTime: number;
  setEstimatedTime: (value: number) => void;
  insurance: boolean | null;
  setInsurance: (value: boolean | null) => void;
  projectDescription: string;
  setProjectDescription: (value: string) => void;
  setDisabled: (value: boolean) => void;
}

const info = `Si las actividades a realizar incluyen:
<ul className="list-disc pl-5">
  <li>* Vallas y/o andamios en vía pública.</li>
  <li>* Limpieza o pintura de fachadas, revoques y/o trabajos similares hacia la calle.</li>
  <li>* Arreglos de cubiertas de techos y/o salientes hacia la calle.</li>
</ul>
</br>
<p>Es obligatorio adjuntar:</p>
<ul className="list-disc pl-5">
  <li>Una póliza de seguro que contemple cobertura contra terceros y responsabilidad civil por el plazo de obra solicitado, incluyendo accidentes con cláusula de no repetición a favor del Municipio de San Isidro.</li>
</ul>
</br>
<p><strong>Importante:</strong></p>
<ul className="list-disc pl-5">
  <li>No se autorizará bajo ningún concepto la ocupación de la vía pública con materiales de construcción y/o escombros, según los Arts. 5.1.1 y 5.13 del Código de Edificación de San Isidro (C.E.S.I).</li>
</ul>`;

export const StepThree = ({
  observations,
  files,
  setFiles,
  selectedActivities,
  setSelectedActivity,
  estimatedTime,
  setEstimatedTime,
  insurance,
  setInsurance,
  projectDescription,
  setProjectDescription,
  setDisabled,
}: StepThreeProps) => {
  const [activities, setActivities] = useState<Activities[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [observationsArray, setObservationsArray] = useState<string[]>([]);

  useEffect(() => {
    const getActivities = async () => {
      const result = await RequestService.getSiteActivities();
      setActivities(result.data);
    };
    getActivities();
  }, []);

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

  const buttonRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (buttonRef.current && !buttonRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
      document.addEventListener("keydown", handleKeyDown);
    }
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen]);

  const handleSelectedOptionChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    if (e.target.value === "yes") {
      setInsurance(true);
      sessionStorage.setItem("insurance", "true");
    } else if (e.target.value === "no") {
      setInsurance(false);
      deleteFile();
      sessionStorage.setItem("insurance", "false");
    }
  };

  const handleProjectDescription = (
    e: React.ChangeEvent<HTMLTextAreaElement>
  ) => {
    setProjectDescription(e.target.value);
    sessionStorage.setItem("projectDescription", e.target.value);
  };

  const handleSelectActivity = (optionId: string) => {
    const auxSelectedActivities = selectedActivities.includes(optionId)
      ? selectedActivities.filter((id) => id !== optionId)
      : [...(selectedActivities || []), optionId];
    setSelectedActivity(auxSelectedActivities);
    sessionStorage.setItem(
      "selectedActivities",
      `[${auxSelectedActivities.join(",")}]`
    );
  };

  useEffect(() => {
    if (selectedActivities.length > 0) {
      sessionStorage.setItem(
        "activityDescription",
        activities
          .reduce((acc: string[], activity: Activities) => {
            if (selectedActivities.includes(activity.id)) {
              acc.push(activity.description);
            }
            return acc;
          }, [])
          .join(", ")
      );
    }
  }, [selectedActivities, activities]);

  const handleSetEstimatedTime = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEstimatedTime(Number(e.target.value));
    sessionStorage.setItem("estimatedTime", e.target.value);
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
    if (!insurance) {
      alert(
        "Usted indicó que la tarea seleccionada no requiere seguro. No necesita adjuntar un documento aquí."
      );
      return;
    }

    if (await checkForSignatureFields(file)) {
      alert("Error en el documento. Contiene campos de firma.");
      return;
    }

    if (!file) return;
    setFiles((prevFiles: FilesType) => ({
      ...prevFiles,
      [key]: {
        content: file,
        name: fileName,
      },
    }));
  };

  const deleteFile = () => {
    setFiles((prevFiles: FilesType) => {
      const newFiles = { ...prevFiles };
      delete newFiles[DocumentTypes.Insurance.id];
      return newFiles;
    });
  };

  useEffect(() => {
    if (
      projectDescription !== "" &&
      insurance !== null &&
      estimatedTime > 0 &&
      selectedActivities.length > 0
    ) {
      if (insurance) {
        const allDocumentsUploaded = files[DocumentTypes.Insurance.id]
          ? true
          : false;
        setDisabled(!allDocumentsUploaded);
        return;
      }

      setDisabled(false);
    } else {
      setDisabled(true);
    }
  }, [projectDescription, files, insurance, estimatedTime, selectedActivities]);

  return (
    <div>
      <StyledTitle title="Memoria descriptiva" />
      {observationsArray.length > 0 && (
        <div className="space-y-2 my-4">
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
        </div>
      )}
      <div className="w-full">
        <div className="mb-2 block">
          <Label
            htmlFor="comment"
            className="font-normal"
            value="Contanos brevemente qué trabajos vas a realizar en tu domicilio."
          />
        </div>
        <Textarea
          id="comment"
          placeholder="Escribir aquí..."
          value={projectDescription}
          onChange={handleProjectDescription}
          required
          maxLength={320}
          rows={4}
        />
        <div className="mt-4">
          <Label htmlFor="activities" value="Actividades de la obra" />
        </div>
        <div className="space-y-2">
          <div className="relative" ref={buttonRef}>
            <Button
              type="button"
              onClick={() => setIsOpen(!isOpen)}
              className={cn(
                "w-full justify-between border-2 border-[#3b5c3f] bg-white text-gray-500 hover:bg-white hover:text-gray-600",
                "rounded-lg px-4 py-6 text-base font-normal overflow-y-hidden"
              )}
            >
              {selectedActivities.length
                ? activities
                    .reduce((acc: string[], activity: Activities) => {
                      if (selectedActivities.includes(activity.id)) {
                        acc.push(activity.description);
                      }
                      return acc;
                    }, [])
                    .join(", ")
                : "Seleccionar una o más opciones según corresponda"}
              <ChevronUp
                className={cn(
                  "h-4 w-4 shrink-0 transition-transform",
                  !isOpen && "rotate-180"
                )}
              />
            </Button>

            {isOpen && (
              <div className="z-50 w-full space-y-4 my-2">
                <div className="max-h-[300px] overflow-auto rounded-xl border-2 bg-white border-gray-300">
                  {activities?.map((option) => (
                    <div
                      key={option.id}
                      className={cn(
                        "flex items-center space-x-3 px-4 py-4",
                        selectedActivities.includes(option.id) &&
                          "bg-[#3E59421A]"
                      )}
                    >
                      <Checkbox
                        id={option.id}
                        checked={selectedActivities.includes(option.id)}
                        onCheckedChange={() => {
                          handleSelectActivity(option.id);
                        }}
                        className="border-[#3b5c3f] data-[state=checked]:bg-white data-[state=checked]:text-[#3E5942]"
                      />
                      <Label
                        htmlFor={option.id}
                        className="flex-grow cursor-pointer text-gray-500"
                      >
                        {option.description}
                      </Label>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
        <div className="mt-4">
          <Label htmlFor="time" value="Tiempo estimado de la obra" />
        </div>
        <div className="mt-2">
          <label>
            <input
              type="radio"
              value="30"
              className="text-green-700 focus:ring-green-600"
              checked={estimatedTime === 30}
              onChange={handleSetEstimatedTime}
            />
            <strong className="ml-2 text-sm text-gray-500">30 días</strong>
          </label>
          <label className="ml-4">
            <input
              type="radio"
              value="60"
              className="text-green-700 focus:ring-green-600"
              checked={estimatedTime === 60}
              onChange={handleSetEstimatedTime}
            />
            <strong className="ml-2 text-sm text-gray-500">60 días</strong>
          </label>
          <label className="ml-4">
            <input
              type="radio"
              value="90"
              className="text-green-700 focus:ring-green-600"
              checked={estimatedTime === 90}
              onChange={handleSetEstimatedTime}
            />
            <strong className="ml-2 text-sm text-gray-500">90 días</strong>
          </label>
        </div>
        <div className="mt-4 flex items-center space-x-0">
          <Label
            htmlFor="building"
            value="¿La tarea seleccionada requiere seguro?"
          />
          <FormInfo
            info={""}
            placement="right-end"
            tooltipContent={{
              title: "Información Importante",
              content: <div dangerouslySetInnerHTML={{ __html: info }} />,
            }}
          />
        </div>
        <div className="mt-2">
          <label>
            <input
              type="radio"
              value="yes"
              className="text-green-700 focus:ring-green-600"
              checked={insurance === true}
              onChange={handleSelectedOptionChange}
            />
            <strong className="ml-2 text-sm text-gray-500">SI</strong>
          </label>
          <label className="ml-4">
            <input
              type="radio"
              value="no"
              className="text-green-700 focus:ring-green-600"
              checked={insurance === false}
              onChange={handleSelectedOptionChange}
            />
            <strong className="ml-2 text-sm text-gray-500">NO</strong>
          </label>
        </div>
        <div className="mt-4">
          <CustomLabel text="Póliza de seguro" mandatory={false} />
          <Dropzone
            id={DocumentTypes.Insurance.description}
            onDrop={(file: File, fileName: string) =>
              handleFileUpload(file, fileName, DocumentTypes.Insurance.id)
            }
          />
          <div className="mt-2">
            {files[DocumentTypes.Insurance.id]?.name && (
              <Toast>
                <div className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg">
                  <FaRegFilePdf className="h-5 w-5" />
                </div>
                <div className="ml-3 text-sm font-normal">
                  Archivo subido: {files[DocumentTypes.Insurance.id]?.name}
                </div>
                <Toast.Toggle onClick={deleteFile} />
              </Toast>
            )}
          </div>
        </div>
        <Label className="font-normal text-xs text-opacity-50">
          (*) La póliza de seguro contra terceros debe incluir responsabilidad
          civil por la duración de la obra y un seguro de accidentes con
          cláusula de no repetición al Municipio de San Isidro.
        </Label>
      </div>
    </div>
  );
};
