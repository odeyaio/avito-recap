import Avatar from "@mui/material/Avatar";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useState } from "react";

import type { Recap, StoryCard } from "../../api/generated/model";
import { ExplanationDialog } from "../molecules/ExplanationDialog";
import { StoryCardContent } from "../molecules/StoryCardContent";

export interface AchievementRevealCardProps {
  card: StoryCard;
  recap: Recap;
}

/**
 * The backend emits one `kind: "achievement"` story card per unlocked
 * achievement (see mapper.go's achievementResponse) - `card.id` is the
 * achievement's `code`, so we look up the matching entry in
 * `recap.achievements` for the explanation/evidence the plain StoryCard
 * fields don't carry.
 */
export function AchievementRevealCard({ card, recap }: AchievementRevealCardProps) {
  const [open, setOpen] = useState(false);
  const achievement = recap.achievements.find((item) => item.code === card.id);

  if (!achievement) {
    return <StoryCardContent card={card} recap={recap} />;
  }

  return (
    <Stack spacing={2} sx={{ maxWidth: 480, textAlign: "center", pointerEvents: "none" }}>
      <Stack
        spacing={2}
        onClick={() => setOpen(true)}
        sx={{ alignItems: "center", cursor: "pointer", pointerEvents: "auto" }}
      >
        <Avatar
          src={achievement.image.url}
          alt={achievement.image.alt}
          sx={{ width: 96, height: 96 }}
        >
          {achievement.name.charAt(0)}
        </Avatar>
        <Typography variant="h4" component="h1">
          {achievement.name}
        </Typography>
        <Typography color="text.secondary">{achievement.description}</Typography>
      </Stack>
      <Typography variant="caption" color="text.secondary">
        Нажмите, чтобы узнать, за что эта ачивка
      </Typography>
      <ExplanationDialog
        open={open}
        onClose={() => setOpen(false)}
        title={achievement.name}
        explanation={achievement.explanation}
      />
    </Stack>
  );
}
