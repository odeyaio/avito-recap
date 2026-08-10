import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";

import type { StoryCard } from "../../api/generated/model";

export interface RecapSummaryCardProps {
  card: StoryCard;
}

export function RecapSummaryCard({ card }: RecapSummaryCardProps) {
  return (
    <Card variant="outlined" sx={{ height: "100%" }}>
      <CardContent>
        <Typography variant="subtitle2">{card.title}</Typography>
        <Typography variant="body2" color="text.secondary">
          {card.text}
        </Typography>
      </CardContent>
    </Card>
  );
}
