import { Box } from '@chakra-ui/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { MainScreen } from './components/MainScreen'
import { Toaster } from './components/ui/toaster'
import './i18n'

export default function App() {
  const { i18n } = useTranslation()
  const [queryClient] = useState(() => new QueryClient())

  useEffect(() => {
    const language = navigator.language.toLowerCase().startsWith('ru') ? 'ru' : 'en'
    void i18n.changeLanguage(language)
  }, [i18n])

  return (
    <Box w="100vw" minH="100vh">
      <QueryClientProvider client={queryClient}>
        <Toaster />
        <MainScreen />
      </QueryClientProvider>
    </Box>
  )
}
