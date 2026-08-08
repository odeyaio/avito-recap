import { Outlet, useNavigation } from "react-router-dom";

export function RootLayout() {
  const navigation = useNavigation();

  return (
    <div className="app-shell">
      {navigation.state === "loading" ? (
        <div className="app-shell__progress" role="progressbar" aria-label="Загрузка" />
      ) : null}
      <Outlet />
    </div>
  );
}
