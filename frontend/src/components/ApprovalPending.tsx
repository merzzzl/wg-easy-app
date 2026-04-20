import { Box, EmptyStateContent, EmptyStateDescription, EmptyStateIndicator, EmptyStateRoot, EmptyStateTitle, VStack } from '@chakra-ui/react'
import { useTranslation } from 'react-i18next'

import type { MeResponse } from '../lib/api'

type Props = {
  me: MeResponse
}

export function ApprovalPending({ me }: Props) {
  const { t } = useTranslation()

  return (
    <Box minH="100vh" display="grid" placeItems="center" py="8">
      <EmptyStateRoot maxW="md">
        <EmptyStateContent>
          <EmptyStateIndicator>{me.user.username.slice(0, 2).toUpperCase()}</EmptyStateIndicator>
          <VStack textAlign="center">
            <EmptyStateTitle>{t('approval.title')}</EmptyStateTitle>
            <EmptyStateDescription>{t('approval.body')}</EmptyStateDescription>
          </VStack>
        </EmptyStateContent>
      </EmptyStateRoot>
    </Box>
  )
}
