import type { PropsWithChildren } from 'react'
import { BrowserRouter } from 'react-router-dom'

export function AppProvider({ children }: PropsWithChildren) {
  return <BrowserRouter>{children}</BrowserRouter>
}
