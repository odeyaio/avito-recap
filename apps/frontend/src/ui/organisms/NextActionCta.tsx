import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import type { NextAction } from "../../api/generated/model";

export interface NextActionCtaProps {
  nextAction: NextAction;
}

export function NextActionCta({ nextAction }: NextActionCtaProps) {
  return (
    <Card
      variant="outlined"
      sx={{ width: "100%", bgcolor: "primary.main", color: "primary.contrastText" }}
    >
      <CardContent>
        <Stack spacing={1.5} sx={{ alignItems: "flex-start" }}>
          <Typography variant="h6">{nextAction.title}</Typography>
          {nextAction.text ? <Typography>{nextAction.text}</Typography> : null}
          <Button
            component="a"
            href={nextAction.href}
            variant="contained"
            sx={{ bgcolor: "background.paper", color: "primary.main" }}
          >
            Перейти
          </Button>
        </Stack>
      </CardContent>
    </Card>
  );
}
