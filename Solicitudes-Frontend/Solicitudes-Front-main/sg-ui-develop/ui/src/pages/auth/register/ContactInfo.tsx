import { Label, TextInput } from "flowbite-react";
import { HiMail } from "react-icons/hi";

import CustomButton from "../../../components/Buttons/CustomButton";
import { useEffect, useState } from "react";
import PhoneInput from "../../../components/Input/PhoneInput";
import { useAuth } from "../provider/useAuth";
import { RequestError, User } from "../types";

interface ContactInfoProps {
  userInfo: User | null;
  onNext: () => void;
  setError: (value: string) => void;
  setUserInfo: (user: User | null) => void;
}

export const ContactInfo = ({
  userInfo,
  onNext,
  setError,
  setUserInfo,
}: ContactInfoProps) => {
  const auth = useAuth();
  const [userData, setUserData] = useState<User>({
    cuit: "",
    dni: "",
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    accepts_notifications: false,
  });
  const [phone, setPhone] = useState<string>("");
  const [email, setEmail] = useState("");
  const [acceptNotifications, setAcceptNotifications] = useState(false);
  const [emailSending, setEmailSending] = useState(false);
  const [emailValid, setEmailValid] = useState(true);

  useEffect(() => {
    if (userInfo) {
      setUserData(userInfo);
    }
  }, [userInfo]);

  const validateEmail = () => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const isValid = validateEmail();
    setEmailValid(isValid);
    if (isValid && phone !== "") {
      const updatedUserData: User = {
        ...userData,
        email,
        phone,
        accepts_notifications: acceptNotifications,
      };
      setUserData(updatedUserData);
      setUserInfo(updatedUserData);
      setEmailSending(true);
      try {
        await auth.completeRegistration(updatedUserData);
        setError("");
        onNext();
      } catch (error) {
        setEmailSending(false);
        if (error instanceof RequestError) {
          setError(error.message);
          return;
        }
        setError("Error al crear el usuario");
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="w-full mt-4">
        <Label
          htmlFor="email1"
          className="font-semibold"
          value="Correo electrónico"
        />
        <div className="w-full mt-1">
          <TextInput
            id="email1"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            icon={HiMail}
            placeholder="Ingresa tu e-mail"
            disabled={emailSending}
            required
            color={emailValid ? "" : "failure"}
            helperText={
              !emailValid && (
                <>
                  <span className="font-medium">Formato inválido!</span>{" "}
                  Asegúrate de que el formato sea de un email válido.
                </>
              )
            }
          />
        </div>
      </div>
      <div className="w-full mt-4">
        <Label
          htmlFor="phone"
          className="font-semibold"
          value="Numero de teléfono (móvil)"
        />
        <div className="w-full mt-1">
          <PhoneInput phone={phone} setPhone={setPhone} />
        </div>
      </div>
      <div className="flex items-center mb-4 mt-4">
        <input
          id="default-checkbox"
          type="checkbox"
          checked={acceptNotifications}
          onChange={(event) => setAcceptNotifications(event.target.checked)}
          className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 focus:ring-2"
        />
        <label
          htmlFor="default-checkbox"
          className="ms-2 text-sm font-medium text-gray-900 dark:text-gray-300"
        >
          Acepto recibir notificaciones sobre el avance de mi solicitud por
          email o celular.
        </label>
      </div>
      <div className="w-full flex flex-col justify-center items-center">
        <CustomButton
          type="submit"
          isLoading={emailSending}
          disabled={emailSending}
          className="bg-title w-full"
        >
          Continuar
        </CustomButton>
      </div>
    </form>
  );
};
