import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { init, isTMA, miniApp, themeParams } from '@tma.js/sdk-react'
import { Provider } from './components/ui/provider'

import App from './App'
import WebApp from './WebApp'
import './index.css'

const initializeTelegramSDK = () => {
  try {
    init()
    themeParams.mount()
    miniApp.mount()
    if (miniApp.ready.isAvailable()) {
      miniApp.ready()
    }
  } catch (error) {
    console.error('Failed to initialize Telegram SDK:', error)
  }
}

const inTelegram = isTMA()

if (inTelegram) {
  initializeTelegramSDK()
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Provider>
      {inTelegram ? <App /> : <WebApp />}
    </Provider>
  </StrictMode>,
)
