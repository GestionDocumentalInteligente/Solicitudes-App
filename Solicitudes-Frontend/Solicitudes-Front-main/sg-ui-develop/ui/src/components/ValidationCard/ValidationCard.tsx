import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { ChevronUp, CircleCheck, CircleX, SquarePen, X } from "lucide-react";
import { cn } from "@/lib/utils";
// TODO: MOVE TYPES TO THIS COMPONENT
import {
  DocumentContent,
  Section,
  VerificationStatus,
} from "@/pages/requestVerification/types.ts";
import { Checkbox } from "@/components/ui/checkbox.tsx";

interface ValidationCardProps {
  section: Section;
  onUpdate: (updates: Partial<Section>, updateNext?: boolean) => void;
  documentContent?: DocumentContent;
  isValidation?: boolean;
  onQuestionClick?: (index: number) => void;
  deleteSteps?: (index: number) => void;
}

export default function ValidationCard({
  section,
  onUpdate,
  documentContent,
  isValidation = false,
  onQuestionClick,
  deleteSteps
}: ValidationCardProps) {
  const [isOpen, setIsOpen] = useState(false);

  const handleInputChange = (value: string) => {
    onUpdate({ value });
  };

  const handleOptionsChange = (optionId: string) => {
    onUpdate({
      selectedObservationOptions: section.selectedObservationOptions?.includes(
        optionId
      )
        ? section.selectedObservationOptions.filter((id) => id !== optionId)
        : [...(section.selectedObservationOptions || []), optionId],
    });
  };

  const handleObservationChange = (observation: string) => {
    onUpdate({ observation });
  };

  const handleConfirm = () => {
    onUpdate({ isConfirmed: true, isExpanded: false }, true);
  };

  const handleEdit = () => {
    if (isValidation) {
      onUpdate({
        isConfirmed: true,
        status: section.status === "valid" ? "invalid" : "valid",
        selectedObservationOptions: [""],
      });
    } else {
      onUpdate({
        isConfirmed: false,
        status: section.status === "valid" ? null : section.status,
        isExpanded: section.status !== "valid",
      });
    }
  };

  const handleValidation = (status: VerificationStatus) => {
    onUpdate(
      {
        status,
        isConfirmed: false,
        observation: undefined,
        isExpanded: status === "invalid",
      },
      true
    );
  };

  const handleExpand = () => {
    onUpdate({ isExpanded: !section.isExpanded });
  };

  const handleQuestionClick = (e: React.MouseEvent<HTMLDivElement>) => {
    const target = e.target as HTMLElement;
    if (target.matches("a[data-doc]")) {
      e.preventDefault();
      if (onQuestionClick) {
        onQuestionClick(parseInt(target.getAttribute("data-doc") || "0"));
      }
    }
  };

  const selectedObservationsDescription = (observedSection: Section) => {
    return observedSection.observationOptions
      ?.reduce((acc: string[], activity) => {
        if (
            observedSection.selectedObservationOptions?.includes(
                activity.id
            )
        ) {
          acc.push(activity.description);
        }
        return acc;
      }, [])
      .join(", ")
  }

  const renderValidationButtons = (section: Section) => (
    <div className="flex gap-2">
      <Button
        variant="outline"
        onClick={() => handleValidation("invalid")}
        className="flex-1"
        disabled={!!section.inputType && !section.value}
      >
        Inválido
      </Button>
      <Button
        variant="outline"
        onClick={() => handleValidation("valid")}
        className="flex-1 bg-[#3b5c3f] hover:bg-[#2e4831] hover:text-white text-white"
        disabled={!!section.inputType && !section.value}
      >
        Válido
      </Button>
    </div>
  );

  const renderInput = (section: Section) => {
    switch (section.inputType) {
      case "radio":
        return (
          // TODO: IMPROVE CHECKBOX SPACING AND DOCUMENT CONTENT NULLITY
          <RadioGroup
            value={section.value}
            onValueChange={(value) => {
              handleInputChange(value);
              handleValidation(
                value === documentContent?.requiresInsurance
                  ? "valid"
                  : value === "No" ? "valid" : "invalid"
              );
              if (value !== documentContent?.requiresInsurance && value === "No" && deleteSteps) {
                deleteSteps(1);
              }
            }}
            className="flex gap-24"
          >
            {section.inputOptions?.map((option) => (
              <div key={option} className="flex items-center space-x-2">
                <RadioGroupItem value={option} id={`${section.id}-${option}`} />
                <Label htmlFor={`${section.id}-${option}`}>{option}</Label>
              </div>
            ))}
          </RadioGroup>
        );
      case "text":
        return (
          <>
            <Input
              className="placeholder:text-wrap placeholder:top-[-7px] placeholder:relative h-12"
              value={section.value}
              onChange={(e) => handleInputChange(e.target.value)}
              placeholder={section.inputPlaceholder}
            />
            {renderValidationButtons(section)}
          </>
        );
      default:
        return renderValidationButtons(section);
    }
  };

  const StatusBadge = ({ status }: { status: VerificationStatus }) => {
    if (!status) return null;

    return (
      <span
        className={cn(
          "px-2 py-0.5 text-xs font-medium rounded-full flex items-center",
          status === "valid" && "bg-[#526D4E] text-white",
          status === "invalid" && "bg-red-500 text-white"
        )}
      >
        {status === "valid" ? (
          <>
            <CircleCheck className="h-3 w-3 mr-1" /> Válido
          </>
        ) : (
          <>
            <CircleX className="h-3 w-3 mr-1" /> Inválido
          </>
        )}
      </span>
    );
  };

  return (
    <>
      <div
        className={cn(
          "border rounded-lg transition-colors",
          section.status === "valid" &&
            "bg-[#3b5c3f]/10 border-solid border-[#3b5c3f]",
          section.status === "invalid" &&
            "bg-destructive/5 border-solid border-destructive"
        )}
      >
        {!section.status ? (
          <div className="p-4 space-y-4">
            <div className="text-sm flex items-center font-bold justify-between">
              <span>{section.title}</span>
            </div>
            <p className="text-sm font-medium text-gray-900">
              {section.question}
            </p>
            {renderInput(section)}
          </div>
        ) : (
          <Accordion
            type="single"
            collapsible
            value={section.isExpanded ? section.id : ""}
            onValueChange={() => handleExpand()}
          >
            <AccordionItem value={section.id} className="border-0">
              <AccordionTrigger className="px-4 hover:no-underline">
                <div className="text-sm flex items-center justify-between w-full">
                  <span>{section.title}</span>
                  <StatusBadge status={section.status} />
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="px-4">
                  {section.status === "valid" ? (
                    <div className="space-y-2">
                      <p
                        className="text-sm text-gray-600"
                        // TODO: CHECK THIS
                        dangerouslySetInnerHTML={{
                          __html:
                            section.value ||
                            selectedObservationsDescription(section) ||
                            section.question ||
                            "",
                        }}
                        onClick={handleQuestionClick}
                      ></p>
                      <Button
                        variant="outline"
                        className="w-full mt-4"
                        onClick={() => handleEdit()}
                      >
                        {/*TODO: IMPROVE COUPLING*/}
                        {isValidation ? (
                          <>
                            <X className="" />
                            Rechazar
                          </>
                        ) : (
                          <>
                            <SquarePen className="h-4 w-4 mr-2" /> Editar{" "}
                          </>
                        )}
                      </Button>
                    </div>
                  ) : section.status === "invalid" && section.isConfirmed ? (
                    <div className="space-y-2">
                      <p
                        className="text-sm text-gray-600"
                        // TODO: CHECK THIS
                        dangerouslySetInnerHTML={{
                          __html:
                            selectedObservationsDescription(section) ||
                            section.value ||
                            section.question ||
                            "",
                        }}
                      ></p>
                      {!section.selectedObservationOptions && (
                        <>
                          <p className="text-sm font-medium text-gray-900">
                            Observaciones:
                          </p>
                          <p className="text-sm text-gray-600">
                            {section.observation}
                          </p>
                        </>
                      )}
                      <Button
                        variant="outline"
                        className="w-full"
                        onClick={() => handleEdit()}
                      >
                        <SquarePen className="h-4 w-4 mr-2" />
                        Editar
                      </Button>
                    </div>
                  ) : section.observationOptions ? (
                    <div className="space-y-4">
                      <p className="text-sm text-gray-600">
                        {section.value || section.question}
                      </p>
                      <div className="space-y-2">
                        <div className="relative">
                          <Button
                            type="button"
                            onClick={() => setIsOpen(!isOpen)}
                            className={cn(
                              "w-full justify-between border-2 border-[#3b5c3f] bg-white text-gray-500 hover:bg-white hover:text-gray-600",
                              "rounded-lg px-4 py-6 text-base font-normal overflow-x-hidden overflow-y-hidden"
                            )}
                          >
                            {/*TODO: CHECK CHEVRON BEHAVIOUR*/}
                            {section.selectedObservationOptions?.length
                              ? section.observationOptions
                                  .reduce((acc: string[], activity) => {
                                    if (
                                      section.selectedObservationOptions?.includes(
                                        activity.id
                                      )
                                    ) {
                                      acc.push(activity.description);
                                    }
                                    return acc;
                                  }, [])
                                  .join(", ")
                              : "Seleccionar una o más opciones"}
                            <ChevronUp
                              className={cn(
                                "h-4 w-4 shrink-0 transition-transform",
                                !isOpen && "rotate-180"
                              )}
                            />
                          </Button>

                          {isOpen && (
                            <div className="z-50 w-full space-y-4 my-2">
                              <div className="max-h-[300px] overflow-auto rounded-xl border-2 bg-white border-gray-300">
                                {section.observationOptions?.map((option) => (
                                  <div
                                    key={option.id}
                                    className={cn(
                                      "flex items-center space-x-3 px-4 py-4",
                                      section.selectedObservationOptions?.includes(
                                        option.id
                                      ) && "bg-[#3E59421A]"
                                    )}
                                  >
                                    <Checkbox
                                      id={option.id}
                                      checked={section.selectedObservationOptions?.includes(
                                        option.id
                                      )}
                                      onCheckedChange={() =>
                                        handleOptionsChange(option.id)
                                      }
                                      className="border-[#3b5c3f] data-[state=checked]:bg-white data-[state=checked]:text-[#3E5942]"
                                    />
                                    <Label
                                      htmlFor={option.id}
                                      className="flex-grow cursor-pointer text-gray-500"
                                    >
                                      {option.description}
                                    </Label>
                                  </div>
                                ))}
                              </div>
                            </div>
                          )}
                          <div className="flex gap-2 pt-2">
                            <Button
                              variant="outline"
                              onClick={() => {
                                setIsOpen(false);
                                handleValidation(null);
                              }}
                              className="flex-1"
                            >
                              Cancelar
                            </Button>
                            <Button
                              onClick={() => {
                                setIsOpen(false);
                                handleObservationChange(
                                    `Tarea asignada: ${ selectedObservationsDescription(section) ?? "" }`
                                );
                                handleConfirm();
                              }}
                              className="flex-1 bg-[#3b5c3f] hover:bg-[#2e4831]"
                              disabled={
                                !section.selectedObservationOptions?.length
                              }
                            >
                              Cambiar
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      <p className="text-sm text-gray-600">
                        {section.value || section.question}
                      </p>
                      <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-900">
                          Observaciones
                        </label>
                        <Textarea
                          placeholder={section.observationPlaceholder}
                          value={section.observation || ""}
                          onChange={(e) =>
                            handleObservationChange(e.target.value)
                          }
                          className="min-h-[100px]"
                        />
                      </div>
                      <div className="flex gap-2">
                        <Button
                          variant="outline"
                          className="flex-1"
                          onClick={() => handleValidation(null)}
                        >
                          Cancelar
                        </Button>
                        <Button
                          className="flex-1 bg-[#3b5c3f] hover:bg-[#2e4831]"
                          disabled={!section.observation}
                          onClick={() => handleConfirm()}
                        >
                          Confirmar
                        </Button>
                      </div>
                    </div>
                  )}
                </div>
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        )}
      </div>
      {(section.isExpanded || !section.status) && (
        <span className="text-[10px] leading-snug block text-gray-400">
          {section.questionClarification}
        </span>
      )}
    </>
  );
}
