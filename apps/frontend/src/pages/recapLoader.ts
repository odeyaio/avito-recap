import { data, type LoaderFunctionArgs } from "react-router-dom";

import { getGetRecapQueryOptions } from "../api/generated/client";
import { queryClient } from "../api/client";

/**
 * Its own module (not co-located with RecapPage) so the share route can
 * import just this and skip pulling in the story player - someone opening
 * a shared link only needs the share card, not the whole recap experience.
 */
export function recapLoader({ params }: LoaderFunctionArgs) {
  const recapId = params.recapId;

  if (!recapId) {
    throw data({ title: "Recap не указан" }, { status: 400 });
  }

  return queryClient.ensureQueryData(getGetRecapQueryOptions(recapId));
}
