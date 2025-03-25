"use client";

import { Button, Modal } from "flowbite-react";
import { HiOutlineExclamationCircle } from "react-icons/hi";

interface BaseModalProps {
  msg: string;
  isOpen: boolean;
  onClose: () => void;
  handleConfirmation: () => void;
}

export function BaseModal({
  msg,
  isOpen,
  onClose,
  handleConfirmation,
}: BaseModalProps) {
  const handleLogout = () => {
    try {
      handleConfirmation();
      onClose();
    } catch (error) {
      console.error("Error al cerrar sesión:", error);
    }
  };

  return (
    <Modal show={isOpen} size="md" onClose={onClose} popup>
      <Modal.Header />
      <Modal.Body>
        <div className="text-center">
          <HiOutlineExclamationCircle className="mx-auto mb-4 h-14 w-14 text-gray-400 dark:text-gray-200" />
          <h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
            {msg}
          </h3>
          <div className="flex justify-center gap-4">
            <Button color="gray" onClick={onClose}>
              Cancelar
            </Button>
            <Button color="failure" onClick={handleLogout}>
              Cerrar sesión
            </Button>
          </div>
        </div>
      </Modal.Body>
    </Modal>
  );
}
