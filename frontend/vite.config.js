import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:1108',
        changeOrigin: true,
        secure: false
      },
      '/mirror': {
        target: 'http://localhost:1108',
        changeOrigin: true,
        secure: false
      }
    }
  }
})
