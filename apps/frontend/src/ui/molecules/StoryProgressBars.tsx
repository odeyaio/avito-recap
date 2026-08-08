import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";

export interface StoryProgressBarsProps {
  count: number;
  activeIndex: number;
}

export function StoryProgressBars({ count, activeIndex }: StoryProgressBarsProps) {
  return (
    <Stack
      direction="row"
      spacing={0.5}
      role="progressbar"
      aria-label="Прогресс просмотра"
      aria-valuenow={activeIndex + 1}
      aria-valuemin={1}
      aria-valuemax={count}
      sx={{ position: "absolute", top: 8, left: 8, right: 8, zIndex: 1 }}
    >
      {Array.from({ length: count }, (_, segmentIndex) => (
        <Box
          key={segmentIndex}
          sx={{
            flex: 1,
            height: 3,
            borderRadius: 999,
            bgcolor:
              segmentIndex <= activeIndex ? "primary.main" : "action.disabledBackground",
          }}
        />
      ))}
    </Stack>
  );
}
