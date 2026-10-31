import { Link } from 'react-router-dom'

export function AuthPage() {
  return (
    <main className="page-shell">
      <section className="panel">
        <p className="eyebrow">Recur account</p>
        <h1>Sign in</h1>
        <p>Authentication form goes here.</p>
        <Link to="/dashboard">Continue to dashboard</Link>
      </section>
    </main>
  )
}
