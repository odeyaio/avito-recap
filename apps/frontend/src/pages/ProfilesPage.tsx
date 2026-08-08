import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { Link as RouterLink } from "react-router-dom";

import { getListProfilesQueryOptions, useListProfiles } from "../api/generated/client";
import { queryClient } from "../api/client";
import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function profilesLoader() {
  return queryClient.ensureQueryData(getListProfilesQueryOptions());
}

export function ProfilesPage() {
  const { data: response } = useListProfiles();

  if (!response || response.status !== 200) {
    return (
      <ScreenLayout>
        <Typography>Не удалось загрузить профили. Попробуйте ещё раз.</Typography>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <Typography variant="h4" component="h1">
        Выберите тестовый профиль
      </Typography>
      <Stack spacing={2} sx={{ width: "100%" }}>
        {response.data.items.map((profile) => {
          const year = Math.max(...profile.availableYears);

          return (
            <Card key={profile.id} variant="outlined">
              <CardActionArea
                component={RouterLink}
                to={`/profiles/${profile.id}/generating?year=${year}`}
              >
                <CardContent>
                  <Typography variant="subtitle1">
                    {profile.displayName} · {profile.region}
                  </Typography>
                  {profile.teaser ? (
                    <Typography variant="body2" color="text.secondary">
                      {profile.teaser}
                    </Typography>
                  ) : null}
                </CardContent>
              </CardActionArea>
            </Card>
          );
        })}
      </Stack>
    </ScreenLayout>
  );
}
