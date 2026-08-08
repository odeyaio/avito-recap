import { act, renderHook, waitFor } from "@testing-library/react";
import { afterEach, expect, test, vi } from "vitest";

import { useShareCard } from "./useShareCard";

vi.mock("html-to-image", () => ({
  toPng: vi.fn().mockResolvedValue("data:image/png;base64,zzz"),
}));

const options = { title: "Заголовок", text: "Текст", url: "https://example.test/recap/1/share" };

function refWithNode() {
  const node = document.createElement("div");
  return { current: node };
}

afterEach(() => {
  vi.unstubAllGlobals();
});

test("download exports a PNG and triggers a link click", async () => {
  const clickSpy = vi
    .spyOn(HTMLAnchorElement.prototype, "click")
    .mockImplementation(() => {});
  const { result } = renderHook(() => useShareCard(refWithNode(), options));

  await act(async () => {
    await result.current.download();
  });

  expect(clickSpy).toHaveBeenCalledOnce();
  expect(result.current.status).toBe("downloaded");
  clickSpy.mockRestore();
});

test("shares a link when Web Share API is available without file support", async () => {
  const share = vi.fn().mockResolvedValue(undefined);
  vi.stubGlobal("navigator", { ...navigator, share });

  const { result } = renderHook(() => useShareCard(refWithNode(), options));

  await act(async () => {
    await result.current.share();
  });

  expect(share).toHaveBeenCalledWith({
    title: options.title,
    text: options.text,
    url: options.url,
  });
  expect(result.current.status).toBe("shared");
});

test("shares the exported image when the browser can share files", async () => {
  const share = vi.fn().mockResolvedValue(undefined);
  const canShare = vi.fn().mockReturnValue(true);
  vi.stubGlobal("navigator", { ...navigator, share, canShare });
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      blob: () => Promise.resolve(new Blob(["x"], { type: "image/png" })),
    }),
  );

  const { result } = renderHook(() => useShareCard(refWithNode(), options));

  await act(async () => {
    await result.current.share();
  });

  expect(canShare).toHaveBeenCalledOnce();
  expect(share).toHaveBeenCalledOnce();
  const [sharedPayload] = share.mock.calls[0] as [{ files: File[] }];
  expect(sharedPayload.files[0].name).toBe("avito-recap.png");
  expect(result.current.status).toBe("shared");
});

test("falls back to copying the link when Web Share API is unavailable", async () => {
  const writeText = vi.fn().mockResolvedValue(undefined);
  vi.stubGlobal("navigator", {
    ...navigator,
    share: undefined,
    canShare: undefined,
    clipboard: { writeText },
  });

  const { result } = renderHook(() => useShareCard(refWithNode(), options));

  await act(async () => {
    await result.current.share();
  });

  expect(writeText).toHaveBeenCalledWith(options.url);
  expect(result.current.status).toBe("copied");
});

test("treats a cancelled native share sheet as idle, not an error", async () => {
  const share = vi.fn().mockRejectedValue(new DOMException("cancelled", "AbortError"));
  vi.stubGlobal("navigator", { ...navigator, share });

  const { result } = renderHook(() => useShareCard(refWithNode(), options));

  await act(async () => {
    await result.current.share();
  });

  await waitFor(() => expect(result.current.status).toBe("idle"));
});
