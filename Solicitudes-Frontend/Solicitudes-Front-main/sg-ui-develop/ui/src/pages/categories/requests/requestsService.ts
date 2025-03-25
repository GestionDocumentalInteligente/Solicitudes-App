import { api } from "../../apiInstance";
import {
  AblResponse,
  Address,
  FilesType,
  RequestDetails,
  RequestDetailsResponse,
  RequestInfo,
  RequestInfoResponse,
  RequestResponse,
  UserResponse,
  ValidationError,
} from "./types";
import { ActivitiesResponse, RequestError, Response } from "./types";

type Payload = {
  address: Address | null;
  ablDebt: boolean | null;
  commonZone: boolean;
  userType: string;
  selectedActivities: string[];
  projectDescription: string;
  estimatedTime: number;
  insurance: boolean | null;
  files: FilesType;
};

export class RequestService {
  static async getAddressInfo(query: string): Promise<Response> {
    try {
      return await api.get<Response>("/admin/address?q=" + query);
    } catch (error) {
      throw new RequestError("Error en la busqueda de direcciones.");
    }
  }

  static async getSiteActivities(): Promise<ActivitiesResponse> {
    try {
      return await api.get<ActivitiesResponse>("/admin/site/activities");
    } catch (error) {
      throw new RequestError("Error en la busqueda de actividades.");
    }
  }

  static async getRequests(): Promise<RequestInfo[]> {
    try {
      const result = await api.get<RequestInfoResponse>("/admin/requests");
      return result.data.requests;
    } catch (error) {
      throw new RequestError("Error en la busqueda de solicitudes.");
    }
  }

  static async getUser(cuit: string): Promise<UserResponse> {
    try {
      return await api.get<UserResponse>(
        "/admin/users?cuit=" + cuit + "&provider=jwt"
      );
    } catch (error) {
      throw new RequestError("Error en la busqueda del usuario.");
    }
  }

  static async checkABL(abl: number): Promise<AblResponse> {
    try {
      return await api.get<AblResponse>("/admin/address/debt?abl=" + abl);
    } catch (error) {
      throw new RequestError("Error en la busqueda del usuario.");
    }
  }

  static async getRequestData(recordNumber: string): Promise<RequestDetails> {
    try {
      const result = await api.get<RequestDetailsResponse>(
        "/admin/requests/" + recordNumber
      );
      return result.data.request;
    } catch (error) {
      throw new RequestError("Error en la busqueda de la solicitud.");
    }
  }

  static async saveRequest(payload: Payload): Promise<RequestResponse> {
    try {
      if (!payload.address) {
        throw new ValidationError("El campo 'address' es requerido.");
      }
      if (!payload.userType) {
        throw new ValidationError("El campo 'userType' es requerido.");
      }
      if (
        !payload.commonZone &&
        (payload.ablDebt === null || payload.ablDebt)
      ) {
        throw new ValidationError("Debe verificar la deuda.");
      }

      if (payload.insurance === null) {
        throw new ValidationError("Debe marcar si tiene seguro o no.");
      }

      const formData = new FormData();

      formData.append("address_street", payload.address.address_street || "");
      formData.append(
        "address_number",
        payload.address.address_number.toString() || ""
      );
      formData.append(
        "address_abl_number",
        payload.address.abl_number.toString() || ""
      );
      formData.append(
        "property_id",
        payload.address.property_id.toString() || ""
      );
      formData.append("ablDebt", payload.ablDebt?.toString() || "");
      formData.append("commonZone", payload.commonZone.toString());
      formData.append("userType", payload.userType);
      formData.append(
        "selectedActivity",
        `[${payload.selectedActivities.join(",")}]`
      );
      formData.append("projectDescription", payload.projectDescription);
      formData.append("estimatedTime", payload.estimatedTime.toString());
      formData.append("insurance", payload.insurance?.toString() || "");

      Object.entries(payload.files).forEach(([key, file]) => {
        formData.append(`files[${key}]`, file.content, file.name);
        formData.append(`fileTypes[${key}]`, key);
      });

      return await api.post("/admin/requests", formData);
    } catch (error) {
      if (error instanceof ValidationError) {
        throw error;
      } else {
        console.error("Error en la solicitud al API:", error);
        throw new RequestError("Error en la creación de la solicitud.");
      }
    }
  }

  static async updateRequest(
    payload: Payload,
    id: string
  ): Promise<RequestResponse> {
    try {
      if (!payload.address) {
        throw new ValidationError("El campo 'address' es requerido.");
      }
      if (!payload.userType) {
        throw new ValidationError("El campo 'userType' es requerido.");
      }
      if (
        !payload.commonZone &&
        (payload.ablDebt === null || payload.ablDebt)
      ) {
        throw new ValidationError("Debe verificar la deuda.");
      }

      if (payload.insurance === null) {
        throw new ValidationError("Debe marcar si tiene seguro o no.");
      }

      const formData = new FormData();

      formData.append("address_street", payload.address.address_street || "");
      formData.append(
        "address_number",
        payload.address.address_number.toString() || ""
      );
      formData.append(
        "address_abl_number",
        payload.address.abl_number.toString() || ""
      );
      formData.append(
        "property_id",
        payload.address.property_id.toString() || ""
      );
      formData.append("ablDebt", payload.ablDebt?.toString() || "");
      formData.append("commonZone", payload.commonZone.toString());
      formData.append("userType", payload.userType);
      formData.append(
        "selectedActivity",
        `[${payload.selectedActivities.join(",")}]`
      );
      formData.append("projectDescription", payload.projectDescription);
      formData.append("estimatedTime", payload.estimatedTime.toString());
      formData.append("insurance", payload.insurance?.toString() || "");

      Object.entries(payload.files).forEach(([key, file]) => {
        formData.append(`files[${key}]`, file.content, file.name);
        formData.append(`fileTypes[${key}]`, key);
      });

      return await api.put("/admin/requests/" + id, formData);
    } catch (error) {
      if (error instanceof ValidationError) {
        throw error;
      } else {
        console.error("Error en la solicitud al API:", error);
        throw new RequestError("Error en la creación de la solicitud.");
      }
    }
  }
}
