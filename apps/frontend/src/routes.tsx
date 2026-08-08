import type { RouteObject } from "react-router-dom";

import { GeneratingPage, generatingLoader } from "./pages/GeneratingPage";
import { IntroPage } from "./pages/IntroPage";
import { ProfilesPage, profilesLoader } from "./pages/ProfilesPage";
import { RecapPage, recapLoader } from "./pages/RecapPage";
import { SharePage } from "./pages/SharePage";
import { RootLayout } from "./ui/templates/RootLayout";
import { RouteErrorBoundary } from "./ui/templates/RouteErrorBoundary";

export const routes: RouteObject[] = [
  {
    path: "/",
    Component: RootLayout,
    ErrorBoundary: RouteErrorBoundary,
    children: [
      { index: true, Component: IntroPage },
      {
        path: "profiles",
        loader: profilesLoader,
        Component: ProfilesPage,
      },
      {
        path: "profiles/:profileId/generating",
        loader: generatingLoader,
        Component: GeneratingPage,
        ErrorBoundary: RouteErrorBoundary,
      },
      {
        path: "recap/:recapId",
        loader: recapLoader,
        Component: RecapPage,
        ErrorBoundary: RouteErrorBoundary,
      },
      {
        path: "recap/:recapId/share",
        loader: recapLoader,
        Component: SharePage,
        ErrorBoundary: RouteErrorBoundary,
      },
    ],
  },
];
