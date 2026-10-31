import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { ClerkLoaded, Show } from "@clerk/react";
import { AuthPage } from "../../features/auth/pages/AuthPage";
import { DashboardPage } from "../../features/dashboard/pages/DashboardPage";
import { BillingsPage } from "../../features/billings/pages/BillingsPage";
import { StaffPage } from "../../features/staff/pages/StaffPage";
import { ServerStaffPage } from "../../features/staff/pages/ServerStaffPage";
import { CatalogPage } from "../../features/admin/pages/CatalogPage";
import { ServerTaxesPage } from "../../features/admin/pages/ServerTaxesPage";
import { CatalogList } from "../../features/admin/catalog/CatalogList";
import { ProductForm } from "../../features/admin/catalog/ProductForm";
import { PlanForm } from "../../features/admin/catalog/PlanForm";
import { usingMock } from "../../features/admin/catalog/api";
import { ServerProductsPage } from "../../features/admin/catalog/ServerProductsPage";
import { ServerProductForm } from "../../features/admin/catalog/ServerProductForm";
import { SubscriptionsPage } from "../../features/subscriptions/pages/SubscriptionsPage";
import { ServerSubscriptionsPage } from "../../features/subscriptions/pages/ServerSubscriptionsPage";
import { NotFoundPage } from "../../pages/notFoundPage/NotFoundPage";
import { MainLayout } from "../../shared/components/mainLayout/MainLayout";
import { useIsAdmin } from "../../shared/auth";

function AdminOnly() {
  const isAdmin = useIsAdmin();
  return isAdmin ? <Outlet /> : <div className="panel"><h1>Not allowed</h1><p>Only an admin can open this page.</p></div>;
}

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
        <Route path="/subscriptions" element={usingMock ? <SubscriptionsPage /> : <ServerSubscriptionsPage />} />
        <Route path="/billing" element={<BillingsPage />} />
        <Route path="/products" element={usingMock ? <CatalogList kind="products" readonly /> : <ServerProductsPage />} />
        {!usingMock && <Route path="/products/:id" element={<ServerProductForm readonly />} />}
        <Route element={<AdminOnly />}>
          <Route path="/admin/products" element={usingMock ? <CatalogList kind="products" /> : <ServerProductsPage admin />} />
          <Route path="/admin/products/new" element={usingMock ? <ProductForm /> : <ServerProductForm />} />
          <Route path="/admin/products/:id/edit" element={usingMock ? <ProductForm /> : <ServerProductForm />} />
          <Route path="/admin/plans" element={<CatalogList kind="plans" />} />
          <Route path="/admin/plans/new" element={<PlanForm />} />
          <Route path="/admin/plans/:id/edit" element={<PlanForm />} />
          <Route path="/plans" element={<Navigate to="/admin/plans" replace />} />
          <Route path="/taxes" element={usingMock ? <CatalogPage section="taxes" /> : <ServerTaxesPage />} />
          <Route path="/discounts" element={<CatalogPage section="discounts" />} />
          <Route path="/staff" element={usingMock ? <StaffPage /> : <ServerStaffPage />} />
          <Route path="/quotation-templates" element={<CatalogPage section="quotation-templates" />} />
        </Route>
      </Route>
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
