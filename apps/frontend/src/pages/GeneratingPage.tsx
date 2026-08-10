import { data, redirect, type LoaderFunctionArgs } from "react-router-dom";

import {
  generateRecap,
  getGetRecapQueryKey,
} from "../api/generated/client";
import { queryClient } from "../api/client";
import { GeneratingOverlay } from "../ui/organisms/GeneratingOverlay";

export async function generatingLoader({ params, request }: LoaderFunctionArgs) {
  const profileId = params.profileId;

  if (!profileId) {
    throw data({ title: "Профиль не указан" }, { status: 400 });
  }

  const year = Number(new URL(request.url).searchParams.get("year"));
  const result = await generateRecap(profileId, { year });

  if (result.status !== 200) {
    throw data(result.data, { status: result.status });
  }

  queryClient.setQueryData(getGetRecapQueryKey(result.data.id), {
    data: result.data,
    status: 200,
    headers: result.headers,
  });

  return redirect(`/recap/${result.data.id}`);
}

export function GeneratingPage() {
  return <GeneratingOverlay />;
}
