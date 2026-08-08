import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, expect, test, vi } from "vitest";

import { stubMatchMedia } from "../../test/stubMatchMedia";
import { useStoryPlayer } from "./useStoryPlayer";

const originalMatchMedia = window.matchMedia;

test("advances and rewinds within bounds", () => {
  const { result } = renderHook(() => useStoryPlayer({ slideCount: 3 }));

  expect(result.current.index).toBe(0);
  expect(result.current.isFirst).toBe(true);

  act(() => result.current.previous());
  expect(result.current.index).toBe(0);

  act(() => result.current.next());
  act(() => result.current.next());
  expect(result.current.index).toBe(2);
  expect(result.current.isLast).toBe(true);
});

test("calls onComplete instead of advancing past the last slide", () => {
  const onComplete = vi.fn();
  const { result } = renderHook(() =>
    useStoryPlayer({ slideCount: 2, onComplete }),
  );

  act(() => result.current.next());
  expect(result.current.index).toBe(1);
  expect(onComplete).not.toHaveBeenCalled();

  act(() => result.current.next());
  expect(result.current.index).toBe(1);
  expect(onComplete).toHaveBeenCalledOnce();
});

test("goTo clamps to the valid slide range", () => {
  const { result } = renderHook(() => useStoryPlayer({ slideCount: 4 }));

  act(() => result.current.goTo(10));
  expect(result.current.index).toBe(3);

  act(() => result.current.goTo(-5));
  expect(result.current.index).toBe(0);
});

beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
  window.matchMedia = originalMatchMedia;
});

test("auto-advances on a timer and stops on the last slide", () => {
  const onComplete = vi.fn();
  const { result } = renderHook(() =>
    useStoryPlayer({ slideCount: 2, autoAdvanceMs: 1000, onComplete }),
  );

  act(() => {
    vi.advanceTimersByTime(1000);
  });
  expect(result.current.index).toBe(1);

  act(() => {
    vi.advanceTimersByTime(5000);
  });
  expect(result.current.index).toBe(1);
  expect(onComplete).not.toHaveBeenCalled();
});

test("pausing stops the auto-advance timer", () => {
  const { result } = renderHook(() =>
    useStoryPlayer({ slideCount: 3, autoAdvanceMs: 1000 }),
  );

  act(() => result.current.pause());
  act(() => {
    vi.advanceTimersByTime(5000);
  });
  expect(result.current.index).toBe(0);

  act(() => result.current.resume());
  act(() => {
    vi.advanceTimersByTime(1000);
  });
  expect(result.current.index).toBe(1);
});

test("never auto-advances when the user prefers reduced motion", () => {
  stubMatchMedia(true);
  const { result } = renderHook(() =>
    useStoryPlayer({ slideCount: 3, autoAdvanceMs: 1000 }),
  );

  act(() => {
    vi.advanceTimersByTime(10000);
  });

  expect(result.current.index).toBe(0);
});
