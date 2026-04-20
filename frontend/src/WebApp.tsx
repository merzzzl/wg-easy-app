import { AbsoluteCenter, Box, EmptyStateContent, EmptyStateDescription, EmptyStateIndicator, EmptyStateRoot, EmptyStateTitle, VStack } from '@chakra-ui/react'

export default function WebApp() {
  return (
    <Box w="100vw" h="100vh">
      <AbsoluteCenter>
        <EmptyStateRoot maxW="420px">
          <EmptyStateContent>
            <EmptyStateIndicator>Telegram</EmptyStateIndicator>
            <VStack textAlign="center">
              <EmptyStateTitle>Откройте приложение внутри Telegram</EmptyStateTitle>
              <EmptyStateDescription>
                Этот интерфейс работает как Telegram Mini App. Перейдите в бот и нажмите кнопку открытия приложения.
              </EmptyStateDescription>
            </VStack>
          </EmptyStateContent>
        </EmptyStateRoot>
      </AbsoluteCenter>
    </Box>
  )
}
