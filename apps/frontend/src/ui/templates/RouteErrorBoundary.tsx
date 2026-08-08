import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import {
  Link as RouterLink,
  isRouteErrorResponse,
  useRevalidator,
  useRouteError,
} from "react-router-dom";

import type { Problem } from "../../api/generated/model";
import { ScreenLayout } from "./ScreenLayout";

interface StatusCopy {
  title: string;
  description: string;
  canRetry: boolean;
}

const STATUS_COPY: Record<number, StatusCopy> = {
  404: {
    title: "Профиль не найден",
    description: "Похоже, такого тестового профиля больше нет.",
    canRetry: false,
  },
  422: {
    title: "Пока недостаточно активности",
    description:
      "Для этого профиля не хватает данных, чтобы собрать интересный recap.",
    canRetry: false,
  },
  503: {
    title: "Сервис временно недоступен",
    description:
      "Не получилось собрать историю прямо сейчас. Попробуйте ещё раз через минуту.",
    canRetry: true,
  },
};

const FALLBACK_COPY: StatusCopy = {
  title: "Что-то пошло не так",
  description: "Попробуйте ещё раз чуть позже.",
  canRetry: true,
};

function isProblem(value: unknown): value is Problem {
  return typeof value === "object" && value !== null && "code" in value;
}

export function RouteErrorBoundary() {
  const error = useRouteError();
  const revalidator = useRevalidator();

  const status = isRouteErrorResponse(error) ? error.status : undefined;
  const problem =
    isRouteErrorResponse(error) && isProblem(error.data) ? error.data : undefined;
  const copy = (status !== undefined ? STATUS_COPY[status] : undefined) ?? FALLBACK_COPY;

  return (
    <ScreenLayout>
      <Alert severity="error" sx={{ width: "100%" }}>
        <AlertTitle>{problem?.title ?? copy.title}</AlertTitle>
        {problem?.detail ?? copy.description}
      </Alert>
      <Stack direction="row" spacing={2}>
        {copy.canRetry ? (
          <Button
            variant="contained"
            loading={revalidator.state === "loading"}
            onClick={() => revalidator.revalidate()}
          >
            Повторить
          </Button>
        ) : null}
        <Button component={RouterLink} to="/profiles" variant="outlined">
          Выбрать другой профиль
        </Button>
      </Stack>
    </ScreenLayout>
  );
}
