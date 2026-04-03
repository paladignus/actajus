import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: 5173,
    strictPort: true,
    cors: true
  },
  build: {
    manifest: true,
    outDir: "../internal/module/web/presentation/http/handler/content/dist",
    emptyOutDir: true,
    rollupOptions: {
      input: "src/main.ts"
    }
  }
});
