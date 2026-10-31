import { SignIn } from '@clerk/react'

export function AuthPage() {
  return (
    <main className="auth-page">
      <SignIn routing="hash" appearance={{ variables: { colorPrimary: 'var(--color-primary)' } }} />
    </main>
  )
}
