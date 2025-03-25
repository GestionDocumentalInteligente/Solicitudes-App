import { Button, Spinner } from "flowbite-react";
import React from "react";

interface CustomButtonProps {
  isLoading?: boolean;
  type?: "button" | "submit" | "reset" | undefined;
  disabled?: boolean;
  className?: string;
  color?: string;
  onClick?: () => void;
  children: React.ReactNode;
}

const CustomButton: React.FC<CustomButtonProps> = ({
  isLoading = false,
  type = "button",
  disabled = false,
  className = "",
  color = "dark",
  onClick,
  children,
}) => {
  const baseStyles = "hover:bg-primary-dark";

  return (
    <Button
      type={type}
      className={`${baseStyles} ${className}`}
      onClick={onClick}
      disabled={disabled}
      color={color}
    >
      {isLoading ? (
        <>
          <Spinner aria-label="Loading spinner" size="sm" />
          <span className="pl-3">Loading...</span>
        </>
      ) : (
        children
      )}
    </Button>
  );
};

export default CustomButton;
