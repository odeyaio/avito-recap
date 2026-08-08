import { useParams } from "react-router-dom";

import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function GeneratingPage() {
  const { profileId } = useParams<{ profileId: string }>();

  return (
    <ScreenLayout>
      <p>Собираем итоги года для профиля {profileId}…</p>
    </ScreenLayout>
  );
}
