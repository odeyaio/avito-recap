import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import { Outlet, useNavigation } from "react-router-dom";

import { GeneratingOverlay } from "../organisms/GeneratingOverlay";

const GENERATING_PATH = /^\/profiles\/[^/]+\/generating$/;

export function RootLayout() {
  const navigation = useNavigation();
  const isNavigatingToGenerating =
    navigation.state === "loading" &&
    GENERATING_PATH.test(navigation.location.pathname);

  return (
    <Box component="main" sx={{ minHeight: "100vh" }}>
      {navigation.state === "loading" && !isNavigatingToGenerating ? (
        <LinearProgress
          aria-label="Загрузка"
          sx={{ position: "fixed", top: 0, left: 0, right: 0, zIndex: 1000 }}
        />
      ) : null}
      <Outlet />
      {isNavigatingToGenerating ? <GeneratingOverlay /> : null}
    </Box>
  );
}
