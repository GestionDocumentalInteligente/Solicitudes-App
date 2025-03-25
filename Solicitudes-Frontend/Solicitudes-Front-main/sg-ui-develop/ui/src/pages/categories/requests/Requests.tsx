import React, { useEffect, useState } from "react";

import { RequestDetails, RequestInfo } from "./types";
import { RequestService } from "./requestsService";
import Home from "./Home";
import CustomButton from "../../../components/Buttons/CustomButton";
import { useNavigate } from "react-router-dom";
import BasicTable from "@/components/Table/BasicTable";
import LoadingScreen from "@/components/LoadingScreen/LoadingScreen";

const headers = [
  {
    id: "file_number",
    title: "N° de solicitud",
  },
  {
    id: "type",
    title: "Tipo de solicitud",
  },
  {
    id: "created_at",
    title: "Fecha de envío",
  },
  {
    id: "status",
    title: "Estado",
  },
  {
    id: "actions",
    title: "",
  },
];

const Requests: React.FC = () => {
  const navigate = useNavigate();

  const [processing, setProcessing] = useState(true);
  const [openHome, setOpenHome] = useState(false);
  const [requests, setRequests] = useState<RequestInfo[] | null>(null);

  const handleInit = () => {
    sessionStorage.clear();
    setOpenHome(true);
  };

  useEffect(() => {
    const getData = async () => {
      setProcessing(true);
      try {
        const result = await RequestService.getRequests();
        setRequests(result);
      } catch (e) {
        console.error(e);
        setRequests(null);
      }
      setProcessing(false);
    };
    getData();
  }, []);

  if (processing) {
    return <LoadingScreen title={["Cargando..."]} />;
  }

  if (openHome) {
    return <Home />;
  }

  return (
    <>
      {requests === null || requests.length === 0 ? (
        <Home />
      ) : (
        <div className="p-6">
          <div className="flex justify-between">
            <h1 className="text-2xl font-bold mb-6">Mis solicitudes</h1>
            <div>
              <CustomButton
                onClick={handleInit}
                className="bg-green-700 text-white"
              >
                + Nueva solicitud
              </CustomButton>
            </div>
          </div>
          <div className="h-px bg-gray-200 mb-6"></div>
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <BasicTable
              headers={headers}
              data={requests}
              actions={[
                {
                  label: "Ver",
                  displayStatus: ["PendingTask"],
                  callback: (request: RequestDetails) => {
                    navigate("/admin/requests/new", {
                      state: request,
                    });
                  },
                },
                {
                  label: "Ver",
                  displayStatus: ["finished"],
                  callback: handleInit,
                },
              ]}
            ></BasicTable>
          </div>
        </div>
      )}
    </>
  );
};

export default Requests;
