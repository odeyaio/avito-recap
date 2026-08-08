import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useState } from "react";

import type { Achievement, Recap } from "../../api/generated/model";
import { AchievementBadge } from "../molecules/AchievementBadge";
import { ExplanationDialog } from "../molecules/ExplanationDialog";
import { RecapSummaryCard } from "../molecules/RecapSummaryCard";
import { NextActionCta } from "./NextActionCta";

export interface RecapDashboardProps {
  recap: Recap;
  onReplay: () => void;
}

export function RecapDashboard({ recap, onReplay }: RecapDashboardProps) {
  const [selected, setSelected] = useState<Achievement | null>(null);

  return (
    <Stack spacing={3} sx={{ width: "100%" }}>
      <Stack spacing={1}>
        <Typography variant="h4" component="h1">
          {recap.story.headline}
        </Typography>
        <Typography color="text.secondary">{recap.story.summary}</Typography>
      </Stack>

      {recap.achievements.length > 0 ? (
        <Stack spacing={1}>
          <Typography variant="h6" component="h2">
            Ачивки
          </Typography>
          <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
            {recap.achievements.map((achievement) => (
              <AchievementBadge
                key={achievement.code}
                achievement={achievement}
                onSelect={setSelected}
              />
            ))}
          </Stack>
        </Stack>
      ) : null}

      <Stack spacing={1}>
        <Typography variant="h6" component="h2">
          Карточки года
        </Typography>
        <Box
          sx={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))",
            gap: 2,
          }}
        >
          {recap.story.cards.map((card) => (
            <RecapSummaryCard key={card.id} card={card} />
          ))}
        </Box>
      </Stack>

      <NextActionCta nextAction={recap.nextAction} />

      <Button variant="outlined" onClick={onReplay} sx={{ alignSelf: "center" }}>
        Смотреть ещё раз
      </Button>

      <ExplanationDialog
        open={selected !== null}
        onClose={() => setSelected(null)}
        title={selected?.name ?? ""}
        explanation={selected?.explanation ?? ""}
      />
    </Stack>
  );
}
