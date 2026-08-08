import type { ComponentType } from "react";

import type { StoryCard } from "../../api/generated/model";
import { StoryCardContent } from "../molecules/StoryCardContent";

/**
 * Every known `kind` renders through the same visual today. As specialized
 * cards (behavior reveal, achievement spotlight, ...) get their own designs
 * in later branches, register them here — an unrecognized `kind` still falls
 * back to the generic content instead of crashing the player.
 */
const CARD_COMPONENTS: Record<string, ComponentType<{ card: StoryCard }>> = {
  intro: StoryCardContent,
  top_category: StoryCardContent,
  activity_streak: StoryCardContent,
  behavior_reveal: StoryCardContent,
  achievement_spotlight: StoryCardContent,
};

const FALLBACK_CARD_COMPONENT: ComponentType<{ card: StoryCard }> = StoryCardContent;

export interface StoryCardRendererProps {
  card: StoryCard;
}

export function StoryCardRenderer({ card }: StoryCardRendererProps) {
  const Component = CARD_COMPONENTS[card.kind] ?? FALLBACK_CARD_COMPONENT;

  return <Component card={card} />;
}
