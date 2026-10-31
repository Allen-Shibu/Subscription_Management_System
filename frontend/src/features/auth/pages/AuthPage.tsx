import { SignIn } from '@clerk/react'

export function AuthPage() {
  return (
    <main className="page-shell">
      <SignIn routing="hash" />
    </main>
  )
}
