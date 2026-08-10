import { HttpResponse } from "msw";

import {
  getCheckHealthMockHandler,
  getGenerateRecapMockHandler,
  getGetRecapMockHandler,
  getListProfilesMockHandler,
} from "./generated/handlers";
import { personas } from "./fixtures/personas";

const recapsByProfileId = new Map(
  personas.map((persona) => [persona.profile.id, persona.recap]),
);
const recapsById = new Map(
  personas.map((persona) => [persona.recap.id, persona.recap]),
);

/**
 * Not part of the visible profile picker: reachable only by navigating
 * directly to their generating URL, so the 422/503 error states can be
 * demoed and tested without a real backend.
 */
export const INSUFFICIENT_ACTIVITY_PROFILE_ID =
  "44444444-4444-4444-8444-444444444444";
export const GENERATION_UNAVAILABLE_PROFILE_ID =
  "55555555-5555-4555-8555-555555555555";

export const mockHandlers = [
  getCheckHealthMockHandler(),
  getListProfilesMockHandler({
    items: personas.map((persona) => persona.profile),
  }),
  getGenerateRecapMockHandler(({ params }) => {
    const profileId = String(params.profileId);

    if (profileId === INSUFFICIENT_ACTIVITY_PROFILE_ID) {
      throw HttpResponse.json(
        {
          type: "https://avito-recap.example/problems/insufficient-activity",
          title: "Недостаточно активности",
          status: 422,
          code: "insufficient_activity",
          detail:
            "Для этого профиля пока накопилось слишком мало действий, чтобы собрать интересный recap.",
        },
        { status: 422 },
      );
    }

    if (profileId === GENERATION_UNAVAILABLE_PROFILE_ID) {
      throw HttpResponse.json(
        {
          type: "https://avito-recap.example/problems/generation-unavailable",
          title: "Сервис генерации недоступен",
          status: 503,
          code: "generation_unavailable",
          detail:
            "Не удалось сформировать персональное представление recap. Попробуйте ещё раз через минуту.",
        },
        { status: 503 },
      );
    }

    const recap = recapsByProfileId.get(profileId);

    if (!recap) {
      throw HttpResponse.json(
        {
          type: "https://avito-recap.example/problems/profile-not-found",
          title: "Профиль не найден",
          status: 404,
          code: "profile_not_found",
        },
        { status: 404 },
      );
    }

    return recap;
  }),
  getGetRecapMockHandler(({ params }) => {
    const recapId = String(params.recapId);
    const recap = recapsById.get(recapId);

    if (!recap) {
      throw HttpResponse.json(
        {
          type: "https://avito-recap.example/problems/recap-not-found",
          title: "Recap не найден",
          status: 404,
          code: "recap_not_found",
        },
        { status: 404 },
      );
    }

    return recap;
  }),
];
