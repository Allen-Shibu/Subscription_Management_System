import { createRoot } from 'react-dom/client'
import { ClerkProvider } from '@clerk/react'
import App from './App.jsx'
import './style.css'

const publishableKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY

if (!publishableKey) {
  throw new Error('Add your VITE_CLERK_PUBLISHABLE_KEY to the .env file')
}

createRoot(document.getElementById('app')).render(
  <ClerkProvider publishableKey={publishableKey}>
    <App />
  </ClerkProvider>,
)
