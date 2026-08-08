import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import { useCyclingMessages } from "../../features/generation-progress/useCyclingMessages";

const MESSAGES = [
  "Считаем ваши действия за год…",
  "Ищем ваш тип поведения…",
  "Собираем ачивки…",
  "Пишем историю…",
] as const;

export function GeneratingOverlay() {
  const message = useCyclingMessages(MESSAGES, 1100);

  return (
    <Box
      role="status"
      aria-live="polite"
      sx={{
        position: "fixed",
        inset: 0,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        bgcolor: "background.default",
        zIndex: 1200,
      }}
    >
      <Stack spacing={2} sx={{ alignItems: "center" }}>
        <CircularProgress size={48} />
        <Typography variant="h6">{message}</Typography>
      </Stack>
    </Box>
  );
}
