import { Navigate, Route, Routes } from "react-router-dom";
import { ClerkLoaded, Show } from "@clerk/react";
import { AuthPage } from "../../features/auth/pages/AuthPage";
import { DashboardPage } from "../../features/dashboard/pages/DashboardPage";
import { NotFoundPage } from "../../pages/notFoundPage/NotFoundPage";
import { MainLayout } from "../../shared/components/mainLayout/MainLayout";

export function RouteManager() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="/auth/*" element={<AuthPage />} />

      <Route
        element={
          <ClerkLoaded>
            <Show when="signed-in" fallback={<Navigate to="/auth" replace />}>
              <MainLayout />
            </Show>
          </ClerkLoaded>
        }
      >
        <Route path="/dashboard" element={<DashboardPage />} />
      </Route>
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
