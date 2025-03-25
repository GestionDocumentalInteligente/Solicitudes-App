export const getStatusOptions = (status: string) => {
  switch (status.toLowerCase()) {
    case "inprogress":
      return {
        label: "En proceso",
        styles: "bg-[#0EA5E9]",
      };
    case "pending":
    case "pendiente":
    case "processing":
      return {
        label: "Pendiente",
        styles: "bg-yellow-400",
      };
    case "pendingtask":
      return {
        label: "Tarea pendiente",
        styles: "bg-yellow-400",
      };
    case "checking":
      return {
        label: "Verificando",
        styles: "bg-yellow-400",
      };
    case "finished":
      return {
        label: "Finalizado",
        styles: "bg-[#3AB266]",
      };
    case "verified":
    case "verificado":
      return {
        label: "Verificado",
        styles: "bg-[#526D4E]",
      };
    case "validated":
      return {
        label: "Validado",
        styles: "bg-[#526D4E]",
      };
    case "observed":
    case "observado":
      return {
        label: "Observado",
        styles: "bg-red-500",
      };
    case "rejected":
      return {
        label: "Rechazado",
        styles: "bg-red-500",
      };
    case "error":
      return {
        label: "Error",
        styles: "bg-red-500",
      };
    default:
      return {
        label: "",
        styles: "",
      };
  }
};
