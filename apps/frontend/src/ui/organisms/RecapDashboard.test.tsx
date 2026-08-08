import { fireEvent, render, screen } from "@testing-library/react";
import { expect, test, vi } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import { RecapDashboard } from "./RecapDashboard";

const { recap } = personas[0];

test("shows the headline, every achievement, every story card and the next action", () => {
  render(<RecapDashboard recap={recap} onReplay={vi.fn()} />);

  expect(
    screen.getByRole("heading", { name: recap.story.headline }),
  ).toBeInTheDocument();

  for (const achievement of recap.achievements) {
    expect(screen.getByText(achievement.name)).toBeInTheDocument();
  }
  for (const card of recap.story.cards) {
    expect(screen.getByText(card.title)).toBeInTheDocument();
  }

  expect(screen.getByRole("link", { name: "Перейти" })).toHaveAttribute(
    "href",
    recap.nextAction.href,
  );
});

test("opens an explanation dialog when an achievement badge is tapped", async () => {
  render(<RecapDashboard recap={recap} onReplay={vi.fn()} />);

  const [firstAchievement] = recap.achievements;
  fireEvent.click(screen.getByRole("button", { name: firstAchievement.name }));

  expect(
    await screen.findByText(firstAchievement.explanation),
  ).toBeInTheDocument();
});

test("calls onReplay when the replay button is tapped", () => {
  const onReplay = vi.fn();
  render(<RecapDashboard recap={recap} onReplay={onReplay} />);

  fireEvent.click(screen.getByRole("button", { name: "Смотреть ещё раз" }));

  expect(onReplay).toHaveBeenCalledOnce();
});
