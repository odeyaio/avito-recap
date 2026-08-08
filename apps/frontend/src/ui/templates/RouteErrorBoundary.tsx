import Typography from "@mui/material/Typography";
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
        <Typography variant="h4" component="h1">
          {error.data.title}
        </Typography>
        <Typography color="text.secondary">
          {error.data.detail ?? "Попробуйте ещё раз чуть позже."}
        </Typography>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <Typography variant="h4" component="h1">
        Что-то пошло не так
      </Typography>
      <Typography color="text.secondary">
        Попробуйте вернуться на главную и повторить попытку.
      </Typography>
    </ScreenLayout>
  );
}
