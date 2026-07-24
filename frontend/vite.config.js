import { defineConfig } from "vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite"; // <-- Import the plugin

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  base: "./",
  build: {
    outDir: "dist",
    // Optional: Clear the output directory before each build
    emptyOutDir: true,
    chunkSizeWarningLimit: 1000,
  },
});
