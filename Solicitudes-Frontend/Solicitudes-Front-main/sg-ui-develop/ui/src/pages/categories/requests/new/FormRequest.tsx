import React, { useEffect, useState } from "react";
import { Stepper, Step, Typography } from "@material-tailwind/react";

import CustomButton from "../../../../components/Buttons/CustomButton";
import { StepOne } from "./StepOne";
import { StepTwo } from "./StepTwo";
import { useLocation, useNavigate } from "react-router-dom";
import { StepThree } from "./StepThree";
import {
  Address,
  FilesType,
  RequestError,
  RequestResponse,
  ValidationError,
} from "../types";
import { StepResume } from "./StepResume";
import { RequestService } from "../requestsService";
import { useLoadingScreenStore } from "@/stores/loadingScreenStore.ts";
import { base64ToFile } from "@/lib/utils.ts";

const FormRequest: React.FC = () => {
  const navigate = useNavigate();

  const [showStepper, setShowStepper] = useState(true);

  const [address, setAdress] = useState<Address | null>(() => {
    const savedAddress = sessionStorage.getItem("address");
    return savedAddress ? (JSON.parse(savedAddress) as Address) : null;
  });

  const [ablDebt, setAblDebt] = useState<boolean | null>(() => {
    const debt = sessionStorage.getItem("ablDebt");
    return debt === "true" ? true : debt === "false" ? false : null;
  });

  const [commonZone, setCommonZone] = useState<boolean>(() => {
    const zone = sessionStorage.getItem("commonZone");
    return zone === "true" ? true : false;
  });

  const [userType, setUserType] = useState<string>(() => {
    const userType = sessionStorage.getItem("userType");
    return userType ? userType : "";
  });

  const [selectedActivities, setSelectedActivities] = useState<string[]>(() => {
    const activities: string[] = JSON.parse(
      sessionStorage.getItem("selectedActivities") as string
    );
    return activities
      ? activities.map((activity: string) => activity.toString())
      : [];
  });

  const [projectDescription, setProjectDescription] = useState<string>(() => {
    const pd = sessionStorage.getItem("projectDescription");
    return pd ? pd : "";
  });

  const [estimatedTime, setEstimatedTime] = useState<number>(() => {
    const time = sessionStorage.getItem("estimatedTime");
    return time ? Number(time) : 0;
  });

  const [insurance, setInsurance] = useState<boolean | null>(() => {
    const ins = sessionStorage.getItem("insurance");
    return ins === "true" ? true : ins === "false" ? false : null;
  });

  const [observations, setObservations] = useState<string>(() => {
    return sessionStorage.getItem("observations") || "";
  });

  const [tasksObservations, setTasksObservations] = useState<string>(() => {
    return sessionStorage.getItem("tasksObservations") || "";
  });

  const [files, setFiles] = useState<FilesType>({});

  const [disabled, setDisabled] = useState<boolean>(true);
  const [update, setUpdate] = useState<boolean>(false);

  const [activeStep, setActiveStep] = useState(() => {
    const savedStep = sessionStorage.getItem("currentStep");
    return savedStep ? parseInt(savedStep) : 0;
  });

  const [isLastStep, setIsLastStep] = React.useState(false);
  const [isFirstStep, setIsFirstStep] = React.useState(false);
  const [creatingRecord, setCreatingRecord] = useState(false);

  const { state: currentRequest } = useLocation();
  const [, setError] = useState("");
  const showLoading = useLoadingScreenStore((state) => state.showLoading);
  const hideLoading = useLoadingScreenStore((state) => state.hideLoading);

  const getRequestData = async () => {
    showLoading();
    try {
      const requestData = await RequestService.getRequestData(
        currentRequest.file_number
      );

      const {
        Address: { Street, Number, ABLNumber },
      } = requestData;

      const addressObj = {
        address_street: Street,
        address_number: parseInt(Number),
        abl_number: ABLNumber,
        property_id: requestData.PropertyID,
      };

      var debt: boolean | null = false;
      if (requestData.UserType === "Admin") {
        debt = null;
      }

      setAdress(addressObj);
      sessionStorage.setItem("address", JSON.stringify(addressObj));
      // TODO ASK BACK TO CHANGE IT TO BOOLEAN IF POSSIBLE

      setAblDebt(debt);
      if (debt === false) {
        sessionStorage.setItem("ablDebt", debt.toString());
      } else {
        setCommonZone(true);
        sessionStorage.setItem("commonZone", "true");
      }

      setUserType(requestData.UserType);
      sessionStorage.setItem("userType", requestData.UserType);
      const files = requestData.Documents.reduce((acc: FilesType, doc) => {
        if (doc.Type === 16 || doc.Type === 17) {
          return acc;
        }
        const file = base64ToFile(doc.Content, doc.Name);
        acc[doc.Type] = {
          content: file,
          name: doc.Name,
        };
        return acc;
      }, {});
      setFiles(files);
      setSelectedActivities(requestData.SelectedActivities);
      sessionStorage.setItem(
        "selectedActivities",
        JSON.stringify(requestData.SelectedActivities)
      );

      setProjectDescription(requestData.ProjectDesc);
      sessionStorage.setItem("projectDescription", requestData.ProjectDesc);
      setEstimatedTime(requestData.EstimatedTime);
      sessionStorage.setItem(
        "estimatedTime",
        String(requestData.EstimatedTime)
      );
      setInsurance(requestData.Insurance);
      sessionStorage.setItem("insurance", String(requestData.Insurance));
      setObservations(requestData.Observations);
      sessionStorage.setItem("observations", requestData.Observations);
      setTasksObservations(requestData.ObservationsTasks);
      sessionStorage.setItem(
        "tasksObservations",
        requestData.ObservationsTasks
      );
      if (requestData.Observations !== "") {
        setActiveStep(1);
      } else {
        setActiveStep(2);
      }
      setUpdate(true);
      setShowStepper(false);
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
    // TODO: CHECK SESSION RESET
    if (currentRequest) {
      getRequestData();
    }
  }, []);

  useEffect(() => {
    sessionStorage.setItem("currentStep", activeStep.toString());
  }, [activeStep]);

  const handleStepClick = (index: number) => {
    if (index > activeStep && disabled) {
      return;
    }
    setActiveStep(index);
  };

  const handleNext = () => {
    if (!validateStep(activeStep)) {
      alert("Complete la información requerida para poder continuar.");
    }

    if (!isLastStep) {
      setDisabled(true);
      var newStep = activeStep + 1;
      if (newStep === 2 && update && tasksObservations === "") {
        newStep++;
      }
      if (newStep === 3) {
        setIsLastStep(true);
      }
      setActiveStep(newStep);
    } else {
      handleSubmit();
    }
  };

  const handleSubmit = async () => {
    setCreatingRecord(true);
    setDisabled(true);

    const payload = {
      address,
      ablDebt,
      commonZone,
      userType,
      selectedActivities,
      projectDescription,
      estimatedTime,
      insurance,
      files,
    };

    try {
      var result: RequestResponse;
      if (currentRequest) {
        result = await RequestService.updateRequest(
          payload,
          currentRequest.file_number
        );
      } else {
        result = await RequestService.saveRequest(payload);
      }

      if (result.success) {
        sessionStorage.clear();
        navigate("/admin/requests/success", {
          state: { updated: currentRequest },
        });
      }
    } catch (error) {
      if (error instanceof ValidationError) {
        alert(`Error en los datos del formulario: ${error.message}`);
      } else if (error instanceof RequestError) {
        alert(`Error al enviar la solicitud: ${error.message}`);
      } else {
        console.error("Error inesperado:", error);
        alert("Ocurrió un error inesperado. Inténtalo de nuevo más tarde.");
      }
    } finally {
      setCreatingRecord(false);
      setDisabled(false);
    }
  };

  const validateStep = (step: number): boolean => {
    switch (step) {
      case 1:
        if (
          address === null ||
          (!commonZone && (ablDebt === null || ablDebt))
        ) {
          return false;
        }
        break;
      case 2:
        if (
          address === null ||
          (!commonZone && (ablDebt === null || ablDebt))
        ) {
          return false;
        }
        break;
      default:
        break;
    }

    return true;
  };

  const handlePrev = () => {
    if (isFirstStep && !update) {
      sessionStorage.clear();
      navigate("/admin/requests");
    }

    var newStep = activeStep - 1;
    if (newStep === 2 && update && tasksObservations === "") {
      newStep--;
    }
    if (
      (newStep === 1 && update && observations === "") ||
      (newStep === 0 && update)
    ) {
      sessionStorage.clear();
      navigate("/admin/requests");
      return;
    }
    setIsLastStep(false);
    setActiveStep(newStep);
  };

  const handleDebtCancel = () => {
    setAblDebt(null);
    setAdress(null);
    window.open("https://boletadepago.gestionmsi.gob.ar/alumbrado", "_blank");
  };

  const steps = [
    {
      component: (
        <StepOne
          setDisabled={setDisabled}
          address={address}
          onAddressChange={setAdress}
          ablDebt={ablDebt}
          setAblDebt={setAblDebt}
          commonZone={commonZone}
          setCommonZone={setCommonZone}
        />
      ),
      label: "Domicilio & ABL",
    },
    {
      component: (
        <StepTwo
          observations={observations}
          setDisabled={setDisabled}
          commonZone={commonZone}
          userType={userType}
          setUserType={setUserType}
          files={files}
          setFiles={setFiles}
        />
      ),
      label: "Documentación",
    },
    {
      component: (
        <StepThree
          files={files}
          observations={tasksObservations}
          setFiles={setFiles}
          selectedActivities={selectedActivities}
          setSelectedActivity={setSelectedActivities}
          projectDescription={projectDescription}
          setProjectDescription={setProjectDescription}
          estimatedTime={estimatedTime}
          setEstimatedTime={setEstimatedTime}
          insurance={insurance}
          setInsurance={setInsurance}
          setDisabled={setDisabled}
        />
      ),
      label: "Memoria descriptiva",
    },
    {
      component: (
        <StepResume
          address={address}
          update={update}
          files={files}
          selectedActivities={selectedActivities}
          projectDescription={projectDescription}
          estimatedTime={estimatedTime}
          insurance={insurance}
          setDisabled={setDisabled}
          setActiveStep={setActiveStep}
        />
      ),
      label: "Resumen",
    },
  ];

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1">
        <div className="grid grid-cols-1 gap-2 max-w-4xl mx-auto mb-2 mt-2">
          <h1 className="text-2xl font-semibold tracking-tight md:text-4xl">
            Aviso de obra
          </h1>
          {showStepper && (
            <div className="mb-10 w-full py-4">
              <Stepper
                activeLineClassName="bg-primary-dark"
                activeStep={activeStep}
                isLastStep={(value) => setIsLastStep(value)}
                isFirstStep={(value) => setIsFirstStep(value)}
              >
                {steps.map((step, index) => (
                  <Step
                    activeClassName="bg-white"
                    key={index}
                    onClick={() => handleStepClick(index)}
                  >
                    <span
                      className={`flex items-center justify-center w-8 h-8 rounded-full cursor-pointer ${
                        index < activeStep
                          ? "bg-green-700 text-white"
                          : index === activeStep
                          ? "border-2 border-green-700"
                          : "border-2 border-gray-400"
                      }`}
                    >
                      {index < activeStep ? (
                        <svg
                          className="w-4 h-4 text-white"
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 16 12"
                        >
                          <path
                            stroke="currentColor"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M1 5.917 5.724 10.5 15 1.5"
                          />
                        </svg>
                      ) : index > activeStep ? (
                        ""
                      ) : (
                        <span className="w-2 h-2 rounded-full bg-green-700"></span>
                      )}
                    </span>
                    <div className="absolute -bottom-[1.5rem] w-max text-center">
                      <Typography
                        color={activeStep === index ? "blue-gray" : "gray"}
                        className="font-normal text-xs"
                        children={step.label}
                      />
                    </div>
                  </Step>
                ))}
              </Stepper>
            </div>
          )}
          <div>{steps[activeStep].component}</div>
        </div>
      </div>
      <div className="w-full mt-16">
        <div className="max-w-4xl mx-auto flex justify-between">
          {ablDebt === null || ablDebt === false ? (
            <>
              <CustomButton
                onClick={handlePrev}
                className="w-full mr-1"
                color="light"
              >
                Volver
              </CustomButton>
              <CustomButton
                onClick={handleNext}
                disabled={disabled}
                className="bg-green-700 text-white w-full ml-1"
                isLoading={isLastStep && creatingRecord}
              >
                {!isLastStep ? "Continuar" : "Firmar Solicitud"}
              </CustomButton>
            </>
          ) : (
            <CustomButton
              onClick={handleDebtCancel}
              disabled={ablDebt === null || address === null ? true : false}
              className="bg-green-700 text-white w-full ml-1"
            >
              Pagar deuda
            </CustomButton>
          )}
        </div>
      </div>
    </div>
  );
};

export default FormRequest;
