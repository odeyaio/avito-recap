import { useNavigate } from "react-router-dom";

import { Button } from "../ui/atoms/Button";
import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function IntroPage() {
  const navigate = useNavigate();

  return (
    <ScreenLayout>
      <h1>Ваши Итоги года на Авито</h1>
      <p>Узнайте, каким был ваш год на площадке.</p>
      <Button onClick={() => navigate("/profiles")}>Смотреть</Button>
    </ScreenLayout>
  );
}
