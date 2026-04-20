import { Box, Button, Dialog, Portal, Spinner, Text, VStack } from '@chakra-ui/react'
import { useTranslation } from 'react-i18next'

type Props = {
  open: boolean
  loading: boolean
  svg: string
  onClose: () => void
}

export function TunnelQrDialog({ open, loading, svg, onClose }: Props) {
  const { t } = useTranslation()

  return (
    <Dialog.Root open={open} onOpenChange={(details) => !details.open && onClose()} placement="center">
      <Portal>
        <Dialog.Positioner bg="blackAlpha.600">
          <Dialog.Content w="90vw" maxW="560px">
            <Dialog.Header>
              <Dialog.Title>{t('qr.title')}</Dialog.Title>
            </Dialog.Header>
            <Dialog.Body>
              <VStack gap="4">
                <Box w="full" minH="240px" display="grid" placeItems="center" rounded="2xl" bg="white" color="black" p="4">
                  {loading ? <Spinner /> : null}
                  {!loading && svg ? <Box w="280px" maxW="full" dangerouslySetInnerHTML={{ __html: svg }} /> : null}
                </Box>
                <Text color="fg.muted" textAlign="center">{t('qr.body')}</Text>
              </VStack>
            </Dialog.Body>
            <Dialog.Footer>
              <Dialog.ActionTrigger asChild>
                <Button variant="outline" onClick={onClose}>{t('dialog.close')}</Button>
              </Dialog.ActionTrigger>
            </Dialog.Footer>
          </Dialog.Content>
        </Dialog.Positioner>
      </Portal>
    </Dialog.Root>
  )
}
