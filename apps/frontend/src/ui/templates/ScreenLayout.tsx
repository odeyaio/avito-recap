import type { ReactNode } from "react";

export interface ScreenLayoutProps {
  children: ReactNode;
}

export function ScreenLayout({ children }: ScreenLayoutProps) {
  return (
    <div className="screen-layout">
      <main className="screen-layout__content">{children}</main>
    </div>
  );
}
