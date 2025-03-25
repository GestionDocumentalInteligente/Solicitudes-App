import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { HiInformationCircle } from "react-icons/hi";
import { Alert } from "flowbite-react";
import { AuthProvider, useAuth } from "../provider/useAuth";
import LoginButton from "../../../components/Buttons/LoginButton";
import { getUser } from "../provider/useLocalStorage";
import LoadingScreen from "@/components/LoadingScreen/LoadingScreen";

const redirectUri = import.meta.env.VITE_REDIRECT_URI as string;

const authUrl = import.meta.env.VITE_AUTH_URL as string;
const autenticarURL = `${authUrl}/auth/realms`;

const protocol = "protocol/openid-connect/auth?response_type=code";

const clientId = import.meta.env.VITE_CLIENT_ID as string;

const realmAfip = import.meta.env.VITE_REALM_AFIP as string;
const realmAnses = import.meta.env.VITE_REALM_ANSES as string;
const realmMiArg = import.meta.env.VITE_REALM_MIARG as string;

const loginProviders = [
  {
    name: "AFIP",
    url: `${autenticarURL}/${realmAfip}/${protocol}&client_id=${clientId}&redirect_uri=${redirectUri}?provider=AFIP`,
  },
  {
    name: "ANSES",
    url: `${autenticarURL}/${realmAnses}/${protocol}&client_id=${clientId}&redirect_uri=${redirectUri}?provider=ANSES`,
  },
  {
    name: "MI ARGENTINA",
    url: `${autenticarURL}/${realmMiArg}/${protocol}&client_id=${clientId}&redirect_uri=${redirectUri}?provider=MI_ARGENTINA`,
  },
];

const Login: React.FC = () => {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>(
    location.state?.errorMessage || ""
  );

  const handleLoginRedirect = (url: string) => {
    auth.login(url);
  };

  useEffect(() => {
    if (getUser()) {
      navigate("/admin/requests");
    }
  }, [navigate]);

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const authCode = params.get("code");
    const provider = params.get("provider");

    if (authCode && provider) {
      setIsLoading(true);
      const authenticateUser = async (authCode: string, provider: string) => {
        try {
          await auth.handleExternalAuth(authCode, provider, redirectUri);
        } catch (error) {
          const errorMessage =
            error instanceof Error ? error.message : "Error inesperado";
          setError(errorMessage);
        }
      };
      authenticateUser(authCode, provider);
    }
  }, [location.search, auth]);

  if (isLoading) {
    return <LoadingScreen title={["Ingresando"]} />;
  }

  return (
    <div className="flex flex-col">
      <div className="flex flex-col mx-auto md:py-0">
        <div className="md:space-y-4 sm:p-4">
          <h1 className="text-3xl mb-4 font-semibold leading-tight tracking-tight md:text-4xl">
            ¡Hola!
          </h1>
          <div className="space-y-4 md:space-y-4 max-w-md">
            <p className="font-sans text-base">
              Bienvenido/a a la gestión de solicitudes online de la
              Municipalidad de San Isidro
            </p>
            {loginProviders.map((provider) => (
              <LoginButton
                key={provider.name}
                name={provider.name}
                onClick={() => handleLoginRedirect(provider.url)}
              />
            ))}
            {error !== "" ? (
              <Alert color="failure" icon={HiInformationCircle}>
                <span className="font-medium">Error!</span> Ocurrió un error al
                intentar iniciar sesión. Por favor, intente nuevamente.
              </Alert>
            ) : (
              ""
            )}
          </div>
        </div>
        <p className="font-sans max-w-lg text-xs text-center w-full">
          Al usar este servicio, aceptás nuestros{" "}
          <a
            href="https://docs.google.com/document/d/1U33075GjuApU84zV5v5wTM5W1-dtD8y-ijj15i-ZQE4/preview"
            className="underline"
            target="_blank"
          >
            términos y condiciones.
          </a>
        </p>
      </div>
    </div>
  );
};

const SignInPage = () => (
  <AuthProvider>
    <Login />
  </AuthProvider>
);

export default SignInPage;
