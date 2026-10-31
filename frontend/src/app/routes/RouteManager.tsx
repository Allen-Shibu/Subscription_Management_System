import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { ClerkLoaded, Show } from "@clerk/react";
import { AuthPage } from "../../features/auth/pages/AuthPage";
import { DashboardPage } from "../../features/dashboard/pages/DashboardPage";
import { BillingsPage } from "../../features/billings/pages/BillingsPage";
import { StaffPage } from "../../features/staff/pages/StaffPage";
import { CatalogPage } from "../../features/admin/pages/CatalogPage";
import { SubscriptionsPage } from "../../features/subscriptions/pages/SubscriptionsPage";
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
        <Route path="/subscriptions" element={<SubscriptionsPage />} />
        <Route path="/billing" element={<BillingsPage />} />
        <Route element={<Show when={{ role: "org:admin" }} fallback={<div className="panel"><h1>Admin access required</h1><p>Switch to an organization where you are an admin to open this page.</p></div>}><Outlet /></Show>}>
          <Route path="/products" element={<CatalogPage section="products" />} />
          <Route path="/plans" element={<CatalogPage section="plans" />} />
          <Route path="/taxes" element={<CatalogPage section="taxes" />} />
          <Route path="/discounts" element={<CatalogPage section="discounts" />} />
          <Route path="/staff" element={<StaffPage />} />
          <Route path="/quotation-templates" element={<CatalogPage section="quotation-templates" />} />
        </Route>
      </Route>
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
