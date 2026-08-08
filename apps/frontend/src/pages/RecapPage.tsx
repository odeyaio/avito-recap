import Typography from "@mui/material/Typography";
import { data, useParams, type LoaderFunctionArgs } from "react-router-dom";

import { getGetRecapQueryOptions, useGetRecap } from "../api/generated/client";
import type { StoryCard } from "../api/generated/model";
import { queryClient } from "../api/client";
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

  if (!response || response.status !== 200) {
    return (
      <ScreenLayout>
        <Typography>Не удалось загрузить recap.</Typography>
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

  return <StoryPlayerLayout cards={slides} recap={response.data} />;
}
