import { AbsoluteCenter, Box, EmptyStateContent, EmptyStateDescription, EmptyStateIndicator, EmptyStateRoot, EmptyStateTitle, VStack } from '@chakra-ui/react'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import './i18n'

export default function WebApp() {
  const { i18n, t } = useTranslation()

  useEffect(() => {
    const language = navigator.language.toLowerCase().startsWith('ru') ? 'ru' : 'en'
    void i18n.changeLanguage(language)
  }, [i18n])

  return (
    <Box w="100vw" h="100vh">
      <AbsoluteCenter>
        <EmptyStateRoot maxW="420px">
          <EmptyStateContent>
            <EmptyStateIndicator>Telegram</EmptyStateIndicator>
            <VStack textAlign="center">
              <EmptyStateTitle>{t('fallback.title')}</EmptyStateTitle>
              <EmptyStateDescription>{t('fallback.body')}</EmptyStateDescription>
            </VStack>
          </EmptyStateContent>
        </EmptyStateRoot>
      </AbsoluteCenter>
    </Box>
  )
}
