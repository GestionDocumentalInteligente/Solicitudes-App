export interface RegistrationData {
  email: string;
  phone: string;
  accepts_notifications?: boolean;
  email_validated?: boolean;
}

export type User = {
  cuit: string;
  dni: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  email_validated?: boolean;
  accepts_notifications?: boolean;
};

export class RequestError extends Error {
  private status: number | undefined;

  constructor(status: number | undefined, msg: string) {
    super(msg);
    this.status = status;
    Object.setPrototypeOf(this, RequestError.prototype);
  }

  getStatus(): number | undefined {
    return this.status;
  }
}

export type TokenResponse = {
  success: boolean;
  data: {
    access_token: string;
    cuil: string;
  };
  error: string;
};

export type Token = {
  access_token: string;
  data: User;
};

export type AutenticarToken = {
  access_token: string;
  refresh_token: string;
};

export type AutenticarInfo = {
  url: string;
  realmId: string;
  clientId: string;
  clientSecret: string;
};

export type JwtPayload = {
  exp: number;
  nbf: number;
  iat: number;
  iss: string;
  aud: string;
  sub: string;
  typ: string;
  azp: string;
  auth_time: number;
  session_state: string;
  acr: string;
  allowed_origins: string[];
  realm_access: {
    roles: string[];
  };
  resource_access: {
    account: {
      roles: string[];
    };
  };
  cuit: string;
  tipo_persona: string;
  proveedor: string;
  preferred_username: string;
  given_name: string;
  nivel: string;
  family_name: string;
};

export type UserResponse = {
  success: boolean;
  admin: boolean;
  data: User;
  error: string;
};

export type ErrorResponse = {
  code: string;
  message?: string;
};
