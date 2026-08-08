import { QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { createMemoryRouter, RouterProvider } from "react-router-dom";
import { expect, test } from "vitest";

import { queryClient } from "./api/client";
import { routes } from "./routes";

const EXPLORER_PROFILE_ID = "11111111-1111-4111-8111-111111111111";

function renderAt(initialPath: string) {
  const router = createMemoryRouter(routes, { initialEntries: [initialPath] });

  return render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  );
}

test("renders the intro screen at /", async () => {
  renderAt("/");

  expect(
    await screen.findByRole("heading", { name: "Ваши Итоги года на Авито" }),
  ).toBeInTheDocument();
});

test("loads mocked profiles on the profile picker screen", async () => {
  renderAt("/profiles");

  expect(await screen.findByText(/Алексей · Омск/)).toBeInTheDocument();
  expect(screen.getByText(/Марина · Санкт-Петербург/)).toBeInTheDocument();
  expect(screen.getByText(/Игорь · Новосибирск/)).toBeInTheDocument();
});

test("generating a recap redirects to the recap screen", async () => {
  renderAt(`/profiles/${EXPLORER_PROFILE_ID}/generating?year=2025`);

  expect(
    await screen.findByRole("heading", {
      name: "Ваш тип года — «Исследователь»",
    }),
  ).toBeInTheDocument();
});

test("generating for an unknown profile renders the error boundary", async () => {
  renderAt("/profiles/00000000-0000-4000-8000-000000000000/generating?year=2025");

  expect(await screen.findByText("Профиль не найден")).toBeInTheDocument();
});
