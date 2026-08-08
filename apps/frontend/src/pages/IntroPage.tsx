import { Link } from "react-router-dom";

import { ScreenLayout } from "../ui/templates/ScreenLayout";

export function IntroPage() {
  return (
    <ScreenLayout>
      <h1>Ваши Итоги года на Авито</h1>
      <p>Узнайте, каким был ваш год на площадке.</p>
      <Link to="/profiles" className="button button--primary">
        Смотреть
      </Link>
    </ScreenLayout>
  );
}
