import { useEffect, useState } from "react";
import { Alert } from "flowbite-react";
import { useNavigate } from "react-router-dom";
import { HiInformationCircle } from "react-icons/hi";

import { AuthProvider, useAuth } from "../provider/useAuth";
import { ContactInfo } from "./ContactInfo";
import { ProfileInfo } from "./ProfileInfo";
import { getUser } from "../provider/useLocalStorage";
import { User } from "../types";
import LoadingScreen from "@/components/LoadingScreen/LoadingScreen";

const Profile = () => {
  const [step, setStep] = useState(1);
  const navigate = useNavigate();
  const [processing, setProcessing] = useState(true);
  const auth = useAuth();
  const [userInfo, setUserInfo] = useState<User | null>(null);

  const [error, setError] = useState("");

  useEffect(() => {
    if (step === 3) {
      navigate("/email-verification-sent", {
        state: { email: userInfo?.email, type: 1 },
      });
    }
  }, [step, navigate, userInfo]);

  useEffect(() => {
    const getData = async () => {
      if (!auth.isRegistrationRequired()) {
        navigate("/admin/requests");
        return;
      }
      setError("");
      setProcessing(true);
      const payload = auth.decodeToken();
      if (!payload) {
        setError("user not found");
        return;
      }

      const userData: User = {
        cuit: payload.cuit,
        dni: payload.cuit,
        first_name: payload.given_name,
        last_name: payload.family_name,
        email: "",
        phone: "",
      };
      setUserInfo(userData);
      setProcessing(false);
    };
    if (getUser()) {
      getData();
    } else {
      navigate("/login");
    }
  }, [navigate, auth]);

  const handleNextStep = () => {
    setStep(step + 1);
  };

  if (processing) {
    return <LoadingScreen title={["Validando información"]} />;
  }

  return (
    <div className="flex flex-col">
      <div className="h-full md:mr-14">
        <h1 className="text-xl font-semibold tracking-tight md:text-3xl">
          Tu perfil
        </h1>
        <p className="font-sans text-sm mt-4">
          {step === 1
            ? "Tus datos provienen de la plataforma con la que iniciaste sesión. Verificá que sean correctos:"
            : "Completá tus datos de contacto para continuar."}
        </p>
        {step === 1 && (
          <ProfileInfo
            userInfo={userInfo}
            onNext={handleNextStep}
            setError={setError}
            error={error}
          />
        )}
        {step === 2 && (
          <ContactInfo
            userInfo={userInfo}
            setUserInfo={setUserInfo}
            onNext={handleNextStep}
            setError={setError}
          />
        )}
        {error !== "" && (
          <Alert
            className="w-full mt-4"
            color="failure"
            icon={HiInformationCircle}
          >
            <span className="font-medium">Error!</span> {error}
          </Alert>
        )}
      </div>
    </div>
  );
};

const CompleteRegistration = () => (
  <AuthProvider>
    <Profile />
  </AuthProvider>
);

export default CompleteRegistration;
