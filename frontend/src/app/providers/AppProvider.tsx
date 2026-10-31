import type { PropsWithChildren } from 'react'
import { BrowserRouter } from 'react-router-dom'
import { ClerkProvider } from '@clerk/react'

const publishableKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY

if (!publishableKey) {
  throw new Error('Add your VITE_CLERK_PUBLISHABLE_KEY to the .env file')
}

export function AppProvider({ children }: PropsWithChildren) {
  return (
    <ClerkProvider
      publishableKey={publishableKey}
      signInUrl="/auth"
      signInFallbackRedirectUrl="/dashboard"
      signUpFallbackRedirectUrl="/dashboard"
    >
      <BrowserRouter>{children}</BrowserRouter>
    </ClerkProvider>
  )
}
