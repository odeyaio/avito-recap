import { useCallback, useEffect, useState } from "react";

import { usePrefersReducedMotion } from "./usePrefersReducedMotion";

export interface UseStoryPlayerOptions {
  slideCount: number;
  autoAdvanceMs?: number;
  onComplete?: () => void;
}

export interface StoryPlayerState {
  index: number;
  isPaused: boolean;
  isFirst: boolean;
  isLast: boolean;
  next: () => void;
  previous: () => void;
  goTo: (index: number) => void;
  pause: () => void;
  resume: () => void;
}

export function useStoryPlayer({
  slideCount,
  autoAdvanceMs,
  onComplete,
}: UseStoryPlayerOptions): StoryPlayerState {
  const [index, setIndex] = useState(0);
  const [isPaused, setIsPaused] = useState(false);
  const prefersReducedMotion = usePrefersReducedMotion();

  const next = useCallback(() => {
    setIndex((current) => {
      if (current >= slideCount - 1) {
        onComplete?.();
        return current;
      }
      return current + 1;
    });
  }, [slideCount, onComplete]);

  const previous = useCallback(() => {
    setIndex((current) => Math.max(0, current - 1));
  }, []);

  const goTo = useCallback(
    (target: number) => {
      setIndex(Math.min(Math.max(target, 0), Math.max(slideCount - 1, 0)));
    },
    [slideCount],
  );

  const pause = useCallback(() => setIsPaused(true), []);
  const resume = useCallback(() => setIsPaused(false), []);

  const isLast = index >= slideCount - 1;

  useEffect(() => {
    if (!autoAdvanceMs || isPaused || prefersReducedMotion || isLast) {
      return;
    }

    const id = setTimeout(next, autoAdvanceMs);
    return () => clearTimeout(id);
  }, [autoAdvanceMs, isPaused, prefersReducedMotion, isLast, index, next]);

  return {
    index,
    isPaused,
    isFirst: index === 0,
    isLast,
    next,
    previous,
    goTo,
    pause,
    resume,
  };
}
