import type { Tunnel } from '../lib/api'

type Props = {
  tunnel: Tunnel
  deleting: boolean
  sending: boolean
  onShowQR: () => void
  onSendConfig: () => void
  onDelete: () => void
}

export function TunnelCard({ tunnel, deleting, sending, onShowQR, onSendConfig, onDelete }: Props) {
  return (
    <article className="tunnel-card">
      <div className="tunnel-card__header">
        <div>
          <p className="eyebrow">Tunnel #{tunnel.id}</p>
          <h3>{tunnel.wg_client_name || 'Создается...'}</h3>
        </div>
        <span className="meta-chip">ID: {tunnel.wg_client_id || 'pending'}</span>
      </div>

      <p className="tunnel-card__timestamp">Создан: {new Date(tunnel.created_at).toLocaleString('ru-RU')}</p>

      <div className="tunnel-card__actions">
        <button className="secondary-button" onClick={onShowQR}>QR</button>
        <button className="secondary-button" disabled={sending} onClick={onSendConfig}>
          {sending ? 'Отправляем...' : 'В чат'}
        </button>
        <button className="danger-button danger-button--soft" disabled={deleting} onClick={onDelete}>
          {deleting ? 'Удаляем...' : 'Удалить'}
        </button>
      </div>
    </article>
  )
}
