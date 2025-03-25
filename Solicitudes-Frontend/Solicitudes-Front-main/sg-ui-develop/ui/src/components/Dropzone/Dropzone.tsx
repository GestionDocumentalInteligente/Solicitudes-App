import { FileInput, Label } from "flowbite-react";
import { useState } from "react";

interface DropzoneProps {
  id: string;
  onDrop: (file: File, fileName: string) => void;
}

export const Dropzone = ({ onDrop, id }: DropzoneProps) => {
  const [error, setError] = useState<string | null>(null);
  const [inputKey, setInputKey] = useState<number>(0);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];

    if (file) {
      if (file.type !== "application/pdf") {
        setError("Solo se permiten archivos en formato PDF.");
        return;
      }

      const maxSizeInBytes = 25 * 1024 * 1024; // 25 MB
      if (file.size > maxSizeInBytes) {
        setError("El archivo excede el tamaño máximo permitido (25 MB).");
        return;
      }

      setError(null);
      onDrop(file, file.name);
      setInputKey((prevKey) => prevKey + 1);
    }
  };

  return (
    <div className="flex w-full items-center justify-center">
      <Label
        htmlFor={`dropzone-${id}`}
        className="relative flex h-30 w-full cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 hover:bg-gray-100"
      >
        <div className="flex flex-col items-center justify-center pb-6 pt-5">
          <svg
            className="mb-4 h-8 w-8 text-gray-500"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 20 16"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"
            />
          </svg>
          <p className="mb-2 text-sm text-gray-500">
            Podés cargar el archivo directamente acá o arrastrarlo y soltarlo
          </p>
          <p className="text-xs text-gray-500">
            Solo admite formato PDF (MAX. 25MB)
          </p>
        </div>
        <FileInput
          key={inputKey}
          id={`dropzone-${id}`}
          className="absolute top-0 left-0 w-full h-full opacity-0 cursor-pointer"
          style={{ height: "100%" }}
          onChange={handleFileChange}
        />
      </Label>
      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
    </div>
  );
};
