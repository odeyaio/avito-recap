import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import { NextActionCta } from "./NextActionCta";

const { nextAction } = personas[0].recap;

test("links to the next action's href", () => {
  render(<NextActionCta nextAction={nextAction} />);

  expect(screen.getByText(nextAction.title)).toBeInTheDocument();
  expect(screen.getByRole("link", { name: "Перейти" })).toHaveAttribute(
    "href",
    nextAction.href,
  );
});
