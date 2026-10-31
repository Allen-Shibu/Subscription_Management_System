import { Navigate, Route, Routes } from 'react-router-dom'
import { AuthPage } from '../../features/auth/pages/AuthPage'
import { DashboardPage } from '../../features/dashboard/pages/DashboardPage'
import { NotFoundPage } from '../../pages/notFoundPage/NotFoundPage'

export function RouteManager() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="/auth/*" element={<AuthPage />} />
      <Route path="/dashboard" element={<DashboardPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}
