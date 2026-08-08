import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import { Outlet, useNavigation } from "react-router-dom";

export function RootLayout() {
  const navigation = useNavigation();

  return (
    <Box sx={{ minHeight: "100vh" }}>
      {navigation.state === "loading" ? (
        <LinearProgress
          aria-label="Загрузка"
          sx={{ position: "fixed", top: 0, left: 0, right: 0, zIndex: 1000 }}
        />
      ) : null}
      <Outlet />
    </Box>
  );
}
