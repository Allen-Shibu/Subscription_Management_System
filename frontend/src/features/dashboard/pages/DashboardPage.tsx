import { Link } from 'react-router-dom'

export function DashboardPage() {
  return (
    <main className="page-shell">
      <header className="page-header">
        <div>
          <p className="eyebrow">Brightloop · Billing</p>
          <h1>Dashboard</h1>
        </div>
        <Link to="/auth/login">Sign out</Link>
      </header>
      <section className="panel" aria-label="Dashboard overview">
        <p>Dashboard feature placeholder</p>
      </section>
    </main>
  )
}
