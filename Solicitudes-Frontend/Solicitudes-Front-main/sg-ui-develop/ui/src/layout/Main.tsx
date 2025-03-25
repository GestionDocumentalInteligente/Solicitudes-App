import { useEffect, useState } from "react";
import { Outlet, useNavigate } from "react-router-dom";

import Navbar from "./Navbar/Navbar";
import { MVPFooter } from "./Footer/MVPFooter";
import { BaseModal } from "../components/Modal/BaseModal";
import { AuthProvider, useAuth } from "../pages/auth/provider/useAuth";
import StoreComponents from "@/layout/StoreComponents/StoreComponents.tsx";
import { AuthService } from "@/pages/auth/authService";

const MainLayout: React.FC = () => {
  const [isLogoutModalOpen, setIsLogoutModalOpen] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    auth.validateRedirection();
  }, [auth]);

  useEffect(() => {
    const validateToken = async () => {
      try {
        await auth.validateToken();
        const storedUser = localStorage.getItem("user");
        if (storedUser) {
          const storage = JSON.parse(storedUser) as Storage;
          const result = await AuthService.getUser(storage.data.cuil);
          if (result.admin) {
            setIsAdmin(true);
          }
        } else {
          navigate("/login");
        }
      } catch (error) {
        console.error(error);
        auth.logout();
      }
    };

    validateToken();
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar
        setIsLogoutModalOpen={() => setIsLogoutModalOpen(true)}
        admin={isAdmin}
      />
      <div className="flex flex-grow">
        {/* <Sidebar
          isOpen={isSidebarOpen}
          setIsLogoutModalOpen={() => setIsLogoutModalOpen(true)}
          onClose={() => setIsSidebarOpen(false)}
          role="local"
        /> */}
        <main className="flex-grow p-4 mt-14 mx-auto">
          <StoreComponents />
          <Outlet />
          <BaseModal
            msg="¿Estás seguro que deseás salir?"
            handleConfirmation={() => auth.logout()}
            isOpen={isLogoutModalOpen}
            onClose={() => setIsLogoutModalOpen(false)}
          />
        </main>
      </div>
      {<MVPFooter />}
    </div>
  );
};

export const ProtectedLayout = () => {
  return (
    <AuthProvider>
      <MainLayout />
    </AuthProvider>
  );
};
