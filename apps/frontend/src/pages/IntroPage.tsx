import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import { Link as RouterLink } from "react-router-dom";

import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function IntroPage() {
  return (
    <ScreenLayout>
      <Typography variant="h3" component="h1">
        Ваши Итоги года на Авито
      </Typography>
      <Typography color="text.secondary">
        Узнайте, каким был ваш год на площадке.
      </Typography>
      <Button
        component={RouterLink}
        to="/profiles"
        variant="contained"
        size="large"
      >
        Смотреть
      </Button>
    </ScreenLayout>
  );
}
