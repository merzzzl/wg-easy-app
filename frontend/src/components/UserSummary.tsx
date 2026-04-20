import { Avatar, Badge, Flex, Heading, HStack, Text, VStack } from '@chakra-ui/react'
import { useTranslation } from 'react-i18next'

import type { MeResponse } from '../lib/api'

type Props = {
  me?: MeResponse
  remaining: number
}

export function UserSummary({ me, remaining }: Props) {
  const { t } = useTranslation()

  return (
    <Flex gap="4" p="5" rounded="3xl" bg="bg.panel" borderWidth="1px" align="center">
      <Avatar.Root size="2xl">
        <Avatar.Fallback name={me?.user.username ?? 'User'} />
      </Avatar.Root>

      <VStack align="start" gap="1" flex="1">
        <Heading size="lg">{me ? `@${me.user.username}` : '...'}</Heading>
        <HStack gap="2" wrap="wrap">
          <Badge colorPalette={me?.user.status === 'approved' ? 'green' : 'orange'}>{me?.user.status ?? 'pending'}</Badge>
          <Badge variant="subtle">{t('profile.limit', { used: me?.used_tunnels ?? 0, max: me?.max_tunnels ?? 0 })}</Badge>
        </HStack>
        <Text color="fg.muted">{t('profile.slots', { count: remaining })}</Text>
      </VStack>
    </Flex>
  )
}
