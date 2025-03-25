import { Button } from "flowbite-react";
import React from "react";

import AfipLogo from "../../assets/afip.svg";
import AnsesLogo from "../../assets/anses.svg";
import MiArgLogo from "../../assets/miarg.svg";

const iconMap: { [key: string]: JSX.Element } = {
  AFIP: <img src={AfipLogo} alt="AFIP logo" className="h-8" />,
  ANSES: <img src={AnsesLogo} alt="AFIP logo" className="h-6" />,
  "MI ARGENTINA": <img src={MiArgLogo} alt="AFIP logo" className="h-8" />,
};

interface CustomButtonProps {
  name: string;
  onClick?: () => void;
}

const LoginButton: React.FC<CustomButtonProps> = ({ onClick, name }) => {
  return (
    <Button
      type="button"
      className="flex items-center justify-center w-full bg-title max-w-lg h-12 hover:bg-primary-dark"
      onClick={onClick}
    >
      <span className="text-base mr-2">Iniciar con</span>
      {iconMap[name]}
    </Button>
  );
};

export default LoginButton;
