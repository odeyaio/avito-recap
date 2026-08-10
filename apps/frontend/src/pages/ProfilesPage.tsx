import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

import { getListProfilesQueryOptions, useListProfiles } from "../api/generated/client";
import { queryClient } from "../api/client";
import { ProfileCard } from "../ui/molecules/ProfileCard";
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
            <ProfileCard
              key={profile.id}
              profile={profile}
              to={`/profiles/${profile.id}/generating?year=${year}`}
            />
          );
        })}
      </Stack>
    </ScreenLayout>
  );
}
