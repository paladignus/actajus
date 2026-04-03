import { defineConfig } from "vite";

export default defineConfig({
  build: {
    manifest: true,
    outDir: "../internal/module/web/presentation/http/handler/content/dist",
    emptyOutDir: false,
    rollupOptions: {
      input: "src/main.ts"
    }
  }
});
