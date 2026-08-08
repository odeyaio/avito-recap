import { useParams } from "react-router-dom";

import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function SharePage() {
  const { recapId } = useParams<{ recapId: string }>();

  return (
    <ScreenLayout>
      <p>Поделиться итогами {recapId}.</p>
    </ScreenLayout>
  );
}
