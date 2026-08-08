import { Link } from "react-router-dom";

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
        <p>Не удалось загрузить профили. Попробуйте ещё раз.</p>
      </ScreenLayout>
    );
  }

  return (
    <ScreenLayout>
      <h1>Выберите тестовый профиль</h1>
      <ul className="profile-list">
        {response.data.items.map((profile) => {
          const year = Math.max(...profile.availableYears);

          return (
            <li key={profile.id}>
              <Link
                to={`/profiles/${profile.id}/generating?year=${year}`}
                className="button"
              >
                {profile.displayName} · {profile.region}
                {profile.teaser ? ` — ${profile.teaser}` : ""}
              </Link>
            </li>
          );
        })}
      </ul>
    </ScreenLayout>
  );
}
