import Typography from "@mui/material/Typography";
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
        <Typography>Нечем поделиться.</Typography>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <Typography variant="h4" component="h1">
        {shareCard.title}
      </Typography>
      <Typography color="text.secondary">{shareCard.subtitle}</Typography>
    </ScreenLayout>
  );
}
