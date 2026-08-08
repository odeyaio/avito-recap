import Avatar from "@mui/material/Avatar";
import ButtonBase from "@mui/material/ButtonBase";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import type { Achievement } from "../../api/generated/model";

export interface AchievementBadgeProps {
  achievement: Achievement;
  onSelect: (achievement: Achievement) => void;
}

export function AchievementBadge({ achievement, onSelect }: AchievementBadgeProps) {
  return (
    <ButtonBase
      aria-label={achievement.name}
      onClick={() => onSelect(achievement)}
      sx={{ borderRadius: 2, p: 1, width: 96 }}
    >
      <Stack spacing={1} sx={{ alignItems: "center" }}>
        <Avatar src={achievement.image.url} alt={achievement.image.alt} sx={{ width: 56, height: 56 }}>
          {achievement.name.charAt(0)}
        </Avatar>
        <Typography variant="caption" sx={{ textAlign: "center" }}>
          {achievement.name}
        </Typography>
      </Stack>
    </ButtonBase>
  );
}
