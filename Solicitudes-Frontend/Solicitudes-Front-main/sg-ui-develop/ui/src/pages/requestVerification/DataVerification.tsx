import { Fragment, useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AlertTriangle } from "lucide-react";
import { cn, replacePlaceholders } from "@/lib/utils";
import { useLocation, useNavigate } from "react-router-dom";
import {
  Document,
  Section,
  Verification,
} from "@/pages/requestVerification/types.ts";
import { useStatusMessageStore } from "@/stores/statusMessageStore.ts";
import {
  docsVerificationParser,
  getRequestVerificationData,
  putRequestVerification,
} from "@/pages/requestVerification/utils.ts";
import { RequestError } from "@/pages/categories/requests/types.ts";
import ownerTemplate from "./docsToSign/ownerTemplate.html?raw";
import administratorTemplate from "./docsToSign/administratorTemplate.html?raw";
import renterTemplate from "./docsToSign/renterTemplate.html?raw";
import observationTemplate from "./docsToSign/observationTemplate.html?raw";
import tasksWithInsuranceTemplate from "./docsToSign/tasksWithInsuranceTemplate.html?raw";
import tasksWithoutInsuranceTemplate from "./docsToSign/tasksWithoutInsuranceTemplate.html?raw";
import { useLoadingScreenStore } from "@/stores/loadingScreenStore.ts";
import ValidationCard from "@/components/ValidationCard/ValidationCard.tsx";
import { getRequestDocument } from "@/pages/requestValidation/utils.ts";
import { useActionDialogStore } from "@/stores/actionDialogStore.ts";
import { RequestService } from "@/pages/categories/requests/requestsService.ts";

export default function DataVerification() {
  const [currentDocIndex, setCurrentDocIndex] = useState(0);
  const navigate = useNavigate();
  const { state: currentVerification }: { state: Verification } = useLocation();
  const [verification, setVerification] =
    useState<Verification>(currentVerification);
  // TODO: CAMBIAR BANDERA
  const [successData, setSuccessData] = useState(false);
  const showMessage = useStatusMessageStore((state) => state.showMessage);
  const hideMessage = useStatusMessageStore((state) => state.hideMessage);
  const showLoading = useLoadingScreenStore((state) => state.showLoading);
  const hideLoading = useLoadingScreenStore((state) => state.hideLoading);
  const showActionDialog = useActionDialogStore((state) => state.showDialog);
  const [, setError] = useState("");

  const getVerificationData = useCallback(async () => {
    showLoading();
    try {
      const updatedVerification = { ...verification };
      updatedVerification.documents = await getRequestVerificationData(
        verification.recordNumber
      );

      const verificationParsed = docsVerificationParser(updatedVerification);
      const {
        documents: [firstDoc, secondDoc],
      } = verificationParsed;

      const updateDocBase64 = async (doc: Document) => {
        if (doc?.gedoCode) {
          try {
            return (await getRequestDocument(doc.gedoCode)).content;
          } catch (e) {
            setError(
              e instanceof RequestError
                ? e.message
                : "Error al buscar información."
            );
            return null;
          }
        }
      };

      if (firstDoc?.gedoCode) {
        firstDoc.base64 = await updateDocBase64(firstDoc);
      }
      if (secondDoc?.gedoCode) {
        secondDoc.base64 = await updateDocBase64(secondDoc);
      }

      if (verificationParsed.verificationCase === "4") {
        try {
          const { data } = await RequestService.getSiteActivities();
          firstDoc.sections[0].observationOptions = data;
        } catch (e) {
          firstDoc.sections[0].observationOptions = [];
        }
      }
      setVerification(verificationParsed);
      setSuccessData(true);
    } catch (e) {
      if (e instanceof RequestError) {
        setError(e.message);
      } else {
        setError("Ocurrió un error en la busqueda de información.");
      }
    }
    hideLoading();
  }, [verification.recordNumber, showLoading, hideLoading]);

  useEffect(() => {
    getVerificationData();
  }, [getVerificationData]);

  const handleSectionUpdate = (
    sectionId: string,
    updates: Partial<Section>,
    updateNext: boolean = false
  ) => {
    setVerification((prev) => {
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
            (sections[sectionIndex].status === "invalid" &&
              sections[sectionIndex].id !== "document")
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

        // TODO: CHECK CASE IDENTIFICATORS
        if (verification.verificationCase === "3") {
          const [
            {
              sections: [, , grantedByAuthSection],
            },
            {
              sections: [, , fullNameSection],
            },
          ] = documents;
          fullNameSection.question = grantedByAuthSection.value;
        }
      }
      return newVerification;
    });
  };

  const hasInvalidSections = verification.documents?.some((doc) =>
    doc.sections.some((section) => section.status === "invalid")
  );

  const handleContinue = async () => {
    if (currentDocIndex < verification.documents.length - 1) {
      setCurrentDocIndex(currentDocIndex + 1);
    } else {
      showLoading();
      let docReference = "VERIFICACIÓN CON OBSERVACIONES";
      let rawTemplate = "";

      let grantedByAuthValue = "";
      if (verification.documents?.[0]?.sections) {
        const grantedByAuthSection = verification.documents[0].sections.find(
          (section) => section.id === "grantedByAuth"
        );
        grantedByAuthValue = grantedByAuthSection?.value || "";
      }

      switch (verification.verificationCase) {
        case "1":
          docReference = "CARÁCTER INVOCADO: TITULAR";
          rawTemplate = ownerTemplate;
          break;
        case "2":
          docReference = "CARÁCTER INVOCADO: ADMINISTRADOR";
          rawTemplate = administratorTemplate;
          break;
        case "3":
          docReference = "CARÁCTER INVOCADO: OCUPANTE O INQUILINO";
          rawTemplate = renterTemplate;
          break;
        case "4":
          if (
            verification.documents[0].sections.find(
              (section: Section) => section.id === "insurance"
            )?.value === "Si"
          ) {
            docReference = "VERIFICACIÓN TAREAS EXITOSA (CON PÓLIZA)";
            rawTemplate = tasksWithInsuranceTemplate;
          } else {
            docReference = "VERIFICACIÓN TAREAS EXITOSA (SIN PÓLIZA)";
            rawTemplate = tasksWithoutInsuranceTemplate;
          }
          break;
      }
      const parser = new DOMParser();
      const doc = parser.parseFromString(
        hasInvalidSections ? observationTemplate : rawTemplate,
        "text/html"
      );
      const template = doc.querySelector("#verification-document");
      if (template) {
        const replacements = {
          requesterFullName: verification.requesterFullName,
          requesterCuil: verification.requesterCuil,
          requesterAddress: verification.requesterAddress,
          firstDocGedoCode: verification.documents?.[0]?.gedoCode ?? " ",
          secondDocGedoCode: verification.documents?.[1]?.gedoCode ?? " ",
          invokedCharacter: verification.invokedCharacter,
          observations: joinObservations(),
          grantedToAuth: verification.requesterFullName ?? " ",
          grantedByAuth: grantedByAuthValue,
        };
        hideLoading();
        const parsedTemplate = replacePlaceholders(
          template.innerHTML,
          replacements
        );
        showActionDialog({
          title: "Verificación final",
          description:
            "Último paso. Verificá que toda la información sea correcta antes de firmar el documento.",
          primaryButtonLabel: "Firmar",
          primaryButtonAction: () => {
            handleSign(parsedTemplate, docReference);
          },
          secondaryButtonLabel: "Volver",
          children: (
            <>
              <ScrollArea className="h-[60vh] rounded-md border p-4">
                <div
                  dangerouslySetInnerHTML={{
                    __html: parsedTemplate,
                  }}
                />
              </ScrollArea>

              {hasInvalidSections && (
                <Alert className="bg-yellow-50 text-yellow-900 border-yellow-200">
                  <AlertTriangle className="h-4 w-4 text-yellow-600" />
                  <AlertDescription>
                    Al firmar, se enviará al vecino el documento con las
                    observaciones realizadas.
                  </AlertDescription>
                </Alert>
              )}
            </>
          ),
        });
      }
    }
  };

  const joinObservations = (): string => {
    return verification.documents.reduce((verifObs: string, doc: Document) => {
      const docObs = doc.sections.reduce((docObs: string, section: Section) => {
        if (section.observation) {
          docObs += `${docObs.length ? "</br>" : ""} ${section.observation}`;
        }
        return docObs;
      }, "");
      if (docObs) {
        verifObs += `${verifObs.length ? "</br>" : ""} ${docObs}`;
      }
      return verifObs;
    }, "");
  };

  const handleSign = async (htmlTemplate: string, docReference: string) => {
    showLoading({
      title: ["¡Ya casi estamos!"],
      description: [
        "Estamos procesando tu información, por favor no cierres la pantalla.",
      ],
    });

    const htmlBase64 = btoa(
      encodeURIComponent(htmlTemplate).replace(
        /%([0-9A-F]{2})/g,
        (_match, p1) => String.fromCharCode(parseInt(p1, 16))
      )
    );

    // TODO: CHECK TYPING
    const body: {
      observations: string;
      finalVerificationDocument: string;
      assignedTask?: string[];
      verificationType: string;
      reference: string;
    } = {
      observations: joinObservations(),
      finalVerificationDocument: htmlBase64,
      verificationType:
        verification.verificationCase === "4" ? "tasks" : "property",
      reference: docReference,
    };

    // TODO: CHECK CASE VALUE
    if (verification.verificationCase === "4") {
      body["assignedTask"] = verification.documents[0].sections.find(
        (section: Section) => section.id === "assignedTask"
      )?.selectedObservationOptions;
    }

    try {
      await putRequestVerification(
        verification.recordNumber,
        JSON.stringify(body)
      );
      showMessage({
        type: "success",
        title: ["¡Listo!"],
        description: [
          "La verificación del expediente ",
          `N° ${verification.recordNumber} `,
          "se completó con éxito.",
        ],
        buttonText: "Finalizar",
        onAction: () => {
          navigate("/admin-panel/request-verification");
          hideMessage();
        },
      });
      hideLoading();
    } catch (e) {
      showMessage({
        type: "error",
        title: ["No se pudo completar la firma"],
        description: [
          "Hubo un problema al procesar la verificación.",
          "Por favor, intentá de nuevo.",
        ],
        buttonText: "Reintentar",
        onAction: () => {
          hideMessage();
        },
      });
      hideLoading();
    }
  };

  return (
    <>
      <div className="flex bg-gray-100 min-h-[calc(100vh-80px)]">
        {/* Left Side - Verification Form */}
        <div className="w-[400px] px-6 pt-6 bg-white border-r">
          <div className="flex flex-col justify-between h-full">
            <div className="space-y-6">
              <div>
                <h1 className="text-2xl font-bold">Datos a verificar</h1>
                <p className="mt-2 text-sm text-gray-600">
                  Verificá que los datos ingresados abajo coincidan con el
                  documento subido por el vecino.
                </p>
                <p className="mt-2 text-sm text-gray-600">
                  Si encontrás diferencias, dejá una observación para que el
                  vecino corrija la información o cargue un nuevo documento
                  según corresponda.
                </p>
              </div>

              {verification.documents?.length > 1 && (
                <div className="flex items-center gap-2 mb-4 mx-20">
                  {verification.documents.map((doc, index) => (
                    <Fragment key={doc.id}>
                      <div
                        onClick={() =>
                          index < currentDocIndex && setCurrentDocIndex(index)
                        }
                        className={cn(
                          "flex items-center justify-center w-8 h-8 rounded-full border-2 font-medium",
                          index === currentDocIndex
                            ? "border-[#3b5c3f] text-[#3b5c3f]"
                            : index < currentDocIndex
                            ? "border-[#3b5c3f] bg-[#3b5c3f] text-white cursor-pointer"
                            : "border-gray-200 text-gray-400"
                        )}
                      >
                        {index + 1}
                      </div>
                      {index < verification.documents.length - 1 && (
                        <div className="flex-1 h-px bg-gray-200" />
                      )}
                    </Fragment>
                  ))}
                </div>
              )}

              <div className="space-y-2">
                {successData &&
                  verification.documents?.[currentDocIndex].sections.map(
                    (section) =>
                      section.show && (
                        <ValidationCard
                          key={section.id}
                          section={section}
                          documentContent={
                            verification.documents?.[currentDocIndex].request
                          }
                          onUpdate={(updates, updateNext) =>
                            handleSectionUpdate(section.id, updates, updateNext)
                          }
                          deleteSteps={(index) =>
                            setVerification((prev: Verification) => {
                              prev.documents = prev.documents.filter(
                                (_verification, idx) => idx != index
                              );
                              return prev;
                            })
                          }
                        ></ValidationCard>
                      )
                  )}
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                className="w-full mt-4"
                onClick={() => navigate(-1)}
              >
                Volver
              </Button>
              <Button
                className="w-full bg-[#3b5c3f] hover:bg-[#2e4831] text-white mt-4"
                disabled={
                  !verification.documents?.[currentDocIndex]?.sections[0]
                    .isConfirmed &&
                  !verification.documents?.[currentDocIndex]?.sections.every(
                    (s) => s.status == "valid" || s.isConfirmed
                  )
                }
                onClick={handleContinue}
              >
                Continuar
              </Button>
            </div>
          </div>
        </div>

        {/* Right Side - PDF Viewer */}
        <div className="flex-1 flex flex-col">
          <div className="bg-white border-b p-4">
            <h1 className="text-xl">
              {verification.documents?.[currentDocIndex]?.title}
            </h1>
            <div className="text-sm text-gray-600">
              {verification.documents?.[currentDocIndex]?.gedoCode}
            </div>
          </div>
          {verification.documents?.[currentDocIndex]?.request ? (
            <div className="p-6 space-y-6">
              <div>
                <h2 className="text-sm text-gray-500 mb-1">
                  Memoria descriptiva:
                </h2>
                <p className="text-sm">
                  {
                    verification.documents?.[currentDocIndex].request
                      .descriptiveMemory
                  }
                </p>
              </div>

              <div>
                <h2 className="text-sm text-gray-500 mb-1">Tarea asignada:</h2>
                <p
                  className="text-sm"
                  dangerouslySetInnerHTML={{
                    __html: verification.documents?.[
                      currentDocIndex
                    ].request.assignedTask
                      .map((task) => task)
                      .join("<br/>"),
                  }}
                ></p>
              </div>

              <div>
                <h2 className="text-sm text-gray-500 mb-1">
                  Tiempo estimado de ejecución:
                </h2>
                <p className="text-sm">
                  {
                    verification.documents?.[currentDocIndex].request
                      .estimatedTime
                  }
                </p>
              </div>

              <div>
                <h2 className="text-sm text-gray-500 mb-1">Requiere seguro:</h2>
                <p className="text-sm">
                  {
                    verification.documents?.[currentDocIndex].request
                      .requiresInsurance
                  }
                </p>
              </div>
            </div>
          ) : (
            verification.documents?.[currentDocIndex]?.base64 && (
              <div className="flex-1 bg-gray-100">
                <object
                  id="pdf"
                  height="100%"
                  width="100%"
                  type="application/pdf"
                  data={`data:application/pdf;base64, ${verification.documents?.[currentDocIndex].base64}`}
                >
                  <span>PDF plugin is not available.</span>
                </object>
              </div>
            )
          )}
        </div>
      </div>
    </>
  );
}
