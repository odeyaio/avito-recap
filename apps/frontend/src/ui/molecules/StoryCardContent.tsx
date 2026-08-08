import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import type { StoryCard } from "../../api/generated/model";

export interface StoryCardContentProps {
  card: StoryCard;
}

export function StoryCardContent({ card }: StoryCardContentProps) {
  return (
    <Stack spacing={2} sx={{ maxWidth: 480, textAlign: "center", pointerEvents: "none" }}>
      {card.mediaUrl ? (
        <Box
          component="img"
          src={card.mediaUrl}
          alt=""
          sx={{ width: "100%", borderRadius: 2 }}
        />
      ) : null}
      <Typography variant="h4" component="h1">
        {card.title}
      </Typography>
      <Typography color="text.secondary">{card.text}</Typography>
    </Stack>
  );
}
