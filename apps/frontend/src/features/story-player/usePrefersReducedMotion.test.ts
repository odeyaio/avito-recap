import { act, renderHook } from "@testing-library/react";
import { afterEach, expect, test } from "vitest";

import { stubMatchMedia } from "../../test/stubMatchMedia";
import { usePrefersReducedMotion } from "./usePrefersReducedMotion";

const originalMatchMedia = window.matchMedia;

afterEach(() => {
  window.matchMedia = originalMatchMedia;
});

test("reads the initial prefers-reduced-motion value", () => {
  stubMatchMedia(true);

  const { result } = renderHook(() => usePrefersReducedMotion());

  expect(result.current).toBe(true);
});

test("returns false when jsdom/the browser has no matchMedia support", () => {
  // @ts-expect-error simulating an environment without matchMedia
  window.matchMedia = undefined;

  const { result } = renderHook(() => usePrefersReducedMotion());

  expect(result.current).toBe(false);
});

test("reacts to the media query changing after mount", () => {
  const media = stubMatchMedia(false);
  const { result } = renderHook(() => usePrefersReducedMotion());

  expect(result.current).toBe(false);

  act(() => media.change(true));

  expect(result.current).toBe(true);
});
