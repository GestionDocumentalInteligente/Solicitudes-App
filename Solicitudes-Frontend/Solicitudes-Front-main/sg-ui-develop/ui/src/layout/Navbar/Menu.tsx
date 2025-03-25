import { Token, User } from "@/pages/auth/types";
import { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

interface NavbarProps {
  admin: boolean;
  setIsLogoutModalOpen: () => void;
}

const Menu: React.FC<NavbarProps> = ({ setIsLogoutModalOpen, admin }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const [user, setUser] = useState<User | null>(null);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);

  const toggleDropdown = () => {
    setIsDropdownOpen(!isDropdownOpen);
  };

  const handleClickOutside = (event: MouseEvent) => {
    if (
      dropdownRef.current &&
      !dropdownRef.current.contains(event.target as Node) &&
      buttonRef.current &&
      !buttonRef.current.contains(event.target as Node)
    ) {
      setIsDropdownOpen(false);
    }
  };

  const handleClickRequests = () => {
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
    setIsDropdownOpen(false);
  };

  useEffect(() => {
    if (isDropdownOpen) {
      if (user === null) {
        const savedUser = localStorage.getItem("user");
        const token = savedUser ? (JSON.parse(savedUser) as Token) : null;
        setUser(token ? token.data : null);
      }
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isDropdownOpen]);

  return (
    <div className="flex items-center ms-3">
      <div>
        <button
          ref={buttonRef}
          type="button"
          className="flex text-sm bg-gray-800 rounded-full focus:ring-4 focus:ring-gray-600"
          aria-expanded={isDropdownOpen ? "true" : "false"}
          onClick={toggleDropdown}
        >
          <div className="w-8 h-8 rounded-full bg-gray-400 flex items-center justify-center">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fillRule="evenodd"
                d="M12 12c2.21 0 4-1.79 4-4S14.21 4 12 4 8 5.79 8 8s1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
                clipRule="evenodd"
              />
            </svg>
          </div>
        </button>
      </div>
      <div
        ref={dropdownRef}
        className={`absolute right-1 top-full mt-1 z-100 ${
          isDropdownOpen ? "" : "hidden"
        } my-4 text-base list-none bg-white divide-y divide-gray-100 rounded shadow`}
      >
        <div className="px-4 py-3" role="none">
          <p className="text-sm text-gray-900 dark:text-white" role="none">
            {user?.first_name} {user?.last_name}
          </p>
          <p
            className="text-sm font-medium text-gray-900 truncate dark:text-gray-300"
            role="none"
          >
            {user?.email}
          </p>
        </div>
        <ul className="py-1" role="none">
          {admin && (
            <li>
              <a
                href="/admin-panel"
                className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                role="menuitem"
              >
                Seleccionar perfil
              </a>
            </li>
          )}
          <li>
            <a
              href="#"
              className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
              role="menuitem"
              onClick={handleClickRequests}
            >
              Ver listado de solicitudes
            </a>
          </li>
          <li>
            <a
              href="#"
              className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
              role="menuitem"
              onClick={setIsLogoutModalOpen}
            >
              Salir
            </a>
          </li>
        </ul>
      </div>
    </div>
  );
};

export default Menu;
