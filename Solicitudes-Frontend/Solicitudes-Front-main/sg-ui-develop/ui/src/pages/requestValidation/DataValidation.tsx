import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { AlertTriangle } from "lucide-react";
import { useLocation, useNavigate } from "react-router-dom";
import { useStatusMessageStore } from "@/stores/statusMessageStore.ts";
import { useLoadingScreenStore } from "@/stores/loadingScreenStore.ts";
import ValidationCard from "@/components/ValidationCard/ValidationCard.tsx";
import { Section, Validation } from "@/pages/requestValidation/types.ts";
import {
  docsValidationParser,
  getDocumentByGedoCode,
  getRequestValidationData,
  putRequestValidation,
} from "@/pages/requestValidation/utils.ts";
import { RequestError } from "@/pages/categories/requests/types.ts";
import { useActionDialogStore } from "@/stores/actionDialogStore.ts";
import { base64ToFile, replacePlaceholders } from "@/lib/utils.ts";
import authorizationWithInsuranceTemplate from "./docsToSign/authorizationWithInsuranceTemplate.html?raw";
import authorizationWithoutInsuranceTemplate from "./docsToSign/authorizationWithoutInsuranceTemplate.html?raw";
import { ScrollArea } from "@/components/ui/scroll-area.tsx";

export default function DataValidation() {
  const navigate = useNavigate();
  // TODO: CHECK IF NECESSARY
  const currentDocIndex = 0;
  const { state: currentValidation }: { state: Validation } = useLocation();
  const [validation, setValidation] = useState<Validation>(currentValidation);
  // TODO: CAMBIAR BANDERA
  const [success, setSuccess] = useState(false);
  const [authorizationDocument, setAuthorizationDocument] =
    useState<string>("");
  const [authorizationReference, setAuthorizationReference] =
    useState<string>("");
  const showMessage = useStatusMessageStore((state) => state.showMessage);
  const hideMessage = useStatusMessageStore((state) => state.hideMessage);
  const showLoading = useLoadingScreenStore((state) => state.showLoading);
  const hideLoading = useLoadingScreenStore((state) => state.hideLoading);
  const showActionDialog = useActionDialogStore((state) => state.showDialog);
  const [, setError] = useState("");

  const getValidationData = async () => {
    showLoading();
    try {
      let validationData = await getRequestValidationData(
        validation.recordNumber
      );
      setValidation((prevState) => {
        const validation = {
          ...prevState,
          ...validationData,
        };
        const parsedValidation = docsValidationParser(validation);
        setAuthorizationDocument(
          parseAuthorizationDocument(parsedValidation) ?? ""
        );
        return parsedValidation;
      });
      setSuccess(true);
    } catch (e) {
      if (e instanceof RequestError) {
        setError(e.message);
      } else {
        setError("Ocurrió un error en la busqueda de información.");
      }
    }
    hideLoading();
  };

  const getData = async () => {
    await getValidationData();
  };

  useEffect(() => {
    getData();
  }, []);

  const parseAuthorizationDocument = (parsedValidation: Validation) => {
    let rawTemplate;
    if (parsedValidation.requiresInsurance) {
      setAuthorizationReference(
        "Acto Administrativo- Aviso de obra (CON PÓLIZA)"
      );
      rawTemplate = authorizationWithInsuranceTemplate;
    } else {
      setAuthorizationReference(
        "Acto Administrativo- Aviso de obra (SIN PÓLIZA)"
      );
      rawTemplate = authorizationWithoutInsuranceTemplate;
    }

    const parser = new DOMParser();
    const doc = parser.parseFromString(rawTemplate, "text/html");
    const template = doc.querySelector("#verification-document");
    if (template) {
      const idxParser: Record<number, string> = {
        1: "firstDocGedoCode",
        2: "secondDocGedoCode",
        3: "thirdDocGedoCode",
        4: "fourthDocGedoCode",
        5: "fifthDocGedoCode",
        6: "sixthDocGedoCode",
        7: "seventhDocGedoCode",
        8: "eighthDocGedoCode",
      };
      const replacements = {
        recordNumber: parsedValidation.recordNumber,
        requesterFullName: parsedValidation.requesterFullName,
        requesterCuil: parsedValidation.requesterCuil,
        requesterAddress: parsedValidation.requesterAddress,
        ...parsedValidation.replacementIFs.reduce(
          (acc: Record<string, string>, replacement, idx) => {
            acc[idxParser[idx + 1]] = replacement.gedoCode;
            return acc;
          },
          {}
        ),
        estimatedTime: parsedValidation.estimatedTime,
        assignedTask: parsedValidation.assignedTask.join("</li><li>"),
      };
      hideLoading();
      return replacePlaceholders(template.innerHTML, replacements);
    }
  };

  const hasInvalidSections = validation.documents?.some((doc) =>
      doc.sections?.some((section) => section.status === "invalid")
  );

  const handleSectionUpdate = (
    sectionId: string,
    updates: Partial<Section>,
    updateNext: boolean = false
  ) => {
    setValidation((prev) => {
      const newVerification = { ...prev };
      const documents = [...prev.documents];
      const documentIndex = currentDocIndex;
      const document = { ...documents[documentIndex] };
      const sections = [...document.sections];
      const sectionIndex = sections.findIndex((s) => s.id === sectionId);

      if (sectionIndex !== -1) {
        sections[sectionIndex] = {
          ...sections[sectionIndex],
          ...updates,
        };

        if (updateNext) {
          if (
            sections[sectionIndex].status === "valid" ||
            (sections[sectionIndex].status === "invalid" && sectionIndex != 0)
          ) {
            const nextSection = sections[sectionIndex + 1];
            if (nextSection) {
              sections[sectionIndex + 1] = {
                ...nextSection,
                show: true,
              };
            }
          }
        }

        document.sections = sections;
        documents[documentIndex] = document;
        newVerification.documents = documents;
      }

      return newVerification;
    });
  };

  const handleSign = async () => {
    if (hasInvalidSections) {
      showActionDialog({
        title: "Validación final",
        description:
            "Estás por rechazar el proceso de validación. Al confirmar, se notificará al empleado municipal que la validación fue rechazada.",
        icon: AlertTriangle,
        primaryButtonLabel: "Confirmar",
        primaryButtonAction: async () => {
          showLoading();
          try {
            await putRequestValidation(
              validation.recordNumber,
              JSON.stringify({
                is_valid: false,
              })
            );
            hideLoading();
            showMessage({
              type: "success",
              title: ["¡Listo!"],
              description: [
                "La validación del expediente ",
                `N° ${validation.recordNumber} `,
                "se rechazó con éxito.",
              ],
              buttonText: "Finalizar",
              onAction: () => {
                navigate("/admin-panel/request-validation");
                hideMessage();
              },
            });
          } catch (e) {
            hideLoading();
            showMessage({
              type: "error",
              title: ["No se pudo completar el rechazo"],
              description: [
                "Hubo un problema al procesar la validación. ",
                "Por favor, intentá de nuevo.",
              ],
              buttonText: "Reintentar",
              onAction: () => {
                hideMessage();
              },
            });
          }
        },
        secondaryButtonLabel: "Cancelar",
      });
    } else {
      const htmlBase64 = btoa(
          encodeURIComponent(authorizationDocument).replace(
              /%([0-9A-F]{2})/g,
              (_match, p1) => String.fromCharCode(parseInt(p1, 16))
          )
      );
      showActionDialog({
        title: "Validación final",
        description:
            "Estás por completar el proceso de validación. Al confirmar, se notificará al vecino que la solicitud se aprobó con éxito.",
        icon: AlertTriangle,
        primaryButtonLabel: "Confirmar",
        primaryButtonAction: async () => {
          showLoading();
          const body: {
            is_valid: boolean;
            reference: string;
            authorizationDocument: string;
          } = {
            is_valid: true,
            reference: authorizationReference,
            authorizationDocument: htmlBase64,
          };
          try {
            await putRequestValidation(
                validation.recordNumber,
                JSON.stringify(body)
            );
            hideLoading();
            showMessage({
              type: "success",
              title: ["¡Listo!"],
              description: [
                "La validación del expediente ",
                `N° ${validation.recordNumber} `,
                "se completó con éxito.",
              ],
              buttonText: "Finalizar",
              onAction: () => {
                navigate("/admin-panel/request-validation");
                hideMessage();
              },
            });
          } catch (e) {
            hideLoading();
            showMessage({
              type: "error",
              title: ["No se pudo completar la firma"],
              description: [
                "Hubo un problema al procesar la validación. ",
                "Por favor, intentá de nuevo.",
              ],
              buttonText: "Reintentar",
              onAction: () => {
                hideMessage();
              },
            });
          }
        },
        secondaryButtonLabel: "Cancelar",
      });
    }
  };

  const handleClose = async () => {
    showActionDialog({
      title: "¿Seguro que querés salir de la revisión?",
      description:
        "No terminaste de validar este aviso de obra. Si salís ahora, la tarea realizada no se guardará en el sistema.",
      icon: AlertTriangle,
      primaryButtonLabel: "Descartar tarea",
      primaryButtonAction: async () => {
        navigate(-1);
      },
      secondaryButtonLabel: "Volver",
    });
  };

  return (
    <>
      <div className="flex">
        {/* Left Side - Validation Form */}
        <div className="w-[400px] px-6 pt-6 border-r">
          <div className="flex flex-col justify-between min-h-[calc(100vh-50px)]">
            <div className="space-y-6">
              <div>
                <h1 className="text-2xl font-bold">Validar solicitud</h1>
                <p className="mt-2 text-sm text-gray-600">
                  A continuación, se listarán los requisitos verificados del
                  expediente.
                </p>
                <p className="mt-2 text-sm text-gray-600">
                  Para ver más detalles, presiona sobre cada uno. En caso de
                  rechazar un documento, este será nuevamente verificado.
                </p>
              </div>

              <div key="keyID" className="space-y-2">
                {/*TODO: FIX MAGIC NUMBER*/}
                {success &&
                  validation.documents[0].sections.map(
                    (section) =>
                      section.show && (
                        <ValidationCard
                          key={section.id}
                          section={section}
                          isValidation={true}
                          onUpdate={(updates, updateNext) =>
                            handleSectionUpdate(section.id, updates, updateNext)
                          }
                          onQuestionClick={async (index) => {
                            const { gedoCode = "", title } =
                              validation.documents[index];

                            const content = await getDocumentByGedoCode(
                              gedoCode
                            );
                            const file = base64ToFile(content, title);

                            const fileURL = URL.createObjectURL(file);
                            window.open(fileURL, "_blank");
                          }}
                        ></ValidationCard>
                      )
                  )}
              </div>
            </div>

            <div>
              <Button
                className="w-full bg-[#3b5c3f] hover:bg-[#2e4831] text-white mt-4"
                onClick={() => handleSign()}
              >
                { hasInvalidSections ? "Rechazar" : "Firmar" }
              </Button>
              <Button
                className="w-full mt-4"
                variant="outline"
                onClick={() => handleClose()}
              >
                Cerrar
              </Button>
            </div>
          </div>
        </div>

        {/* Right Side - PDF Viewer */}
        <div className="flex-1 flex flex-col">
          <div className="border-b p-4">
            <h1 className="text-xl font-bold">
              Expediente: {validation.recordNumber}
            </h1>
          </div>

          {success && (
            <ScrollArea>
              <div
                className="pt-4 max-h-[calc(100vh-50px)]"
                dangerouslySetInnerHTML={{
                  __html: authorizationDocument,
                }}
              />
            </ScrollArea>
          )}
        </div>
      </div>
    </>
  );
}
