import { AppProvider } from './providers/AppProvider'
import { RouteManager } from './routes/RouteManager'

export default function App() {
  return (
    <AppProvider>
      <RouteManager />
    </AppProvider>
  )
}
