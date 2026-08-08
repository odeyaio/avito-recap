import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useState } from "react";

import type { Achievement, Recap, StoryCard } from "../../api/generated/model";
import { AchievementBadge } from "../molecules/AchievementBadge";
import { ExplanationDialog } from "../molecules/ExplanationDialog";

export interface AchievementSpotlightCardProps {
  card: StoryCard;
  recap: Recap;
}

export function AchievementSpotlightCard({ recap }: AchievementSpotlightCardProps) {
  const [selected, setSelected] = useState<Achievement | null>(null);

  return (
    <Stack spacing={2} sx={{ maxWidth: 480, textAlign: "center", pointerEvents: "none" }}>
      <Typography variant="h4" component="h1">
        Открытые ачивки
      </Typography>
      <Stack
        direction="row"
        spacing={2}
        sx={{ flexWrap: "wrap", justifyContent: "center", pointerEvents: "auto" }}
      >
        {recap.achievements.map((achievement) => (
          <AchievementBadge
            key={achievement.code}
            achievement={achievement}
            onSelect={setSelected}
          />
        ))}
      </Stack>
      <Typography variant="caption" color="text.secondary">
        Нажмите на ачивку, чтобы узнать, за что она
      </Typography>
      <ExplanationDialog
        open={selected !== null}
        onClose={() => setSelected(null)}
        title={selected?.name ?? ""}
        explanation={selected?.explanation ?? ""}
      />
    </Stack>
  );
}
