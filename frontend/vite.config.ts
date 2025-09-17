import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import wasm from "vite-plugin-wasm";

export default defineConfig({
  plugins: [react(), wasm()],
  server: {
    fs: {
      strict: false
    }
  },
  build: {
    target: "esnext"
  },
  define: {
    "process.env": {}
  }
});
