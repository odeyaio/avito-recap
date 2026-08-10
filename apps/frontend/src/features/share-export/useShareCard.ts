import { toPng } from "html-to-image";
import { useCallback, useState, type RefObject } from "react";

export type ShareStatus =
  | "idle"
  | "exporting"
  | "shared"
  | "copied"
  | "downloaded"
  | "error";

export interface UseShareCardOptions {
  title: string;
  text: string;
  url: string;
}

export interface ShareCardControls {
  status: ShareStatus;
  share: () => Promise<void>;
  download: () => Promise<void>;
}

const FILE_NAME = "avito-recap.png";

async function exportPng(node: HTMLElement): Promise<string> {
  return toPng(node, { pixelRatio: 2 });
}

async function dataUrlToFile(dataUrl: string): Promise<File> {
  const blob = await (await fetch(dataUrl)).blob();
  return new File([blob], FILE_NAME, { type: "image/png" });
}

export function useShareCard(
  nodeRef: RefObject<HTMLElement | null>,
  { title, text, url }: UseShareCardOptions,
): ShareCardControls {
  const [status, setStatus] = useState<ShareStatus>("idle");

  const download = useCallback(async () => {
    if (!nodeRef.current) {
      return;
    }

    setStatus("exporting");
    try {
      const dataUrl = await exportPng(nodeRef.current);
      const link = document.createElement("a");
      link.href = dataUrl;
      link.download = FILE_NAME;
      link.click();
      setStatus("downloaded");
    } catch {
      setStatus("error");
    }
  }, [nodeRef]);

  const share = useCallback(async () => {
    setStatus("exporting");
    try {
      if (nodeRef.current && typeof navigator.canShare === "function") {
        const dataUrl = await exportPng(nodeRef.current);
        const file = await dataUrlToFile(dataUrl);

        if (navigator.canShare({ files: [file] })) {
          await navigator.share({ files: [file], title, text });
          setStatus("shared");
          return;
        }
      }

      if (typeof navigator.share === "function") {
        await navigator.share({ title, text, url });
        setStatus("shared");
        return;
      }

      await navigator.clipboard.writeText(url);
      setStatus("copied");
    } catch (error) {
      if (error instanceof DOMException && error.name === "AbortError") {
        setStatus("idle");
        return;
      }
      setStatus("error");
    }
  }, [nodeRef, title, text, url]);

  return { status, share, download };
}
