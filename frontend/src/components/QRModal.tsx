type Props = {
  open: boolean
  loading: boolean
  tunnelName: string
  svg: string
  onClose: () => void
}

export function QRModal({ open, loading, tunnelName, svg, onClose }: Props) {
  if (!open) {
    return null
  }

  return (
    <div className="modal-backdrop" role="presentation">
      <div className="modal-card" role="dialog" aria-modal="true" aria-label="QR code">
        <div className="modal-card__header">
          <div>
            <p className="eyebrow">QR</p>
            <h3>{tunnelName}</h3>
          </div>
          <button className="icon-button" onClick={onClose} aria-label="Закрыть">
            ×
          </button>
        </div>

        <div className="qr-surface">
          {loading ? <div className="qr-loader">Загружаем QR...</div> : null}
          {!loading && svg ? <div className="qr-markup" dangerouslySetInnerHTML={{ __html: svg }} /> : null}
        </div>

        <p className="modal-copy">Откройте QR в приложении WireGuard или AmneziaVPN для быстрого импорта конфигурации.</p>
      </div>
    </div>
  )
}
