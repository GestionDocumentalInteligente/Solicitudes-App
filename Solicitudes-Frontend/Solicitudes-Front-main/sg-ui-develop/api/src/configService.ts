import dotenv from "dotenv";

dotenv.config();

class ConfigService {
  public readonly baseUsersApi = process.env.BASE_USERS_API as string;
  public readonly baseLoginApi = process.env.BASE_LOGIN_API as string;
  public readonly baseMailingApi = process.env.BASE_MAILING_API as string;
  public readonly baseRequestApi = process.env.BASE_REQUEST_API as string;
  public readonly baseFileManagerApi = process.env.BASE_MANAGER_API as string;

  public readonly clientSecretAfip: string;
  public readonly clientSecretAnses: string;
  public readonly clientSecretMiArg: string;
  public readonly clientId: string;
  public readonly realmAfip: string;
  public readonly realmAnses: string;
  public readonly realmMiArg: string;
  public readonly authUrl: string;

  constructor() {
    this.clientSecretAfip = process.env.CLIENT_SECRET_AFIP as string;
    this.clientSecretAnses = process.env.CLIENT_SECRET_ANSES as string;
    this.clientSecretMiArg = process.env.CLIENT_SECRET_MIARG as string;
    this.clientId = process.env.CLIENT_ID as string;
    this.realmAfip = process.env.REALM_AFIP as string;
    this.realmAnses = process.env.REALM_ANSES as string;
    this.realmMiArg = process.env.REALM_MIARG as string;
    this.authUrl = process.env.AUTH_URL as string;
  }
}

export const configService = new ConfigService();
