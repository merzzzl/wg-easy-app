import {
  Alert,
  Badge,
  Box,
  Button,
  Container,
  Dialog,
  EmptyStateContent,
  EmptyStateDescription,
  EmptyStateIndicator,
  EmptyStateRoot,
  EmptyStateTitle,
  Flex,
  Heading,
  HStack,
  Portal,
  SimpleGrid,
  Skeleton,
  Spinner,
  Stack,
  Text,
  VStack,
} from '@chakra-ui/react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useMemo, useState } from 'react'

import { type Tunnel, useApi } from '../lib/api'

export function MainScreen() {
  const api = useApi()
  const queryClient = useQueryClient()
  const [feedback, setFeedback] = useState<string>('')
  const [selectedTunnel, setSelectedTunnel] = useState<Tunnel | null>(null)
  const [qrSvg, setQrSvg] = useState<string>('')
  const [deletingTunnelId, setDeletingTunnelId] = useState<number | null>(null)

  const meQuery = useQuery({
    queryKey: ['me'],
    queryFn: api.getMe,
  })

  const tunnelsQuery = useQuery({
    queryKey: ['tunnels'],
    queryFn: api.listTunnels,
  })

  const createTunnel = useMutation({
    mutationFn: api.createTunnel,
    onSuccess: async () => {
      setFeedback('Туннель создан.')
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['me'] }),
        queryClient.invalidateQueries({ queryKey: ['tunnels'] }),
      ])
    },
    onError: (error: Error) => {
      setFeedback(error.message)
    },
  })

  const deleteTunnel = useMutation({
    mutationFn: api.deleteTunnel,
    onSuccess: async () => {
      setFeedback('Туннель удален.')
      setDeletingTunnelId(null)
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['me'] }),
        queryClient.invalidateQueries({ queryKey: ['tunnels'] }),
      ])
    },
    onError: (error: Error) => {
      setFeedback(error.message)
    },
  })

  const sendConfig = useMutation({
    mutationFn: api.sendTunnelConfig,
    onSuccess: () => {
      setFeedback('Конфиг отправлен в чат.')
    },
    onError: (error: Error) => {
      setFeedback(error.message)
    },
  })

  const qrMutation = useMutation({
    mutationFn: api.getTunnelQR,
    onSuccess: (payload) => {
      setQrSvg(payload.svg)
    },
    onError: (error: Error) => {
      setFeedback(error.message)
      setSelectedTunnel(null)
      setQrSvg('')
    },
  })

  useEffect(() => {
    if (!feedback) {
      return
    }

    const timeout = window.setTimeout(() => setFeedback(''), 3200)
    return () => window.clearTimeout(timeout)
  }, [feedback])

  const remaining = useMemo(() => {
    if (!meQuery.data) {
      return 0
    }

    return meQuery.data.max_tunnels - meQuery.data.used_tunnels
  }, [meQuery.data])

  const loading = meQuery.isLoading || tunnelsQuery.isLoading
  const bootError = meQuery.error || tunnelsQuery.error

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

  return (
    <Box minH="100vh" py="6">
      <Container maxW="5xl">
        <Stack gap="6">
          <Flex
            direction={{ base: 'column', md: 'row' }}
            justify="space-between"
            align={{ base: 'stretch', md: 'center' }}
            gap="4"
            p="6"
            rounded="3xl"
            bg="bg.panel"
            borderWidth="1px"
            boxShadow="lg"
          >
            <VStack align="start" gap="2">
              <Text fontSize="xs" textTransform="uppercase" color="blue.400" letterSpacing="widest">
                WireGuard Mini App
              </Text>
              <Heading size="2xl">Ваши туннели</Heading>
              <Text color="fg.muted" maxW="2xl">
                Управляйте конфигами в одном окне: создавайте новые туннели, открывайте QR и отправляйте `.conf` прямо в чат.
              </Text>
            </VStack>

            <VStack align={{ base: 'stretch', md: 'end' }} gap="3">
              <Badge colorPalette="blue" px="3" py="2" rounded="full" fontSize="sm">
                Лимит {meQuery.data?.used_tunnels ?? 0}/{meQuery.data?.max_tunnels ?? 0}
              </Badge>
              <Button colorPalette="blue" onClick={handleCreateTunnel} disabled={createTunnel.isPending || loading || remaining <= 0}>
                {createTunnel.isPending ? 'Создаем...' : 'Добавить туннель'}
              </Button>
            </VStack>
          </Flex>

          {feedback ? (
            <Alert.Root status="info" rounded="2xl">
              <Alert.Indicator />
              <Alert.Content>{feedback}</Alert.Content>
            </Alert.Root>
          ) : null}

          {bootError ? (
            <Alert.Root status="error" rounded="2xl">
              <Alert.Indicator />
              <Alert.Content>
                {bootError instanceof Error ? bootError.message : 'Unknown error'}
              </Alert.Content>
            </Alert.Root>
          ) : null}

          <Flex
            justify="space-between"
            align={{ base: 'start', md: 'center' }}
            direction={{ base: 'column', md: 'row' }}
            gap="4"
            p="5"
            rounded="3xl"
            bg="bg.panel"
            borderWidth="1px"
          >
            <VStack align="start" gap="1">
              <Text fontSize="xs" textTransform="uppercase" color="blue.400" letterSpacing="widest">
                Профиль
              </Text>
              <Heading size="lg">{meQuery.data ? `@${meQuery.data.user.username}` : 'Загрузка...'}</Heading>
            </VStack>
            <HStack gap="3" wrap="wrap">
              <Badge colorPalette={meQuery.data?.user.status === 'pending' ? 'orange' : 'green'} px="3" py="2" rounded="full">
                {meQuery.data?.user.status ?? 'pending'}
              </Badge>
              <Badge px="3" py="2" rounded="full" variant="subtle">
                Осталось слотов: {remaining}
              </Badge>
            </HStack>
          </Flex>

          <Box p="5" rounded="3xl" bg="bg.panel" borderWidth="1px">
            <VStack align="start" gap="4">
              <VStack align="start" gap="1">
                <Text fontSize="xs" textTransform="uppercase" color="blue.400" letterSpacing="widest">
                  Список
                </Text>
                <Heading size="lg">Туннели пользователя</Heading>
              </VStack>

              {loading ? (
                <SimpleGrid columns={{ base: 1, md: 2 }} gap="4" w="full">
                  <Skeleton height="160px" rounded="2xl" />
                  <Skeleton height="160px" rounded="2xl" />
                </SimpleGrid>
              ) : null}

              {!loading && tunnelsQuery.data?.length === 0 ? (
                <EmptyStateRoot w="full">
                  <EmptyStateContent>
                    <EmptyStateIndicator>WG</EmptyStateIndicator>
                    <VStack textAlign="center">
                      <EmptyStateTitle>Пока нет ни одного туннеля</EmptyStateTitle>
                      <EmptyStateDescription>
                        Создайте первый конфиг, и он появится здесь с QR и отправкой в чат.
                      </EmptyStateDescription>
                    </VStack>
                  </EmptyStateContent>
                </EmptyStateRoot>
              ) : null}

              {!loading && tunnelsQuery.data?.length ? (
                <SimpleGrid columns={{ base: 1, md: 2 }} gap="4" w="full">
                  {tunnelsQuery.data.map((tunnel) => (
                    <Box key={tunnel.id} p="5" rounded="2xl" borderWidth="1px">
                      <Flex justify="space-between" align="start" gap="3">
                        <VStack align="start" gap="1">
                          <Text fontSize="xs" textTransform="uppercase" color="blue.400" letterSpacing="widest">
                            Tunnel #{tunnel.id}
                          </Text>
                          <Heading size="md">{tunnel.wg_client_name || 'Создается...'}</Heading>
                        </VStack>
                        <Badge variant="subtle">ID: {tunnel.wg_client_id || 'pending'}</Badge>
                      </Flex>

                      <Text mt="4" color="fg.muted">
                        Создан: {new Date(tunnel.created_at).toLocaleString('ru-RU')}
                      </Text>

                      <HStack mt="4" gap="3" wrap="wrap">
                        <Button variant="outline" onClick={() => handleShowQR(tunnel)}>QR</Button>
                        <Button variant="outline" disabled={sendConfig.isPending} onClick={() => sendConfig.mutate(tunnel.id)}>
                          {sendConfig.isPending ? 'Отправляем...' : 'В чат'}
                        </Button>
                        <Button colorPalette="red" variant="subtle" disabled={deleteTunnel.isPending && deletingTunnelId === tunnel.id} onClick={() => handleDeleteTunnel(tunnel)}>
                          {deleteTunnel.isPending && deletingTunnelId === tunnel.id ? 'Удаляем...' : 'Удалить'}
                        </Button>
                      </HStack>
                    </Box>
                  ))}
                </SimpleGrid>
              ) : null}
            </VStack>
          </Box>
        </Stack>
      </Container>

      <Dialog.Root open={selectedTunnel !== null} onOpenChange={(details) => !details.open && setSelectedTunnel(null)} placement="center">
        <Portal>
          <Dialog.Positioner bg="blackAlpha.600">
            <Dialog.Content w="90vw" maxW="560px">
              <Dialog.Header>
                <Dialog.Title>{selectedTunnel?.wg_client_name ?? 'QR'}</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <VStack gap="4">
                  <Box w="full" minH="240px" display="grid" placeItems="center" rounded="2xl" bg="white" color="black" p="4">
                    {qrMutation.isPending ? <Spinner /> : null}
                    {!qrMutation.isPending && qrSvg ? <Box w="280px" maxW="full" dangerouslySetInnerHTML={{ __html: qrSvg }} /> : null}
                  </Box>
                  <Text color="fg.muted" textAlign="center">
                    Откройте QR в приложении WireGuard или AmneziaVPN для быстрого импорта конфигурации.
                  </Text>
                </VStack>
              </Dialog.Body>
              <Dialog.Footer>
                <Dialog.ActionTrigger asChild>
                  <Button variant="outline" onClick={() => {
                    setSelectedTunnel(null)
                    setQrSvg('')
                    qrMutation.reset()
                  }}>
                    Закрыть
                  </Button>
                </Dialog.ActionTrigger>
              </Dialog.Footer>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>

      <Dialog.Root open={deletingTunnelId !== null} onOpenChange={(details) => !details.open && setDeletingTunnelId(null)} placement="center">
        <Portal>
          <Dialog.Positioner bg="blackAlpha.600">
            <Dialog.Content w="90vw" maxW="420px">
              <Dialog.Header>
                <Dialog.Title>Удалить туннель?</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <Text>Конфиг будет удален и из приложения, и из `wg-easy`.</Text>
              </Dialog.Body>
              <Dialog.Footer>
                <Dialog.ActionTrigger asChild>
                  <Button variant="outline" onClick={() => setDeletingTunnelId(null)}>Отмена</Button>
                </Dialog.ActionTrigger>
                <Button colorPalette="red" onClick={confirmDeleteTunnel} disabled={deleteTunnel.isPending}>
                  {deleteTunnel.isPending ? 'Удаляем...' : 'Удалить'}
                </Button>
              </Dialog.Footer>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </Box>
  )
}
