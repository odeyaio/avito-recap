import Typography from "@mui/material/Typography";
import { data, useParams, type LoaderFunctionArgs } from "react-router-dom";

import { getGetRecapQueryOptions, useGetRecap } from "../api/generated/client";
import { queryClient } from "../api/client";
import { ScreenLayout } from "../ui/templates/ScreenLayout";

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

  return (
    <ScreenLayout>
      <Typography variant="h4" component="h1">
        {response.data.story.headline}
      </Typography>
      <Typography color="text.secondary">
        {response.data.story.summary}
      </Typography>
    </ScreenLayout>
  );
}
