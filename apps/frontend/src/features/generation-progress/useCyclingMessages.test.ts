import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, expect, test, vi } from "vitest";

import { useCyclingMessages } from "./useCyclingMessages";

beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

test("cycles through messages on an interval and wraps around", () => {
  const messages = ["one", "two", "three"];
  const { result } = renderHook(() => useCyclingMessages(messages, 1000));

  expect(result.current).toBe("one");

  act(() => {
    vi.advanceTimersByTime(1000);
  });
  expect(result.current).toBe("two");

  act(() => {
    vi.advanceTimersByTime(2000);
  });
  expect(result.current).toBe("one");
});
