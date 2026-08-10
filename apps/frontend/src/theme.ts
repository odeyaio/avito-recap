import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  cssVariables: { colorSchemeSelector: "media" },
  colorSchemes: {
    light: {
      palette: {
        primary: { main: "#965eeb" },
        background: { default: "#ffffff", paper: "#f4f4f5" },
        text: { primary: "#18181b", secondary: "#6b7280" },
        divider: "#e4e4e7",
      },
    },
    dark: {
      palette: {
        primary: { main: "#b18af0" },
        background: { default: "#0b0b0d", paper: "#18181b" },
        text: { primary: "#f4f4f5", secondary: "#a1a1aa" },
        divider: "#27272a",
      },
    },
  },
  shape: { borderRadius: 8 },
  typography: {
    fontFamily:
      '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  },
  components: {
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: { textTransform: "none" },
      },
    },
  },
});
