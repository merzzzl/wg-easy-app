import { Button, Dialog, Portal, Text } from '@chakra-ui/react'
import { useTranslation } from 'react-i18next'

type Props = {
  open: boolean
  loading: boolean
  onClose: () => void
  onConfirm: () => void
}

export function TunnelDeleteDialog({ open, loading, onClose, onConfirm }: Props) {
  const { t } = useTranslation()

  return (
    <Dialog.Root open={open} onOpenChange={(details) => !details.open && onClose()} placement="center">
      <Portal>
        <Dialog.Positioner bg="blackAlpha.600">
          <Dialog.Content w="90vw" maxW="420px">
            <Dialog.Header>
              <Dialog.Title>{t('dialog.deleteTitle')}</Dialog.Title>
            </Dialog.Header>
            <Dialog.Body>
              <Text>{t('dialog.deleteBody')}</Text>
            </Dialog.Body>
            <Dialog.Footer>
              <Dialog.ActionTrigger asChild>
                <Button variant="outline" onClick={onClose}>{t('dialog.cancel')}</Button>
              </Dialog.ActionTrigger>
              <Button colorPalette="red" onClick={onConfirm} disabled={loading}>
                {loading ? t('actions.deleting') : t('actions.delete')}
              </Button>
            </Dialog.Footer>
          </Dialog.Content>
        </Dialog.Positioner>
      </Portal>
    </Dialog.Root>
  )
}
