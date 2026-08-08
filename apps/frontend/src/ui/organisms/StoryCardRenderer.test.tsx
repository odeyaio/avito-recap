import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import type { StoryCard } from "../../api/generated/model";
import { StoryCardRenderer } from "./StoryCardRenderer";

const { recap } = personas[0];

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
  render(
    <StoryCardRenderer card={buildCard({ kind: "top_category" })} recap={recap} />,
  );

  expect(screen.getByText("Заголовок")).toBeInTheDocument();
  expect(screen.getByText("Текст карточки")).toBeInTheDocument();
});

test("falls back without crashing for an unrecognized kind", () => {
  render(
    <StoryCardRenderer
      card={buildCard({ kind: "some_future_kind_the_backend_invented" })}
      recap={recap}
    />,
  );

  expect(screen.getByText("Заголовок")).toBeInTheDocument();
  expect(screen.getByText("Текст карточки")).toBeInTheDocument();
});
