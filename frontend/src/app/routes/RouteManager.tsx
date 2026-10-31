import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { ClerkLoaded, Show } from "@clerk/react";
import { AuthPage } from "../../features/auth/pages/AuthPage";
import { DashboardPage } from "../../features/dashboard/pages/DashboardPage";
import { BillingsPage } from "../../features/billings/pages/BillingsPage";
import { StaffPage } from "../../features/staff/pages/StaffPage";
import { CatalogPage } from "../../features/admin/pages/CatalogPage";
import { CatalogList } from "../../features/admin/catalog/CatalogList";
import { ProductForm } from "../../features/admin/catalog/ProductForm";
import { PlanForm } from "../../features/admin/catalog/PlanForm";
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
        <Route element={<Show when={{ role: "org:admin" }} fallback={<div className="panel"><h1>Not allowed</h1><p>Only an organization admin can open this page.</p></div>}><Outlet /></Show>}>
          <Route path="/admin/products" element={<CatalogList kind="products" />} />
          <Route path="/admin/products/new" element={<ProductForm />} />
          <Route path="/admin/products/:id/edit" element={<ProductForm />} />
          <Route path="/admin/plans" element={<CatalogList kind="plans" />} />
          <Route path="/admin/plans/new" element={<PlanForm />} />
          <Route path="/admin/plans/:id/edit" element={<PlanForm />} />
          <Route path="/products" element={<Navigate to="/admin/products" replace />} />
          <Route path="/plans" element={<Navigate to="/admin/plans" replace />} />
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
