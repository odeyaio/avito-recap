import type { ComponentType } from "react";

import type { Recap, StoryCard } from "../../api/generated/model";
import { StoryCardContent } from "../molecules/StoryCardContent";
import { AchievementSpotlightCard } from "./AchievementSpotlightCard";
import { BehaviorRevealCard } from "./BehaviorRevealCard";

export interface StoryCardComponentProps {
  card: StoryCard;
  recap: Recap;
}

/**
 * `behavior_reveal` and `achievement_spotlight` cross-reference recap.behavior
 * and recap.achievements for explainability; every other known `kind` renders
 * through the same generic visual. An unrecognized `kind` still falls back to
 * the generic content instead of crashing the player, so the backend can
 * introduce new kinds without a matching frontend release.
 */
const CARD_COMPONENTS: Record<string, ComponentType<StoryCardComponentProps>> = {
  intro: StoryCardContent,
  top_category: StoryCardContent,
  activity_streak: StoryCardContent,
  behavior_reveal: BehaviorRevealCard,
  achievement_spotlight: AchievementSpotlightCard,
};

const FALLBACK_CARD_COMPONENT: ComponentType<StoryCardComponentProps> = StoryCardContent;

export interface StoryCardRendererProps {
  card: StoryCard;
  recap: Recap;
}

export function StoryCardRenderer({ card, recap }: StoryCardRendererProps) {
  const Component = CARD_COMPONENTS[card.kind] ?? FALLBACK_CARD_COMPONENT;

  return <Component card={card} recap={recap} />;
}
