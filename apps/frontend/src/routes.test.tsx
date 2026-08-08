import { QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen } from "@testing-library/react";
import { createMemoryRouter, RouterProvider } from "react-router-dom";
import { expect, test } from "vitest";

import { queryClient } from "./api/client";
import {
  GENERATION_UNAVAILABLE_PROFILE_ID,
  INSUFFICIENT_ACTIVITY_PROFILE_ID,
} from "./api/mocks/handlers";
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

test("steps through story cards via tap and keyboard", async () => {
  renderAt(`/profiles/${EXPLORER_PROFILE_ID}/generating?year=2025`);

  await screen.findByRole("heading", {
    name: "Ваш тип года — «Исследователь»",
  });

  fireEvent.click(screen.getByRole("button", { name: "Следующая карточка" }));
  expect(
    await screen.findByRole("heading", {
      name: "Электроника — ваша главная тема",
    }),
  ).toBeInTheDocument();

  fireEvent.keyDown(window, { key: "ArrowRight" });
  expect(
    await screen.findByRole("heading", { name: "18 дней подряд" }),
  ).toBeInTheDocument();

  fireEvent.keyDown(window, { key: "ArrowLeft" });
  expect(
    await screen.findByRole("heading", {
      name: "Электроника — ваша главная тема",
    }),
  ).toBeInTheDocument();
});

test("opens an explanation dialog for a behavior card reached through the player", async () => {
  renderAt(`/profiles/${EXPLORER_PROFILE_ID}/generating?year=2025`);

  await screen.findByRole("heading", {
    name: "Ваш тип года — «Исследователь»",
  });

  const next = screen.getByRole("button", { name: "Следующая карточка" });
  fireEvent.click(next); // top_category
  fireEvent.click(next); // activity_streak
  fireEvent.click(next); // behavior_reveal

  const behaviorHeading = await screen.findByRole("heading", {
    name: "Исследователь",
  });
  fireEvent.click(behaviorHeading);

  expect(
    await screen.findByText(
      "Вы посмотрели 132 уникальных объявления в 9 разных категориях и почти не выходили на контакт с продавцами — вам было интереснее изучать, чем покупать.",
    ),
  ).toBeInTheDocument();
});

test("generating for an unknown profile renders the error boundary", async () => {
  renderAt("/profiles/00000000-0000-4000-8000-000000000000/generating?year=2025");

  expect(await screen.findByText("Профиль не найден")).toBeInTheDocument();
  expect(screen.queryByRole("button", { name: "Повторить" })).not.toBeInTheDocument();
  expect(
    screen.getByRole("link", { name: "Выбрать другой профиль" }),
  ).toBeInTheDocument();
});

test("insufficient activity renders a friendly non-retryable error", async () => {
  renderAt(
    `/profiles/${INSUFFICIENT_ACTIVITY_PROFILE_ID}/generating?year=2025`,
  );

  expect(
    await screen.findByText("Недостаточно активности"),
  ).toBeInTheDocument();
  expect(screen.queryByRole("button", { name: "Повторить" })).not.toBeInTheDocument();
});

test("generation unavailable renders a retryable error", async () => {
  renderAt(
    `/profiles/${GENERATION_UNAVAILABLE_PROFILE_ID}/generating?year=2025`,
  );

  expect(
    await screen.findByText("Сервис генерации недоступен"),
  ).toBeInTheDocument();
  expect(
    screen.getByRole("button", { name: "Повторить" }),
  ).toBeInTheDocument();
});
