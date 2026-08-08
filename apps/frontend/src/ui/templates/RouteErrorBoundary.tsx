import { isRouteErrorResponse, useRouteError } from "react-router-dom";

import type { Problem } from "../../api/generated/model";
import { ScreenLayout } from "./ScreenLayout";

function isProblem(value: unknown): value is Problem {
  return typeof value === "object" && value !== null && "code" in value;
}

export function RouteErrorBoundary() {
  const error = useRouteError();

  if (isRouteErrorResponse(error) && isProblem(error.data)) {
    return (
      <ScreenLayout>
        <h1>{error.data.title}</h1>
        <p>{error.data.detail ?? "Попробуйте ещё раз чуть позже."}</p>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <h1>Что-то пошло не так</h1>
      <p>Попробуйте вернуться на главную и повторить попытку.</p>
    </ScreenLayout>
  );
}
