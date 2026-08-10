import { useEffect, useState } from "react";

export function useCyclingMessages(
  messages: readonly string[],
  intervalMs: number,
): string {
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const id = setInterval(() => {
      setIndex((current) => (current + 1) % messages.length);
    }, intervalMs);

    return () => clearInterval(id);
  }, [messages, intervalMs]);

  return messages[index];
}
