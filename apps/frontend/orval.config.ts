import { defineConfig } from "orval";

export default defineConfig({
  api: {
    input: "../../contracts/openapi.yaml",
    output: {
      target: "./src/api/generated/client.ts",
      schemas: "./src/api/generated/model",
      client: "fetch",
      clean: true,
      formatter: "prettier",
    },
  },
});
