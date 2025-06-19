/* eslint-disable @typescript-eslint/ban-ts-comment */
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import * as path from "path";
import { ConfigEnv, defineConfig } from "vite";
import i18nextLoader from "vite-plugin-i18next-loader";
function getDefines(env: ConfigEnv) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const defines: Record<string, any> = {};

  defines["process.env.NODE_ENV"] =
    env.mode === "development" ? `"development"` : `"production"`; // A lot of libraries need this constant to know if we're in production
  return defines;
}

// https://vite.dev/config/
// @ts-ignore
export default defineConfig((env) => ({
  define: getDefines(env),
  experimental: {
    renderBuiltUrl(filename) {
      return {
        runtime: `window.__fransAssetUrl('${filename}')`,
        relative: true,
      };
    },
  },
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "/client"),
      // https://github.com/tabler/tabler-icons/issues/1233#issuecomment-2428245119
      "@tabler/icons-react": "@tabler/icons-react/dist/esm/icons/index.mjs",
    },
  },
  // @ts-ignore
  plugins: [
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
      routesDirectory: path.resolve(__dirname, "./client/routes"),
      generatedRouteTree: path.resolve(__dirname, "./client/routeTree.gen.ts"),
    }),
    react(),
    i18nextLoader({ paths: ["./locales"], namespaceResolution: "basename" }),
  ],
  server: {
    port: 3000,
  },
  build: {
    emptyOutDir: true,
    manifest: true,
    outDir: path.resolve(__dirname, "pkg/routes/client/assets"),
    rollupOptions: {
      input: "client/main.tsx",
    },
  },
}));
