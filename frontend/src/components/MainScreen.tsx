import { Alert, Box, Button, Container, Stack } from '@chakra-ui/react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { type Tunnel, useApi } from '../lib/api'
import { ApprovalPending } from './ApprovalPending'
import { TunnelDeleteDialog } from './TunnelDeleteDialog'
import { TunnelList } from './TunnelList'
import { TunnelQrDialog } from './TunnelQrDialog'
import { UserSummary } from './UserSummary'
import { toaster } from './ui/toaster'

export function MainScreen() {
  const { t } = useTranslation()
  const api = useApi()
  const queryClient = useQueryClient()
  const [selectedTunnel, setSelectedTunnel] = useState<Tunnel | null>(null)
  const [qrSvg, setQrSvg] = useState<string>('')
  const [deletingTunnelId, setDeletingTunnelId] = useState<number | null>(null)

  const meQuery = useQuery({
    queryKey: ['me'],
    queryFn: api.getMe,
  })

  const isApproved = meQuery.data?.user.status === 'approved'

  const tunnelsQuery = useQuery({
    queryKey: ['tunnels'],
    queryFn: api.listTunnels,
    enabled: isApproved,
  })

  const createTunnel = useMutation({
    mutationFn: api.createTunnel,
    onSuccess: async () => {
      toaster.create({ title: t('feedback.created'), type: 'success' })
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['me'] }),
        queryClient.invalidateQueries({ queryKey: ['tunnels'] }),
      ])
    },
    onError: (error: Error) => {
      toaster.create({ title: error.message, type: 'error' })
    },
  })

  const deleteTunnel = useMutation({
    mutationFn: api.deleteTunnel,
    onSuccess: async () => {
      toaster.create({ title: t('feedback.deleted'), type: 'success' })
      setDeletingTunnelId(null)
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['me'] }),
        queryClient.invalidateQueries({ queryKey: ['tunnels'] }),
      ])
    },
    onError: (error: Error) => {
      toaster.create({ title: error.message, type: 'error' })
    },
  })

  const sendConfig = useMutation({
    mutationFn: api.sendTunnelConfig,
    onSuccess: () => {
      toaster.create({ title: t('feedback.configSent'), type: 'success' })
    },
    onError: (error: Error) => {
      toaster.create({ title: error.message, type: 'error' })
    },
  })

  const qrMutation = useMutation({
    mutationFn: api.getTunnelQR,
    onSuccess: (payload) => {
      setQrSvg(payload.svg)
    },
    onError: (error: Error) => {
      toaster.create({ title: error.message, type: 'error' })
      setSelectedTunnel(null)
      setQrSvg('')
    },
  })

  const remaining = useMemo(() => {
    if (!meQuery.data) {
      return 0
    }

    return meQuery.data.max_tunnels - meQuery.data.used_tunnels
  }, [meQuery.data])

  const loading = meQuery.isLoading || (isApproved && tunnelsQuery.isLoading)
  const bootError = meQuery.error || (isApproved ? tunnelsQuery.error : null)

  const handleCreateTunnel = () => {
    createTunnel.mutate()
  }

  const handleShowQR = (tunnel: Tunnel) => {
    setSelectedTunnel(tunnel)
    setQrSvg('')
    qrMutation.mutate(tunnel.id)
  }

  const handleDeleteTunnel = (tunnel: Tunnel) => {
    setDeletingTunnelId(tunnel.id)
  }

  const confirmDeleteTunnel = () => {
    if (deletingTunnelId == null) {
      return
    }

    deleteTunnel.mutate(deletingTunnelId)
  }

  if (meQuery.data && !isApproved) {
    return <ApprovalPending me={meQuery.data} />
  }

  return (
    <Box minH="100vh" py="6">
      <Container maxW="container.md">
        <Stack gap="6">
        <UserSummary me={meQuery.data} />

          {bootError ? (
            <Alert.Root status="error" rounded="2xl">
              <Alert.Indicator />
              <Alert.Content>
                {t('error.load')}: {bootError instanceof Error ? bootError.message : 'Unknown error'}
              </Alert.Content>
            </Alert.Root>
          ) : null}

          <Button colorPalette="blue" size="lg" onClick={handleCreateTunnel} disabled={createTunnel.isPending || loading || remaining <= 0}>
            {createTunnel.isPending ? t('actions.creating') : t('actions.create')}
          </Button>

          <TunnelList
            tunnels={tunnelsQuery.data ?? []}
            loading={loading}
            deletingTunnelId={deletingTunnelId}
            deletePending={deleteTunnel.isPending}
            sendPending={sendConfig.isPending}
            onShowQR={handleShowQR}
            onSendConfig={(tunnel) => sendConfig.mutate(tunnel.id)}
            onDelete={handleDeleteTunnel}
          />
        </Stack>
      </Container>

      <TunnelQrDialog
        open={selectedTunnel !== null}
        loading={qrMutation.isPending}
        svg={qrSvg}
        onClose={() => {
          setSelectedTunnel(null)
          setQrSvg('')
          qrMutation.reset()
        }}
      />

      <TunnelDeleteDialog
        open={deletingTunnelId !== null}
        loading={deleteTunnel.isPending}
        onClose={() => setDeletingTunnelId(null)}
        onConfirm={confirmDeleteTunnel}
      />
    </Box>
  )
}
