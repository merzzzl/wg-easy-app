import {
  Badge,
  Box,
  Button,
  EmptyStateContent,
  EmptyStateDescription,
  EmptyStateIndicator,
  EmptyStateRoot,
  EmptyStateTitle,
  Flex,
  Heading,
  HStack,
  Separator,
  Skeleton,
  Stack,
  Text,
  VStack,
} from '@chakra-ui/react'
import { useTranslation } from 'react-i18next'

import type { Tunnel } from '../lib/api'

type Props = {
  tunnels: Tunnel[]
  loading: boolean
  deletingTunnelId: number | null
  deletePending: boolean
  sendPending: boolean
  onShowQR: (tunnel: Tunnel) => void
  onSendConfig: (tunnel: Tunnel) => void
  onDelete: (tunnel: Tunnel) => void
}

export function TunnelList({
  tunnels,
  loading,
  deletingTunnelId,
  deletePending,
  sendPending,
  onShowQR,
  onSendConfig,
  onDelete,
}: Props) {
  const { t } = useTranslation()

  return (
    <VStack align="start" gap="4" w="full">
      <Heading size="md">{t('list.title')}</Heading>

      {loading ? (
        <Stack w="full" gap="3">
          <Skeleton height="62px" rounded="xl" />
          <Skeleton height="62px" rounded="xl" />
          <Skeleton height="62px" rounded="xl" />
        </Stack>
      ) : null}

      {!loading && tunnels.length === 0 ? (
        <EmptyStateRoot w="full">
          <EmptyStateContent>
            <EmptyStateIndicator>WG</EmptyStateIndicator>
            <VStack textAlign="center">
              <EmptyStateTitle>{t('list.emptyTitle')}</EmptyStateTitle>
              <EmptyStateDescription>{t('list.emptyBody')}</EmptyStateDescription>
            </VStack>
          </EmptyStateContent>
        </EmptyStateRoot>
      ) : null}

      {!loading && tunnels.length > 0 ? (
        <Stack w="full" gap="0">
          {tunnels.map((tunnel) => (
            <Box key={tunnel.id} py="4">
              <Flex direction={{ base: 'column', md: 'row' }} justify="space-between" align={{ base: 'start', md: 'center' }} gap="3">
                <VStack align="start" gap="1">
                  <HStack gap="2" wrap="wrap">
                    <Text fontWeight="semibold">{tunnel.wg_client_name || 'pending'}</Text>
                    <Badge variant="subtle">#{tunnel.id}</Badge>
                  </HStack>
                  <Text fontSize="sm" color="fg.muted">
                    {t('tunnel.createdAt', { value: new Date(tunnel.created_at).toLocaleString() })}
                  </Text>
                </VStack>

                <HStack gap="2" wrap="wrap">
                  <Button size="sm" variant="outline" onClick={() => onShowQR(tunnel)}>{t('actions.showQr')}</Button>
                  <Button size="sm" variant="outline" disabled={sendPending} onClick={() => onSendConfig(tunnel)}>
                    {sendPending ? t('actions.sendingConfig') : t('actions.sendConfig')}
                  </Button>
                  <Button
                    size="sm"
                    colorPalette="red"
                    variant="subtle"
                    disabled={deletePending && deletingTunnelId === tunnel.id}
                    onClick={() => onDelete(tunnel)}
                  >
                    {deletePending && deletingTunnelId === tunnel.id ? t('actions.deleting') : t('actions.delete')}
                  </Button>
                </HStack>
              </Flex>
              <Separator mt="4" />
            </Box>
          ))}
        </Stack>
      ) : null}
    </VStack>
  )
}
