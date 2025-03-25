export interface Request {
  id: number;
  name: string;
  description: string;
  category_id: number;
  requires_documentation: boolean;
  is_active: boolean;
}

export type RequestInfoResponse = {
  success: boolean;
  data: {
    requests: RequestInfo[];
  };
};

export interface RequestInfo {
  FileNumber: string;
  CreatedAt: string;
  status: string;
}

export type RequestDetailsResponse = {
  success: boolean;
  data: {
    request: RequestDetails;
  };
};

export interface RequestDetails {
  ID: number;
  UserID: number;
  PropertyID: number;
  Cuil: string;
  Dni: string;
  FirstName: string;
  LastName: string;
  Email: string;
  Phone: string;
  Address: {
    Street: string;
    Number: string;
    ABLNumber: number;
  };
  ABLDebt: string;
  CommonZone: boolean;
  UserType: string;
  SelectedActivities: string[];
  Activities: string[];
  ProjectDesc: string;
  EstimatedTime: number;
  Insurance: boolean;
  FileNumber: string;
  Documents: Array<{
    Name: string;
    Type: number;
    Content: string;
  }>;
  StatusName: string;
  CreatedAt: string;
  Observations: string;
  ObservationsTasks: string;
  VerifyBy: string;
  VerifyByTasks: string;
  VerifyDate: string;
  VerifyDateTask: string;
}

export interface Address {
  address_street: string;
  address_number: number;
  abl_number: number;
  property_id: number;
}

export interface Response {
  success: boolean;
  data: Suggestions;
}

export interface Suggestions {
  suggestions: Address[];
}

export interface Activities {
  id: string;
  description: string;
}

export interface ActivitiesResponse {
  success: boolean;
  data: Activities[];
}

export type User = {
  cuil: string;
  dni: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
};

export type UserResponse = {
  success: boolean;
  data: User;
};

export type AblResponse = {
  success: boolean;
  data: {
    dbt: boolean;
  };
};

export type RequestResponse = {
  success: boolean;
};

export const UserDocumentRequirements = {
  Admin: [9, 14],
  Owner: [10],
  Occupant: [10, 11],
};

type FileInfo = {
  content: File;
  name: string;
};

export type FilesType = Record<number, FileInfo>;

export const DocumentTypes = {
  CoOwnership: {
    id: 9,
    description: "co_ownership_regulation",
  },
  PropertyTitle: {
    id: 10,
    description: "property_title_or_ownership_report",
  },
  OwnerAuthorization: {
    id: 11,
    description: "owner_authorization",
  },
  Insurance: {
    id: 12,
    description: "insurance",
  },
  AppointmentCertificate: {
    id: 14,
    description: "appointment_certificate",
  },
};

class RequestError extends Error {}

class ValidationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ValidationError";
  }
}

export { RequestError, ValidationError };
