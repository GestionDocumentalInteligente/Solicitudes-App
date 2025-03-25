import { Outlet } from "react-router-dom";
import CityView from "../assets/bg_san_isidro-min.png";
import SGLogo from "../assets/sglogo.svg";

interface BaseLayoutProps {
  children?: React.ReactNode;
}

export const BaseLayout: React.FC<BaseLayoutProps> = ({ children }) => {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 h-screen overflow-hidden">
      <div className="flex flex-col h-full">
        <div className="sm:p-4">
          <div className="flex">
            <img src={SGLogo} width={270} alt="Logo" />
          </div>
          <hr />
          <main className="flex-grow px-12 pt-2">{children || <Outlet />}</main>
        </div>
      </div>
      <div className="hidden lg:block">
        <img
          src={CityView}
          alt="Catedral de San Isidro"
          className="object-cover h-screen w-full"
        />
      </div>
    </div>
  );
};
