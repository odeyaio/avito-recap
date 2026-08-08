import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import type { StoryCard } from "../../api/generated/model";
import { StoryCardRenderer } from "./StoryCardRenderer";

function buildCard(overrides: Partial<StoryCard>): StoryCard {
  return {
    id: "card-1",
    kind: "top_category",
    title: "Заголовок",
    text: "Текст карточки",
    shareable: true,
    ...overrides,
  };
}

test("renders a known kind", () => {
  render(<StoryCardRenderer card={buildCard({ kind: "top_category" })} />);

  expect(screen.getByText("Заголовок")).toBeInTheDocument();
  expect(screen.getByText("Текст карточки")).toBeInTheDocument();
});

test("falls back without crashing for an unrecognized kind", () => {
  render(
    <StoryCardRenderer
      card={buildCard({ kind: "some_future_kind_the_backend_invented" })}
    />,
  );

  expect(screen.getByText("Заголовок")).toBeInTheDocument();
  expect(screen.getByText("Текст карточки")).toBeInTheDocument();
});
