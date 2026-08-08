import Alert from "@mui/material/Alert";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { useRef } from "react";
import { Link as RouterLink, useParams } from "react-router-dom";

import { useGetRecap } from "../api/generated/client";
import { useShareCard } from "../features/share-export/useShareCard";
import { ShareCardPreview } from "../ui/organisms/ShareCardPreview";
import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function SharePage() {
  const { recapId } = useParams<{ recapId: string }>();
  const { data: response } = useGetRecap(recapId ?? "");
  const shareCard =
    response && response.status === 200 ? response.data.shareCard : undefined;
  const previewRef = useRef<HTMLDivElement>(null);
  const shareUrl = `${window.location.origin}/recap/${recapId}/share`;
  const { status, share, download } = useShareCard(previewRef, {
    title: shareCard?.title ?? "",
    text: shareCard?.subtitle ?? "",
    url: shareUrl,
  });

  if (!shareCard) {
    return (
      <ScreenLayout>
        <Typography>Нечем поделиться.</Typography>
      </ScreenLayout>
    );
  }

  const isExporting = status === "exporting";

  return (
    <ScreenLayout>
      <ShareCardPreview ref={previewRef} shareCard={shareCard} />
      <Stack direction="row" spacing={2}>
        <Button variant="contained" loading={isExporting} onClick={share}>
          Поделиться
        </Button>
        <Button variant="outlined" loading={isExporting} onClick={download}>
          Скачать картинку
        </Button>
      </Stack>
      {status === "copied" ? (
        <Alert severity="success">Ссылка скопирована в буфер обмена</Alert>
      ) : null}
      {status === "error" ? (
        <Alert severity="error">
          Не получилось поделиться. Скопируйте ссылку вручную: {shareUrl}
        </Alert>
      ) : null}
      <Button component={RouterLink} to={`/recap/${recapId}`}>
        Назад к итогам
      </Button>
    </ScreenLayout>
  );
}
