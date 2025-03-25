import React from "react";

import SGLogo from "../../assets/sglogo.svg";
import Menu from "./Menu";
import { useLocation, useNavigate } from "react-router-dom";

interface NavbarProps {
  admin: boolean;
  setIsLogoutModalOpen: () => void;
}

const Navbar: React.FC<NavbarProps> = ({ setIsLogoutModalOpen, admin }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogoClick = () => {
    const { pathname } = location;

    if (pathname.startsWith("/admin-panel")) {
      if (pathname.includes("validation")) {
        navigate("/admin-panel/request-validation");
      } else {
        navigate("/admin-panel/request-verification");
      }
    } else {
      navigate("/admin/requests");
    }
  };

  return (
    <nav className="fixed top-0 z-50 w-full bg-navbar-background border-gray-500">
      <div className="px-3 py-2 lg:px-5 lg:pl-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center justify-start rtl:justify-end">
            {/* <button
              onClick={toggleSidebar}
              type="button"
              className="inline-flex items-center p-2 text-sm rounded-lg  focus:outline-none focus:ring-2 text-gray-400 bg-gray-500 hover:bg-gray-700 focus:ring-gray-600"
            >
              <svg
                className="w-6 h-6"
                aria-hidden="true"
                fill="currentColor"
                viewBox="0 0 20 20"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  clipRule="evenodd"
                  fillRule="evenodd"
                  d="M2 4.75A.75.75 0 012.75 4h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 4.75zm0 10.5a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5a.75.75 0 01-.75-.75zM2 10a.75.75 0 01.75-.75h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 10z"
                ></path>
              </svg>
            </button> */}
            <button
              onClick={handleLogoClick}
              className="bg-transparent border-none"
            >
              <img className="h-11 me-4" src={SGLogo} alt="Logo" />
            </button>
          </div>
          <div className="flex items-center">
            <Menu setIsLogoutModalOpen={setIsLogoutModalOpen} admin={admin} />
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
