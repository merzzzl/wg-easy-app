import { Avatar, Badge, Flex, Heading, Progress, VStack } from '@chakra-ui/react'
import { useLaunchParams, type User as TelegramUser } from '@tma.js/sdk-react'
import { useTranslation } from 'react-i18next'

import type { MeResponse } from '../lib/api'

type Props = {
  me?: MeResponse
}

export function UserSummary({ me }: Props) {
  const { t } = useTranslation()
  const launchParams = useLaunchParams()
  const telegramUser = launchParams.tgWebAppData?.user as TelegramUser | undefined
  const fullName = [telegramUser?.first_name, telegramUser?.last_name].filter(Boolean).join(' ').trim()
  const avatarName = fullName || me?.user.username || 'User'
  const progressValue = me ? (me.used_tunnels / Math.max(me.max_tunnels, 1)) * 100 : 0

  return (
    <Flex gap="4" p="5" rounded="3xl" bg="bg.panel" borderWidth="1px" align="center">
      <Avatar.Root size="2xl">
        {telegramUser?.photo_url ? <Avatar.Image src={telegramUser.photo_url} /> : null}
        <Avatar.Fallback name={avatarName} />
      </Avatar.Root>

      <VStack align="start" gap="1" flex="1">
        <Heading size="lg">{me ? `@${me.user.username}` : '...'}</Heading>
        <Badge colorPalette={me?.user.status === 'approved' ? 'green' : 'orange'}>{me?.user.status ?? 'pending'}</Badge>
        <Badge variant="subtle">{t('profile.limit', { used: me?.used_tunnels ?? 0, max: me?.max_tunnels ?? 0 })}</Badge>
        <Progress.Root value={progressValue} w="full" maxW="260px" size="sm" rounded="full">
          <Progress.Track>
            <Progress.Range />
          </Progress.Track>
        </Progress.Root>
      </VStack>
    </Flex>
  )
}
