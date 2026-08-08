import Typography from "@mui/material/Typography";
import { useState } from "react";
import { useParams } from "react-router-dom";

import { useGetRecap } from "../api/generated/client";
import type { StoryCard } from "../api/generated/model";
import { RecapDashboard } from "../ui/organisms/RecapDashboard";
import { ScreenLayout } from "../ui/templates/ScreenLayout";
import { StoryPlayerLayout } from "../ui/templates/StoryPlayerLayout";

export function RecapPage() {
  const { recapId } = useParams<{ recapId: string }>();
  const { data: response } = useGetRecap(recapId ?? "");
  const [showDashboard, setShowDashboard] = useState(false);

  if (!response || response.status !== 200) {
    return (
      <ScreenLayout>
        <Typography>Не удалось загрузить recap.</Typography>
      </ScreenLayout>
    );
  }

  if (showDashboard) {
    return (
      <ScreenLayout key={recapId}>
        <RecapDashboard
          recap={response.data}
          onReplay={() => setShowDashboard(false)}
        />
      </ScreenLayout>
    );
  }

  const { story } = response.data;
  const introCard: StoryCard = {
    id: "intro",
    kind: "intro",
    title: story.headline,
    text: story.summary,
    shareable: true,
  };
  const slides = [introCard, ...story.cards];

  return (
    <StoryPlayerLayout
      key={recapId}
      cards={slides}
      recap={response.data}
      onComplete={() => setShowDashboard(true)}
    />
  );
}
