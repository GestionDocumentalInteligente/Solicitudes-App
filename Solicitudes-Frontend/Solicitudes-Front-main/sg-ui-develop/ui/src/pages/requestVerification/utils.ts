import { getRequestOptions } from "@/lib/utils.ts";
import {
  RequestResponse,
  RequestError,
  Verification,
  Document,
} from "@/pages/requestVerification/types.ts";

export const getRequestsVerification = async (): Promise<Verification[]> => {
  const response = await fetch(
    "/api/admin/requests/verifications",
    getRequestOptions("GET", "")
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  const result: RequestResponse = await response.json();
  if (response.status >= 200 && response.status < 300) {
    return result.data;
  } else {
    const { error } = result;
    throw new RequestError(error);
  }
};

export const getRequestVerificationData = async (
  recordNumber: string
): Promise<Document[]> => {
  const response = await fetch(
    `/api/admin/requests/verifications/${recordNumber}`,
    getRequestOptions("GET", "")
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  // TODO: CHECK RESULT TYPES
  const result = await response.json();
  if (response.status >= 200 && response.status < 300) {
    return result.data.documents;
  } else {
    const { error } = result;
    throw new RequestError(error);
  }
};

export const putRequestVerification = async (
  recordNumber: string,
  body: string
): Promise<Verification[]> => {
  const response = await fetch(
    `/api/admin/requests/verifications/${recordNumber}`,
    getRequestOptions("PUT", body)
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  const result: RequestResponse = await response.json();
  if (response.status >= 200 && response.status < 300) {
    return result.data;
  } else {
    const { error } = result;
    throw new RequestError(error);
  }
};

export const detachRequest = (
  verifications: Verification[]
): Verification[] => {
  const verificationClone = structuredClone(verifications);
  return [
    ...verifications.map((verification: Verification) => {
      verification.status = verification.status_property;
      verification.documentType = "Potestad";
      return verification;
    }),
    ...verificationClone.map((verification: Verification) => {
      verification.status = verification.status_task;
      verification.documentType = "Tareas";
      return verification;
    }),
  ];
};

const defineCase = (verification: Verification) => {
  const verificationDocuments: Document[] = [];
  const propertyTitle = verification.documents.find(
    (doc) =>
      doc.title.trim().toLowerCase() ===
      "título de propiedad o informe de dominio"
  );
  const propertyAuthorization = verification.documents.find(
    (doc) => doc.title.trim().toLowerCase() === "autorización del propietario"
  );
  const propertyRegulation = verification.documents.find(
    (doc) => doc.title.trim().toLowerCase() === "reglamento de co-propiedad"
  );
  const designationAct = verification.documents.find(
    (doc) => doc.title.trim().toLowerCase() === "acta de designación"
  );
  if (propertyTitle) {
    if (propertyAuthorization) {
      verification.invokedCharacter = "Ocupante o Inquilino";
    } else {
      verification.invokedCharacter = "Titular";
    }
  } else {
    if (propertyRegulation && designationAct) {
      verification.invokedCharacter = "Administrador";
    }
  }
  if (verification.documentType.trim().toLowerCase() === "tareas") {
    verification.verificationCase = "4";
    const documentContent = verification.documents.find((doc) => doc?.request);
    if (documentContent) {
      verificationDocuments[0] = documentContent;
      if (
        documentContent.request.requiresInsurance.trim().toLowerCase() === "si"
      ) {
        verificationDocuments[1] = <Document>(
          verification.documents.find(
            (doc) => doc.title.trim().toLowerCase() === "póliza de seguro"
          )
        );
      }
    }
  } else {
    if (propertyTitle) {
      if (propertyAuthorization) {
        verification.verificationCase = "3";
        verificationDocuments[0] = propertyAuthorization;
        verificationDocuments[1] = propertyTitle;
      } else {
        verification.verificationCase = "1";
        verificationDocuments[0] = propertyTitle;
      }
    } else {
      if (propertyRegulation && designationAct) {
        verification.verificationCase = "2";
        verificationDocuments[0] = propertyRegulation;
        verificationDocuments[1] = designationAct;
      }
    }
  }
  verification.documents = verificationDocuments;
};

export const docsVerificationParser = (
  verification: Verification
): Verification => {
  defineCase(verification);
  const {
    documents: [firstDoc, secondDoc],
  } = verification;
  switch (verification.verificationCase) {
    case "1":
      firstDoc.sections = [
        {
          id: "document",
          title: "Documentación",
          status: null,
          question:
            "¿El documento cuenta con las características propias de un titulo de propiedad o Informe de dominio?",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no cumple con las características necesarias de un título de propiedad porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: true,
        },
        {
          id: "address",
          title: "Dirección",
          status: null,
          question: verification.requesterAddress,
          observationPlaceholder:
            'Ejemplo: "La dirección ingresada no coincide con la dirección del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
        {
          id: "fullName",
          title: "Nombre y Apellido",
          status: null,
          question: verification.requesterFullName,
          observationPlaceholder:
            'Ejemplo: "El nombre y apellido ingresado no coincide con el nombre y apellido del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
      ];
      break;
    case "2":
      firstDoc.sections = [
        {
          id: "document",
          title: "Documentación",
          status: null,
          question:
            "¿El documento presentado cuenta con las características propias de un Reglamento de copropiedad?",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no cumple con las características necesarias de un reglamento de copropiedad porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: true,
        },
        {
          id: "address",
          title: "Dirección",
          status: null,
          question: verification.requesterAddress,
          observationPlaceholder:
            'Ejemplo: "La dirección ingresada no coincide con la dirección del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
      ];
      secondDoc.sections = [
        {
          id: "document",
          title: "Documentación",
          status: null,
          question:
            "¿El documento presentado cuenta con las características propias de un Acta de Designación?*",
          questionClarification:
            "*Si el Reglamento de Copropiedad establece la designación del administrador y el mandato se encuentra vigente, dicho documento sería suficiente para acreditar su rol.",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no cumple con las características necesarias de un reglamento de copropiedad porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: true,
        },
        {
          id: "address",
          title: "Dirección",
          status: null,
          question: verification.requesterAddress,
          observationPlaceholder:
            'Ejemplo: "La dirección ingresada no coincide con la dirección del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
        {
          id: "fullName",
          title: "Nombre y Apellido",
          status: null,
          question: verification.requesterFullName,
          observationPlaceholder:
            'Ejemplo: "El nombre y apellido ingresado no coincide con el nombre y apellido del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
        {
          id: "validity",
          title: "Vigencia",
          status: null,
          question: "¿Tiene vigencia al día de la fecha?",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no se encuentra vigente. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
      ];
      break;
    case "3":
      firstDoc.sections = [
        {
          id: "document",
          title: "Documentación",
          status: null,
          question:
            "¿El documento cuenta con las características propias de una Autorización del propietario?",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no cumple con las características necesarias de una autorización del propietario porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: true,
        },
        {
          id: "address",
          title: "Dirección",
          status: null,
          question: verification.requesterAddress,
          observationPlaceholder:
            'Ejemplo: "La dirección ingresada no coincide con la dirección del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
        {
          id: "grantedByAuth",
          title: "Autorización otorgada por",
          status: null,
          observationPlaceholder:
            'Ejemplo: "El nombre y apellido ingresado no coincide con el nombre y apellido del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          inputType: "text",
          inputPlaceholder:
            "Escribir el nombre y apellido de la persona que firma la autorización",
          isExpanded: false,
          show: false,
        },
        {
          id: "grantedToAuth",
          title: "Autorización otorgada a",
          status: null,
          question: verification.requesterFullName,
          observationPlaceholder:
            'Ejemplo: "El nombre y apellido ingresado no coincide con el nombre y apellido del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
        {
          id: "validity",
          title: "Vigencia",
          status: null,
          question: `¿La autorización se encuentra vigente al día de la fecha?`,
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no se encuentra vigente. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
      ];
      secondDoc.sections = [
        {
          id: "document",
          title: "Documentación",
          status: null,
          question:
            "¿El documento presentado reúne las características propias de un título de propiedad?",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no cumple con las características necesarias de un título de propiedad porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: true,
        },
        {
          id: "address",
          title: "Dirección",
          status: null,
          question: verification.requesterAddress,
          observationPlaceholder:
            'Ejemplo: "La dirección ingresada no coincide con la dirección del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
        {
          id: "fullName",
          title: "Nombre y Apellido",
          status: null,
          question: "Should be grantedByAuth name",
          observationPlaceholder:
            'Ejemplo: "El nombre y apellido ingresado no coincide con el nombre y apellido del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
          isExpanded: false,
          show: false,
        },
      ];
      break;
    case "4":
      firstDoc.sections = [
        {
          id: "assignedTask",
          title: "Tarea asignada",
          status: null,
          question: "¿La memoria descriptiva coincide con la tarea asignada?",
          observationOptions: [],
          selectedObservationOptions: [],
          isExpanded: false,
          show: true,
        },
        {
          id: "insurance",
          title: "Seguro",
          status: null,
          question: "¿La tarea seleccionada requiere seguro?",
          questionClarification:
            "En caso de requerir, es necesario que el responsable adjunte una póliza de seguro.",
          observationPlaceholder:
            'Ejemplo: "El documento adjuntado no cumple con las características necesarias de un reglamento de copropiedad porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
          inputType: "radio",
          inputOptions: ["Si", "No"],
          isExpanded: false,
          show: false,
        },
      ];
      if (secondDoc) {
        secondDoc.sections = [
          {
            id: "document",
            title: "Documentación",
            status: null,
            question:
              "¿El documento cuenta con las características propias de una Póliza de seguro?",
            observationPlaceholder:
              'Ejemplo: "El documento adjuntado no cumple con las características necesarias de una Póliza de seguro porque [especificar razón]. Por favor, revise y cargue el documento correspondiente. Gracias."',
            isExpanded: false,
            show: true,
          },
          {
            id: "address",
            title: "Dirección",
            status: null,
            question: verification.requesterAddress,
            observationPlaceholder:
              'Ejemplo: "La dirección ingresada no coincide con la dirección del documento adjuntado. Por favor, revise y cargue el documento correspondiente. Gracias."',
            isExpanded: false,
            show: false,
          },
          {
            id: "clause",
            title: "Cláusula",
            status: null,
            question:
              "¿El documento incluye una cláusula de no repetición contra la municipalidad?",
            observationPlaceholder:
              'Ejemplo: "El documento adjuntado no incluye una cláusula de no repetición contra la municipalidad. Por favor, revise y cargue el documento correspondiente. Gracias."',
            isExpanded: false,
            show: false,
          },
          {
            id: "validity",
            title: "Vigencia",
            status: null,
            question: `¿La vigencia de la póliza cubre el período de ${firstDoc.request.estimatedTime} declarados por el vecino?`,
            observationPlaceholder:
              'Ejemplo: "El documento adjuntado no se encuentra vigente. Por favor, revise y cargue el documento correspondiente. Gracias."',
            isExpanded: false,
            show: false,
          },
        ];
      }
  }
  return verification;
};
