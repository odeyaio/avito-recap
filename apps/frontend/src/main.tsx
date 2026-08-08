import { QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";

import { App } from "./App";
import { queryClient } from "./api/client";
import "./styles.css";

async function enableMockingIfNeeded() {
  if (!import.meta.env.DEV || import.meta.env.VITE_API_MOCK !== "true") {
    return;
  }

  const { worker } = await import("./api/mocks/browser");
  await worker.start({ onUnhandledRequest: "bypass" });
}

const root = document.getElementById("root");

if (!root) {
  throw new Error("Root element is missing");
}

enableMockingIfNeeded().then(() => {
  createRoot(root).render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </QueryClientProvider>
    </StrictMode>,
  );
});
