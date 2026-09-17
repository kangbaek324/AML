import { createBrowserRouter } from "react-router-dom";
import { Layout } from "./components/layout";
import { DashboardPage, AlertsPage } from "./pages";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: "alerts", element: <AlertsPage /> },
    ],
  },
]);
