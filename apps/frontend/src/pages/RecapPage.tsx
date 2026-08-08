import { useParams } from "react-router-dom";

import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function RecapPage() {
  const { recapId } = useParams<{ recapId: string }>();

  return (
    <ScreenLayout>
      <p>Recap {recapId}.</p>
    </ScreenLayout>
  );
}
