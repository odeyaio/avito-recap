import { fireEvent, render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import type { StoryCard } from "../../api/generated/model";
import { AchievementRevealCard } from "./AchievementRevealCard";

const { recap } = personas[0];
const [firstAchievement] = recap.achievements;
const matchingCard = recap.story.cards.find((card) => card.id === firstAchievement.code)!;

test("shows the matching achievement's name and description", () => {
  render(<AchievementRevealCard card={matchingCard} recap={recap} />);

  expect(
    screen.getByRole("heading", { name: firstAchievement.name }),
  ).toBeInTheDocument();
  expect(screen.getByText(firstAchievement.description)).toBeInTheDocument();
});

test("opens the explanation dialog when tapped", async () => {
  render(<AchievementRevealCard card={matchingCard} recap={recap} />);

  fireEvent.click(screen.getByRole("heading", { name: firstAchievement.name }));

  expect(
    await screen.findByText(firstAchievement.explanation),
  ).toBeInTheDocument();
});

test("falls back to the generic card if no achievement matches the card id", () => {
  const orphanCard: StoryCard = {
    id: "does_not_exist",
    kind: "achievement",
    title: "Заголовок",
    text: "Текст карточки",
    shareable: true,
  };

  render(<AchievementRevealCard card={orphanCard} recap={recap} />);

  expect(screen.getByText("Заголовок")).toBeInTheDocument();
  expect(screen.getByText("Текст карточки")).toBeInTheDocument();
});
