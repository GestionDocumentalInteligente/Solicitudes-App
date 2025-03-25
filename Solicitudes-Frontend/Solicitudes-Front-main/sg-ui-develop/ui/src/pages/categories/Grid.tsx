import React, { useEffect, useState } from "react";
import { BiBuildings } from "react-icons/bi";
import { FaCarAlt } from "react-icons/fa";

import StyledTitle from "../../components/Title/StyledTitle";
import Loader from "../../components/Loader/Loader";
import { getCategoriesInfo } from "./utils";
import { Category } from "./types";
import { Card } from "../../components/Card/Card";

const iconMap: { [key: number]: JSX.Element } = {
  1: <BiBuildings size={80} />,
  2: <FaCarAlt size={80} />,
};

const CategoriesGrid: React.FC = () => {
  const [processing, setProcessing] = useState(true);
  const [categories, setCategories] = useState<Category[] | null>(null);

  const getData = async () => {
    setProcessing(true);
    try {
      const result = await getCategoriesInfo();
      setCategories(result);
    } catch (e) {
      console.error(e);
      setCategories(null);
    }
    setProcessing(false);
  };

  useEffect(() => {
    getData();
  }, []);

  if (processing) {
    return <Loader />;
  }

  return (
    <>
      <StyledTitle title="Solicitudes" />
      <br />
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-5xl mx-auto">
        {categories?.map((category) => (
          <Card
            key={category.id}
            info={category}
            icon={iconMap[category.id]}
            to={`/admin/categories/${category.id}/requests`}
          />
        ))}
      </div>
    </>
  );
};

export default CategoriesGrid;
