import { Card } from "flowbite-react";

import { Request } from "./types";
import { Card as CustomCard } from "../../../components/Card/Card";
import BuilderImg from "../../../assets/aviso_de_obra_banner-min.png";
import CustomButton from "../../../components/Buttons/CustomButton";
import { useNavigate } from "react-router-dom";

const requests: Request[] = [
  {
    id: 1,
    name: "Autorización y gestión de torres grúa",
    description: "",
    category_id: 1,
    requires_documentation: true,
    is_active: false,
  },
  {
    id: 2,
    name: "Permiso de demolición",
    description: "",
    category_id: 1,
    requires_documentation: true,
    is_active: false,
  },
  {
    id: 3,
    name: "Habilitación y gestión comercial o industrial",
    description: "",
    category_id: 1,
    requires_documentation: true,
    is_active: false,
  },
  {
    id: 4,
    name: "Instalación y gestión de medio de traslación",
    description: "",
    category_id: 1,
    requires_documentation: true,
    is_active: false,
  },
  {
    id: 5,
    name: "Permiso y gestión de obra",
    description: "",
    category_id: 1,
    requires_documentation: true,
    is_active: false,
  },
  {
    id: 6,
    name: "Extracción de forestación",
    description: "",
    category_id: 1,
    requires_documentation: true,
    is_active: false,
  },
];

const Home = () => {
  const navigate = useNavigate();

  const handleInit = () => {
    navigate("/admin/requests/new");
  };

  return (
    <div className="grid grid-cols-1 gap-4 max-w-5xl mx-auto mb-2 mt-2">
      <h1 className="text-2xl font-semibold tracking-tight md:text-4xl">
        Solicitudes
      </h1>
      <Card className="max-w-full relative overflow-hidden min-h-72">
        <img
          src={BuilderImg}
          alt="Imagen de Aviso de Obra"
          className="w-full h-full object-cover absolute inset-0"
        />
        <div className="absolute inset-0 flex flex-col justify-end p-4 h-full">
          <div className="w-full flex flex-col md:flex-row justify-between items-center md:items-end space-y-2 md:space-y-0">
            <span className="w-full md:w-1/2 justify-end text-3xl font-bold text-white">
              Aviso de obra
            </span>
            <div className="flex space-x-2">
              <CustomButton
                onClick={handleInit}
                className="bg-green-700 text-white"
              >
                Comenzar
              </CustomButton>
              <CustomButton
                color="light"
                onClick={() => {
                  window.open(
                    "https://docs.google.com/document/d/1q9daijaznqR8ngh4YK9lm_uxqZiw6tpirobkGG8tb8k/preview",
                    "_blank",
                    "noopener,noreferrer"
                  );
                }}
              >
                Más info
              </CustomButton>
            </div>
          </div>
        </div>
      </Card>
      <hr />
      <p className="text-gray-600 text-sm">Próximamente</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-w-5xl mx-auto mt-2">
        {requests?.map((request) => (
          <CustomCard key={request.id} info={request} icon={null} to="" />
        ))}
      </div>
    </div>
  );
};

export default Home;
