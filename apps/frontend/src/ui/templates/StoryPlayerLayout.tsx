import Box from "@mui/material/Box";
import ButtonBase from "@mui/material/ButtonBase";
import { useEffect } from "react";

import type { StoryCard } from "../../api/generated/model";
import { useStoryPlayer } from "../../features/story-player/useStoryPlayer";
import { StoryCardRenderer } from "../organisms/StoryCardRenderer";
import { StoryProgressBars } from "../molecules/StoryProgressBars";

const AUTO_ADVANCE_MS = 6000;

export interface StoryPlayerLayoutProps {
  cards: StoryCard[];
  onComplete?: () => void;
}

export function StoryPlayerLayout({ cards, onComplete }: StoryPlayerLayoutProps) {
  const player = useStoryPlayer({
    slideCount: cards.length,
    autoAdvanceMs: AUTO_ADVANCE_MS,
    onComplete,
  });

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "ArrowRight") {
        player.next();
      } else if (event.key === "ArrowLeft") {
        player.previous();
      }
    }

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [player]);

  const currentCard = cards[player.index];

  return (
    <Box
      sx={{
        position: "relative",
        minHeight: "100vh",
        overflow: "hidden",
        bgcolor: "background.default",
      }}
      onPointerDown={player.pause}
      onPointerUp={player.resume}
      onPointerCancel={player.resume}
    >
      <StoryProgressBars count={cards.length} activeIndex={player.index} />
      <Box sx={{ position: "absolute", inset: 0, display: "flex" }}>
        <ButtonBase
          aria-label="Предыдущая карточка"
          onClick={player.previous}
          disabled={player.isFirst}
          sx={{ flex: 1, height: "100%" }}
        />
        <ButtonBase
          aria-label="Следующая карточка"
          onClick={player.next}
          sx={{ flex: 2, height: "100%" }}
        />
      </Box>
      <Box
        sx={{
          position: "relative",
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          p: 3,
          pointerEvents: "none",
        }}
      >
        {currentCard ? <StoryCardRenderer card={currentCard} /> : null}
      </Box>
    </Box>
  );
}
