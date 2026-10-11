import { defineConfig } from "vite";

const disableProxy = process.env.VITE_DISABLE_PROXY === '1' || process.env.CI === 'true';

export default defineConfig({
  server: {
    port: 5173,
    ...(disableProxy
      ? {}
      : {
          proxy: {
            "/api": "http://localhost:8080",
            "/healthz": "http://localhost:8080",
            "/readyz": "http://localhost:8080",
          },
        }),
  },
  build: {
    sourcemap: true,
    target: "es2022",
  },
});
