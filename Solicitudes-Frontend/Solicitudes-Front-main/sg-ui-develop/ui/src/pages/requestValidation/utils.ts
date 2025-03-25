import { getRequestOptions } from "@/lib/utils.ts";
import {
  RequestResponse,
  RequestError,
  Validation,
} from "@/pages/requestValidation/types.ts";

export const getRequestsValidation = async (): Promise<Validation[]> => {
  const response = await fetch(
    "/api/admin/requests/validations",
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

export const getRequestValidationData = async (
  recordNumber: string
): Promise<Validation> => {
  const response = await fetch(
    `/api/admin/requests/validations/${recordNumber}`,
    getRequestOptions("GET", "")
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  // TODO: CHECK RESULT TYPES
  const result = await response.json();
  if (response.status >= 200 && response.status < 300) {
    return result.data;
  } else {
    const { error } = result;
    throw new RequestError(error);
  }
};

// TODO: MOVE TO VERIFICATION? OR GENERAL FILE
export const getRequestDocument = async (
  documentCode: string
): Promise<any> => {
  const response = await fetch(
    `/api/admin/requests/document/${documentCode}`,
    getRequestOptions("GET", "")
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  const result = await response.json();
  if (response.status >= 200 && response.status < 300) {
    return result.data;
  } else {
    const { error } = result;
    throw new RequestError(error);
  }
};

export const getDocumentByGedoCode = async (code: string): Promise<string> => {
  const response = await fetch(
    `/api/admin/requests/document/${code}`,
    getRequestOptions("GET", "")
  );
  if (!response.ok) {
    throw new Error("La solicitud no tuvo éxito.");
  }

  const { data } = await response.json();

  if (response.status >= 200 && response.status < 300) {
    return data.content;
  } else {
    const { error } = data;
    throw new RequestError(error);
  }
};

export const putRequestValidation = async (
  recordNumber: string,
  body: string
): Promise<Validation[]> => {
  const response = await fetch(
    `/api/admin/requests/validations/${recordNumber}`,
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

export const docsValidationParser = (validation: Validation): Validation => {
  const {
    documents: [firstDoc, secondDoc, thirdDoc, fourthDoc],
  } = validation;

  validation.documents[0] = {
    ...firstDoc,
    sections: [
      {
        id: "property",
        title: "Potestad sobre el inmueble",
        status: "valid",
        question: `
                    Verificado por: ${secondDoc.verifiedBy} <br/>
                    Fecha de verificación: ${secondDoc.verifiedDate} <br/>
                    Presentación: <a class="text-xs underline cursor-pointer" data-doc="0">${firstDoc.gedoCode}</a> <br/>
                    Verificación: <a class="text-xs underline cursor-pointer" data-doc="1">${secondDoc.gedoCode}</a>
                `,
        isExpanded: false,
        show: true,
      },
      {
        id: "tasks",
        title: "Tareas a realizar",
        status: "valid",
        question: `
                    Verificado por: ${fourthDoc.verifiedBy} <br/>
                    Fecha de verificación: ${fourthDoc.verifiedDate} <br/>
                    Presentación: <a class="text-xs underline cursor-pointer" data-doc="2">${thirdDoc.gedoCode}</a> <br/>
                    Verificación: <a class="text-xs underline cursor-pointer" data-doc="3">${fourthDoc.gedoCode}</a>
                `,
        isExpanded: false,
        show: true,
      },
    ],
  };
  return validation;
};
