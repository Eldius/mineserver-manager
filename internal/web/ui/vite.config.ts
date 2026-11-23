import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// The build output goes straight into internal/web/static, which the Go
// API server embeds (see internal/web/api.go). Commit the built assets:
// Go builds must work without node installed.
export default defineConfig({
  plugins: [react()],
  base: './',
  build: {
    outDir: '../static',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
