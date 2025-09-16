import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: "./",
  plugins: [react()],
  server: {
    port: 3000,
    strictPort: true,
    open: false,
    cors: true,
    hmr: {
      protocol: "ws",
      host: "localhost",
      port: 3000
    },
    fs: {
      strict: false
    },
    watch: {
      usePolling: true
    }
  },
  build: {
    outDir: "../backend/assets/frontend/dist",
    emptyOutDir: true,
    sourcemap: true,
    minify: true,
    target: "esnext",
    rollupOptions: {
      input: {
        main: "./index.html"
      },
      output: {
        entryFileNames: "assets/[name]-[hash].js",
        chunkFileNames: "assets/[name]-[hash].js",
        assetFileNames: "assets/[name]-[hash][extname]"
      }
    }
  },
  optimizeDeps: {
    esbuildOptions: {
      target: "es2020"
    }
  },
  define: {
    "process.env": {}
  },
  resolve: {
    alias: {
      "@": "/src"
    }
  }
});
