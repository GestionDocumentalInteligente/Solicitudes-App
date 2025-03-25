import { ProtectedLayout } from "./layout/Main";
import { AuthProvider } from "./pages/auth/provider/useAuth";
import ErrorPage from "./pages/ErrorPage";
import { Navigate } from "react-router-dom";
import { BaseLayout } from "./layout/BaseLayout";
import CompleteRegistration from "./pages/auth/register/Register";
import Requests from "./pages/categories/requests/Requests";
import FormRequest from "./pages/categories/requests/new/FormRequest";
import RequestVerification from "./pages/requestVerification/RequestVerification";
import DataVerification from "./pages/requestVerification/DataVerification";
import EmailConfirmation from "./pages/auth/register/EmailConfirmation";
import ActivateAccount from "./pages/auth/register/ActivateAccount";
import SignInPage from "./pages/auth/login/Login";
import RequestSuccess from "./pages/categories/requests/new/RequestSuccess";
import RequestValidation from "@/pages/requestValidation/RequestValidation.tsx";
import DataValidation from "@/pages/requestValidation/DataValidation.tsx";
import { ProtectedAdminLayout } from "./layout/AdminLayout";
import RolPage from "./pages/auth/login/Rol";

export default [
  {
    path: "",
    element: <Navigate to="/admin/requests" />,
    errorElement: <ErrorPage />,
  },
  {
    path: "login",
    element: (
      <BaseLayout>
        <SignInPage />
      </BaseLayout>
    ),
  },
  {
    path: "email-verification-sent",
    element: (
      <BaseLayout>
        <EmailConfirmation />
      </BaseLayout>
    ),
  },
  {
    path: "activate-account",
    element: (
      <AuthProvider>
        <BaseLayout>
          <ActivateAccount />
        </BaseLayout>
      </AuthProvider>
    ),
  },
  {
    path: "complete-registration",
    element: (
      <BaseLayout>
        <CompleteRegistration />
      </BaseLayout>
    ),
  },
  {
    path: "/admin",
    element: <ProtectedLayout />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "requests",
        element: <Requests />,
      },
      {
        path: "requests/new",
        element: <FormRequest />,
      },
      {
        path: "requests/success",
        element: <RequestSuccess />,
      },
      {
        path: "",
        element: <Navigate to="requests" />,
        errorElement: <ErrorPage />,
      },
      {
        path: "*",
        element: <Navigate to="requests" />,
        errorElement: <ErrorPage />,
      },
    ],
  },
  {
    path: "/admin-panel",
    element: <ProtectedAdminLayout />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "",
        element: <RolPage />,
      },
      {
        path: "request-verification",
        element: <RequestVerification />,
      },
      {
        path: "data-verification",
        element: <DataVerification />,
      },
      {
        path: "request-validation",
        element: <RequestValidation />,
      },
      {
        path: "data-validation",
        element: <DataValidation />,
      },
      {
        path: "*",
        element: <Navigate to="request-verification" />,
        errorElement: <ErrorPage />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/admin/requests" />,
    errorElement: <ErrorPage />,
  },
];
