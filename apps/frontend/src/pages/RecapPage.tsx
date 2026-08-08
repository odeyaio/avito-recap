import Typography from "@mui/material/Typography";
import { useState } from "react";
import { data, useParams, type LoaderFunctionArgs } from "react-router-dom";

import { getGetRecapQueryOptions, useGetRecap } from "../api/generated/client";
import type { StoryCard } from "../api/generated/model";
import { queryClient } from "../api/client";
import { RecapDashboard } from "../ui/organisms/RecapDashboard";
import { ScreenLayout } from "../ui/templates/ScreenLayout";
import { StoryPlayerLayout } from "../ui/templates/StoryPlayerLayout";

export function recapLoader({ params }: LoaderFunctionArgs) {
  const recapId = params.recapId;

  if (!recapId) {
    throw data({ title: "Recap не указан" }, { status: 400 });
  }

  return queryClient.ensureQueryData(getGetRecapQueryOptions(recapId));
}

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
