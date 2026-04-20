import { ChakraProvider, createSystem, defaultConfig, defineConfig } from '@chakra-ui/react'
import type { PropsWithChildren } from 'react'

const config = defineConfig({
  theme: {
    semanticTokens: {
      colors: {
        bg: {
          canvas: {
            value: {
              _light: '#f4f7fb',
              _dark: '#0b1020',
            },
          },
          panel: {
            value: {
              _light: '#ffffff',
              _dark: '#11182b',
            },
          },
        },
      },
    },
  },
  globalCss: {
    'html, body': {
      margin: 0,
      minHeight: '100%',
      bg: 'bg.canvas',
    },
    '#root': {
      minHeight: '100vh',
    },
  },
})

const system = createSystem(defaultConfig, config)

export function Provider({ children }: PropsWithChildren) {
  return <ChakraProvider value={system}>{children}</ChakraProvider>
}
