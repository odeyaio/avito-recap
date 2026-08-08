import { useParams } from "react-router-dom";

import { useGetRecap } from "../api/generated/client";
import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function SharePage() {
  const { recapId } = useParams<{ recapId: string }>();
  const { data: response } = useGetRecap(recapId ?? "");
  const shareCard =
    response && response.status === 200 ? response.data.shareCard : undefined;

  if (!shareCard) {
    return (
      <ScreenLayout>
        <p>Нечем поделиться.</p>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <h1>{shareCard.title}</h1>
      <p>{shareCard.subtitle}</p>
    </ScreenLayout>
  );
}
