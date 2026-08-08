import { Route, Routes } from "react-router-dom";

import { GeneratingPage } from "./pages/GeneratingPage";
import { IntroPage } from "./pages/IntroPage";
import { ProfilesPage } from "./pages/ProfilesPage";
import { RecapPage } from "./pages/RecapPage";
import { SharePage } from "./pages/SharePage";

export function App() {
  return (
    <Routes>
      <Route path="/" element={<IntroPage />} />
      <Route path="/profiles" element={<ProfilesPage />} />
      <Route
        path="/profiles/:profileId/generating"
        element={<GeneratingPage />}
      />
      <Route path="/recap/:recapId" element={<RecapPage />} />
      <Route path="/recap/:recapId/share" element={<SharePage />} />
    </Routes>
  );
}
