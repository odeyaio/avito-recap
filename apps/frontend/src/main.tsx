import { QueryClientProvider } from "@tanstack/react-query";
import CssBaseline from "@mui/material/CssBaseline";
import { ThemeProvider } from "@mui/material/styles";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";

import { queryClient } from "./api/client";
import { routes } from "./routes";
import { theme } from "./theme";

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

const router = createBrowserRouter(routes);

enableMockingIfNeeded().then(() => {
  createRoot(root).render(
    <StrictMode>
      <ThemeProvider theme={theme} defaultMode="system">
        <CssBaseline enableColorScheme />
        <QueryClientProvider client={queryClient}>
          <RouterProvider router={router} />
        </QueryClientProvider>
      </ThemeProvider>
    </StrictMode>,
  );
});
