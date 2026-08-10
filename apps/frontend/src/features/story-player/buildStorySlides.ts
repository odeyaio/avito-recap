import type { Recap, StoryCard } from "../../api/generated/model";

/**
 * The backend's `story.cards` today only ever contains `kind: "achievement"`
 * entries (see apps/backend mapper.go) - there's no headline/summary or
 * behavior card in the sequence, even though both are meant to open the
 * story. Synthesize them client-side from fields the contract already
 * guarantees (`story.headline`/`summary`, `behavior.primary`) instead of
 * waiting on a backend change.
 */
export function buildStorySlides(recap: Recap): StoryCard[] {
  const introCard: StoryCard = {
    id: "intro",
    kind: "intro",
    title: recap.story.headline,
    text: recap.story.summary,
    shareable: true,
  };

  const behaviorCard: StoryCard = {
    id: "behavior",
    kind: "behavior_reveal",
    title: recap.behavior.primary.name,
    text: recap.behavior.primary.description,
    shareable: false,
  };

  return [introCard, behaviorCard, ...recap.story.cards];
}
