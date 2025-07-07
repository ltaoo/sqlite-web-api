import path from "path";

import { UserConfigExport, defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";

const config = defineConfig(({ mode }) => {
  return {
    base: "/",
    plugins: [solidPlugin()],
    resolve: {
      alias: {
        "hls.js": "hls.js/dist/hls.min.js",
        "lucide-solid": require.resolve("lucide-solid").replace("cjs", "esm"),
        "@": path.resolve(__dirname, "./src"),
      },
    },
    esbuild: {
      drop: mode === "production" ? ["console", "debugger"] : [],
    },
    build: {
      target: "esnext",
      rollupOptions: {
        output: {
          manualChunks(filepath) {
            // if (filepath.includes("hls.js")) {
            //   return "hls";
            // }
            if (filepath.includes("node_modules") && !filepath.includes("hls")) {
              return "vendor";
            }
          },
        },
      },
    },
  };
});

export default config;
