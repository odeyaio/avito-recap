import { fireEvent, render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import { BehaviorRevealCard } from "./BehaviorRevealCard";

const { recap } = personas[0];

test("shows the primary behavior and its traits", () => {
  render(<BehaviorRevealCard card={recap.story.cards[0]} recap={recap} />);

  expect(
    screen.getByRole("heading", { name: recap.behavior.primary.name }),
  ).toBeInTheDocument();
  expect(screen.getByText(recap.behavior.traits[0].name)).toBeInTheDocument();
});

test("opens the explanation dialog when the primary behavior is tapped", async () => {
  render(<BehaviorRevealCard card={recap.story.cards[0]} recap={recap} />);

  fireEvent.click(
    screen.getByRole("heading", { name: recap.behavior.primary.name }),
  );

  expect(
    await screen.findByText(recap.behavior.primary.explanation),
  ).toBeInTheDocument();
  expect(
    screen.getByText(recap.behavior.primary.evidence[0].description),
  ).toBeInTheDocument();
});

test("opens the explanation dialog when a trait chip is tapped", async () => {
  render(<BehaviorRevealCard card={recap.story.cards[0]} recap={recap} />);

  fireEvent.click(screen.getByText(recap.behavior.traits[0].name));

  expect(
    await screen.findByText(recap.behavior.traits[0].explanation),
  ).toBeInTheDocument();
});
