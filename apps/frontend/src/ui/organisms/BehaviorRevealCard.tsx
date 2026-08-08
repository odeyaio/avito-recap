import Chip from "@mui/material/Chip";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useState } from "react";

import type { BehaviorMatch, Recap, StoryCard } from "../../api/generated/model";
import { ExplanationDialog } from "../molecules/ExplanationDialog";

export interface BehaviorRevealCardProps {
  card: StoryCard;
  recap: Recap;
}

export function BehaviorRevealCard({ recap }: BehaviorRevealCardProps) {
  const [selected, setSelected] = useState<BehaviorMatch | null>(null);
  const { primary, traits } = recap.behavior;

  return (
    <Stack spacing={2} sx={{ maxWidth: 480, textAlign: "center", pointerEvents: "none" }}>
      <Typography
        variant="h4"
        component="h1"
        onClick={() => setSelected(primary)}
        sx={{ cursor: "pointer", pointerEvents: "auto" }}
      >
        {primary.name}
      </Typography>
      <Typography color="text.secondary">{primary.description}</Typography>
      {traits.length > 0 ? (
        <Stack
          direction="row"
          spacing={1}
          sx={{ flexWrap: "wrap", justifyContent: "center", pointerEvents: "auto" }}
        >
          {traits.map((trait) => (
            <Chip
              key={trait.code}
              label={trait.name}
              variant="outlined"
              onClick={() => setSelected(trait)}
            />
          ))}
        </Stack>
      ) : null}
      <Typography variant="caption" color="text.secondary">
        Нажмите, чтобы узнать почему
      </Typography>
      <ExplanationDialog
        open={selected !== null}
        onClose={() => setSelected(null)}
        title={selected?.name ?? ""}
        explanation={selected?.explanation ?? ""}
        evidence={selected?.evidence}
      />
    </Stack>
  );
}
