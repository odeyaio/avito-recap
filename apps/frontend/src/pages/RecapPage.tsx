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
        <p>Не удалось загрузить recap.</p>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <h1>{response.data.story.headline}</h1>
      <p>{response.data.story.summary}</p>
    </ScreenLayout>
  );
}
