import type { RouteObject } from "react-router-dom";

import { IntroPage } from "./pages/IntroPage";
import { recapLoader } from "./pages/recapLoader";
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
        lazy: () =>
          import("./pages/ProfilesPage").then((module) => ({
            Component: module.ProfilesPage,
            loader: module.profilesLoader,
          })),
      },
      {
        path: "profiles/:profileId/generating",
        ErrorBoundary: RouteErrorBoundary,
        lazy: () =>
          import("./pages/GeneratingPage").then((module) => ({
            Component: module.GeneratingPage,
            loader: module.generatingLoader,
          })),
      },
      {
        // `recapLoader` is a shared, tiny module (not part of the
        // RecapPage/SharePage chunks) so visiting a shared link only
        // downloads the share card, never the story player.
        path: "recap/:recapId",
        loader: recapLoader,
        ErrorBoundary: RouteErrorBoundary,
        lazy: () =>
          import("./pages/RecapPage").then((module) => ({
            Component: module.RecapPage,
          })),
      },
      {
        path: "recap/:recapId/share",
        loader: recapLoader,
        ErrorBoundary: RouteErrorBoundary,
        lazy: () =>
          import("./pages/SharePage").then((module) => ({
            Component: module.SharePage,
          })),
      },
    ],
  },
];
