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

export const mockHandlers = [
  getCheckHealthMockHandler(),
  getListProfilesMockHandler({
    items: personas.map((persona) => persona.profile),
  }),
  getGenerateRecapMockHandler(({ params }) => {
    const profileId = String(params.profileId);
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
