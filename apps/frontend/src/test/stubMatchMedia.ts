import { vi } from "vitest";

export function stubMatchMedia(initialMatches: boolean) {
  const listeners = new Set<(event: MediaQueryListEvent) => void>();
  const state = { matches: initialMatches };
  const mediaQueryList = {
    get matches() {
      return state.matches;
    },
    addEventListener: (_type: string, listener: (event: MediaQueryListEvent) => void) => {
      listeners.add(listener);
    },
    removeEventListener: (_type: string, listener: (event: MediaQueryListEvent) => void) => {
      listeners.delete(listener);
    },
  } as MediaQueryList;

  window.matchMedia = vi.fn().mockReturnValue(mediaQueryList);

  return {
    change(nextMatches: boolean) {
      state.matches = nextMatches;
      for (const listener of listeners) {
        listener({ matches: nextMatches } as MediaQueryListEvent);
      }
    },
  };
}
