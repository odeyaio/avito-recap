import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { expect, test } from "vitest";

import { App } from "./App";

function renderApp(initialPath: string) {
  const queryClient = new QueryClient();

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialPath]}>
        <App />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

test("renders the intro screen at /", () => {
  renderApp("/");

  expect(
    screen.getByRole("heading", { name: "Ваши Итоги года на Авито" }),
  ).toBeInTheDocument();
});

test("renders the generating screen with the selected profile id", () => {
  renderApp("/profiles/11111111-1111-4111-8111-111111111111/generating");

  expect(
    screen.getByText(/11111111-1111-4111-8111-111111111111/),
  ).toBeInTheDocument();
});

test("loads mocked profiles on the profile picker screen", async () => {
  renderApp("/profiles");

  expect(
    await screen.findByText(/Алексей · Омск/),
  ).toBeInTheDocument();
  expect(screen.getByText(/Марина · Санкт-Петербург/)).toBeInTheDocument();
  expect(screen.getByText(/Игорь · Новосибирск/)).toBeInTheDocument();
});
