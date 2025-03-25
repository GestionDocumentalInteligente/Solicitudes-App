import { useEffect, useState } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import Navbar from "./Navbar/Navbar";
// import Footer from "./Footer/Footer";
import LoadingScreen from "@/components/LoadingScreen/LoadingScreen.tsx";
import { AuthProvider, useAuth } from "@/pages/auth/provider/useAuth";
import { BaseModal } from "@/components/Modal/BaseModal";
import { AuthService } from "@/pages/auth/authService";
import { User } from "@/pages/categories/requests/types";
import StoreComponents from "@/layout/StoreComponents/StoreComponents.tsx";

type Storage = {
  data: User;
};

const AdminLayout: React.FC = () => {
  const [isLogoutModalOpen, setIsLogoutModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const validateAdminAccess = async () => {
      try {
        await auth.validateToken();
        const storedUser = localStorage.getItem("user");
        if (storedUser) {
          const storage = JSON.parse(storedUser) as Storage;
          const result = await AuthService.getUser(storage.data.cuil);
          if (!result.admin) {
            navigate("/admin/requests");
          }
        } else {
          navigate("/login");
        }
      } catch (error) {
        console.error(error);
        auth.logout();
      } finally {
        setIsLoading(false);
      }
    };

    validateAdminAccess();
  }, [auth]);

  if (isLoading) {
    return <LoadingScreen title={["Validando acceso..."]} description={[""]} />;
  }

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar
        setIsLogoutModalOpen={() => setIsLogoutModalOpen(true)}
        admin={true}
      />
      <div className="flex flex-grow">
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
      {/*<Footer/>*/}
    </div>
  );
};

export const ProtectedAdminLayout = () => {
  return (
    <AuthProvider>
      <AdminLayout />
    </AuthProvider>
  );
};
