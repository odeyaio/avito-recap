import { useNavigate } from "react-router-dom";

import { useListProfiles } from "../api/generated/client";
import { Button } from "../ui/atoms/Button";
import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function ProfilesPage() {
  const navigate = useNavigate();
  const { data: response, isLoading } = useListProfiles();

  if (isLoading) {
    return (
      <ScreenLayout>
        <p>Загружаем тестовые профили…</p>
      </ScreenLayout>
    );
  }

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
        {response.data.items.map((profile) => (
          <li key={profile.id}>
            <Button
              onClick={() => navigate(`/profiles/${profile.id}/generating`)}
            >
              {profile.displayName} · {profile.region}
              {profile.teaser ? ` — ${profile.teaser}` : ""}
            </Button>
          </li>
        ))}
      </ul>
    </ScreenLayout>
  );
}
