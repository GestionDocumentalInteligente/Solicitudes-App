import React from "react";

interface StyledTitleProps {
  title: string;
}

const StyledTitle: React.FC<StyledTitleProps> = ({ title }) => {
  return <h1 className="text-xl font-semibold text-gray-900">{title}</h1>;
};

export default StyledTitle;
