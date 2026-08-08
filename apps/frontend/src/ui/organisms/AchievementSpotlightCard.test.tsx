import { fireEvent, render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import { AchievementSpotlightCard } from "./AchievementSpotlightCard";

const { recap } = personas[0];

test("shows every achievement as a badge", () => {
  render(<AchievementSpotlightCard card={recap.story.cards[3]} recap={recap} />);

  for (const achievement of recap.achievements) {
    expect(screen.getByText(achievement.name)).toBeInTheDocument();
  }
});

test("opens the explanation dialog when a badge is tapped", async () => {
  render(<AchievementSpotlightCard card={recap.story.cards[3]} recap={recap} />);

  const [firstAchievement] = recap.achievements;
  fireEvent.click(screen.getByRole("button", { name: firstAchievement.name }));

  expect(
    await screen.findByText(firstAchievement.explanation),
  ).toBeInTheDocument();
});
