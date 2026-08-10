import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";

import type { ShareCard } from "../../api/generated/model";
import { ShareCardPreview } from "./ShareCardPreview";

test("renders the share card's title and subtitle", () => {
  const shareCard: ShareCard = {
    title: "Алексей — Исследователь года",
    subtitle: "132 объявления, 9 категорий, 11 активных месяцев",
  };

  render(<ShareCardPreview shareCard={shareCard} />);

  expect(screen.getByText(shareCard.title)).toBeInTheDocument();
  expect(screen.getByText(shareCard.subtitle)).toBeInTheDocument();
});
