import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { RequestError } from "@/pages/categories/requests/types.ts";
import BasicTable from "@/components/Table/BasicTable.tsx";
import { useLoadingScreenStore } from "@/stores/loadingScreenStore.ts";
import { Validation } from "@/pages/requestValidation/types.ts";
import { getRequestsValidation } from "@/pages/requestValidation/utils.ts";

const requestTypeOptions = ["Aviso de obra"];
const statusOptions = ["Pendiente", "Verificado", "Observado"];

export default function RequestValidation() {
  const [validations, setValidations] = useState<Validation[]>([]);
  const [selectedRequestType, setSelectedRequestType] = useState<string | null>(
    null
  );
  const [selectedStatus, setSelectedStatus] = useState<
    (typeof statusOptions)[number] | null
  >(null);
  const showLoading = useLoadingScreenStore((state) => state.showLoading);
  const hideLoading = useLoadingScreenStore((state) => state.hideLoading);
  // TODO: MAKE GLOBAL ERROR / USE STATUS MESSAGE
  const [, setError] = useState("");
  const navigate = useNavigate();

  const headers = [
    {
      id: "recordNumber",
      title: "N° Expediente",
    },
    {
      id: "requesterAddress",
      title: "Domicilio",
    },
    {
      id: "requestType",
      title: "Tipo de solicitud",
      filter: {
        typeOptions: requestTypeOptions,
        selectedType: selectedRequestType,
        setSelectedType: setSelectedRequestType,
      },
    },
    {
      id: "deliveryDate",
      title: "Fecha de solicitud",
    },
    {
      id: "status",
      title: "Estado",
      filter: {
        typeOptions: statusOptions,
        selectedType: selectedStatus,
        setSelectedType: setSelectedStatus,
      },
    },
    {
      id: "actions",
      title: "",
    },
  ];

  const getData = async () => {
    showLoading();
    try {
      const result = await getRequestsValidation();
      setValidations(
        result.sort((a, b) => {
          const dateA = new Date(a.deliveryDate);
          const dateB = new Date(b.deliveryDate);
          return dateB.getTime() - dateA.getTime();
        })
      );
    } catch (e) {
      if (e instanceof RequestError) {
        setError(e.message);
      } else {
        setError("Ocurrió un error en la busqueda de información.");
      }
    }
    hideLoading();
  };

  useEffect(() => {
    getData();
  }, []);

  const filteredVerifications: Validation[] = validations.filter(
    (verification) => {
      if (
        selectedRequestType &&
        verification.requestType !== selectedRequestType
      )
        return false;
      if (selectedStatus && verification.status !== selectedStatus)
        return false;
      return true;
    }
  );

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Solicitudes</h1>
      <div className="h-px bg-gray-200 mb-6"></div>
      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <BasicTable
          headers={headers}
          data={filteredVerifications}
          actions={[
            {
              label: "Validar",
              // TODO: GENERATE ENUM WITH STATUS TYPES
              displayStatus: ["pending"],
              callback: (validation: Validation) => {
                navigate("/admin-panel/data-validation", {
                  state: validation,
                });
              },
            },
          ]}
        ></BasicTable>
      </div>
    </div>
  );
}
