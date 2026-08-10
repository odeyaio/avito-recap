import { expect, test } from "vitest";

import { personas } from "../../api/mocks/fixtures/personas";
import { buildStorySlides } from "./buildStorySlides";

const { recap } = personas[0];

test("opens with an intro slide, then a behavior slide, then the backend's own cards", () => {
  const slides = buildStorySlides(recap);

  expect(slides[0]).toMatchObject({
    kind: "intro",
    title: recap.story.headline,
    text: recap.story.summary,
  });
  expect(slides[1]).toMatchObject({
    kind: "behavior_reveal",
    title: recap.behavior.primary.name,
    text: recap.behavior.primary.description,
  });
  expect(slides.slice(2)).toEqual(recap.story.cards);
});
