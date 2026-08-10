import { defineConfig } from "orval";

export default defineConfig({
  api: {
    input: "../../contracts/openapi.yaml",
    output: {
      target: "./src/api/generated/client.ts",
      schemas: "./src/api/generated/model",
      client: "react-query",
      httpClient: "fetch",
      clean: true,
      formatter: "prettier",
    },
  },
  apiMocks: {
    input: "../../contracts/openapi.yaml",
    output: {
      target: "./src/api/mocks/generated/handlers.ts",
      client: "fetch",
      httpClient: "fetch",
      mock: true,
      clean: true,
      formatter: "prettier",
    },
  },
});
