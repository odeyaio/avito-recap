import Container from "@mui/material/Container";
import Stack from "@mui/material/Stack";
import type { ReactNode } from "react";

export interface ScreenLayoutProps {
  children: ReactNode;
}

export function ScreenLayout({ children }: ScreenLayoutProps) {
  return (
    <Container maxWidth="sm" sx={{ py: { xs: 4, sm: 6 } }}>
      <Stack spacing={2} sx={{ alignItems: "flex-start" }}>
        {children}
      </Stack>
    </Container>
  );
}
