import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { forwardRef } from "react";

import type { ShareCard } from "../../api/generated/model";

export interface ShareCardPreviewProps {
  shareCard: ShareCard;
}

/**
 * Renders only what `ShareCard` exposes (title, subtitle, imageUrl) - never
 * the full `Recap`. The prop type is the enforcement: there is no `metrics`
 * or `achievements` field to accidentally thread through here, so this stays
 * safe to render on the public share surface even as the rest of the recap
 * evolves.
 */
export const ShareCardPreview = forwardRef<HTMLDivElement, ShareCardPreviewProps>(
  function ShareCardPreview({ shareCard }, ref) {
    return (
      <Box
        ref={ref}
        sx={{
          width: 320,
          aspectRatio: "9 / 16",
          borderRadius: 3,
          overflow: "hidden",
          position: "relative",
          display: "flex",
          flexDirection: "column",
          justifyContent: "flex-end",
          p: 3,
          color: "#fff",
          backgroundImage: shareCard.imageUrl
            ? `linear-gradient(180deg, rgba(0,0,0,0) 0%, rgba(0,0,0,0.75) 100%), url(${shareCard.imageUrl})`
            : "linear-gradient(160deg, #965eeb 0%, #4c2889 100%)",
          backgroundSize: "cover",
          backgroundPosition: "center",
        }}
      >
        <Stack spacing={1}>
          <Typography variant="overline" sx={{ opacity: 0.8 }}>
            Avito · Итоги года
          </Typography>
          <Typography variant="h5" component="p" sx={{ fontWeight: 700 }}>
            {shareCard.title}
          </Typography>
          <Typography variant="body2">{shareCard.subtitle}</Typography>
        </Stack>
      </Box>
    );
  },
);
