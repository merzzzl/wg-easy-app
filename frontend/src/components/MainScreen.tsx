import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useMemo, useState } from 'react'

import { type Tunnel, useApi } from '../lib/api'
import { QRModal } from './QRModal'
import { TunnelCard } from './TunnelCard'

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
    <main className="shell">
      <section className="hero-card">
        <div>
          <p className="eyebrow">WireGuard Mini App</p>
          <h1>Ваши туннели</h1>
          <p className="hero-copy">
            Управляйте конфигами в одном окне: создавайте новые туннели, открывайте QR и отправляйте `.conf` прямо в чат.
          </p>
        </div>
        <div className="hero-metrics">
          <div className="metric-pill">
            <span>Лимит</span>
            <strong>{meQuery.data?.used_tunnels ?? 0}/{meQuery.data?.max_tunnels ?? 0}</strong>
          </div>
          <button
            className="primary-button"
            disabled={createTunnel.isPending || loading || remaining <= 0}
            onClick={handleCreateTunnel}
          >
            {createTunnel.isPending ? 'Создаем...' : 'Добавить туннель'}
          </button>
        </div>
      </section>

      {feedback ? <div className="feedback-banner">{feedback}</div> : null}

      {bootError ? (
        <section className="panel error-panel">
          <h2>Не удалось загрузить данные</h2>
          <p>{bootError instanceof Error ? bootError.message : 'Unknown error'}</p>
        </section>
      ) : null}

      <section className="panel profile-panel">
        <div>
          <p className="eyebrow">Профиль</p>
          <h2>{meQuery.data ? `@${meQuery.data.user.username}` : 'Загрузка...'}</h2>
        </div>
        <div className="profile-meta">
          <span className={`status-chip status-chip--${meQuery.data?.user.status ?? 'pending'}`}>
            {meQuery.data?.user.status ?? 'pending'}
          </span>
          <span className="meta-chip">Осталось слотов: {remaining}</span>
        </div>
      </section>

      <section className="panel list-panel">
        <div className="panel-heading">
          <div>
            <p className="eyebrow">Список</p>
            <h2>Туннели пользователя</h2>
          </div>
        </div>

        {loading ? <div className="skeleton-grid"><div className="skeleton-card" /><div className="skeleton-card" /></div> : null}

        {!loading && tunnelsQuery.data?.length === 0 ? (
          <div className="empty-state empty-state--compact">
            <div className="empty-state__icon">WG</div>
            <h3>Пока нет ни одного туннеля</h3>
            <p>Создайте первый конфиг, и он появится здесь с QR и отправкой в чат.</p>
          </div>
        ) : null}

        {!loading && tunnelsQuery.data?.length ? (
          <div className="tunnel-grid">
            {tunnelsQuery.data.map((tunnel) => (
              <TunnelCard
                key={tunnel.id}
                tunnel={tunnel}
                deleting={deleteTunnel.isPending && deletingTunnelId === tunnel.id}
                sending={sendConfig.isPending}
                onShowQR={() => handleShowQR(tunnel)}
                onSendConfig={() => sendConfig.mutate(tunnel.id)}
                onDelete={() => handleDeleteTunnel(tunnel)}
              />
            ))}
          </div>
        ) : null}
      </section>

      <QRModal
        open={selectedTunnel !== null}
        loading={qrMutation.isPending}
        tunnelName={selectedTunnel?.wg_client_name ?? ''}
        svg={qrSvg}
        onClose={() => {
          setSelectedTunnel(null)
          setQrSvg('')
          qrMutation.reset()
        }}
      />

      {deletingTunnelId !== null ? (
        <div className="modal-backdrop" role="presentation">
          <div className="modal-card modal-card--compact" role="dialog" aria-modal="true">
            <div className="modal-card__header">
              <p className="eyebrow">Удаление</p>
              <h3>Удалить туннель?</h3>
            </div>
            <p className="modal-copy">Конфиг будет удален и из приложения, и из `wg-easy`.</p>
            <div className="modal-actions">
              <button className="ghost-button" onClick={() => setDeletingTunnelId(null)}>
                Отмена
              </button>
              <button className="danger-button" disabled={deleteTunnel.isPending} onClick={confirmDeleteTunnel}>
                {deleteTunnel.isPending ? 'Удаляем...' : 'Удалить'}
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </main>
  )
}
