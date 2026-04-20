import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'

import { MainScreen } from './components/MainScreen'

export default function App() {
  const [queryClient] = useState(() => new QueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      <MainScreen />
    </QueryClientProvider>
  )
}
