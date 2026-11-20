import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// The page calls /api/...; the dev server passes those requests on to the
// ticketing service, so the browser sees one origin and CORS never applies.
// Point it elsewhere with TICKETING_URL (default http://localhost:8080).
const target = process.env.TICKETING_URL ?? "http://localhost:8080";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api": { target, changeOrigin: true, rewrite: (path) => path.replace(/^\/api/, "") },
    },
  },
});
