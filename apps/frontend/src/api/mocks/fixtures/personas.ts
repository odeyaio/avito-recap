import type { ProfileSummary, Recap } from "../../generated/model";

/**
 * Hand-written (not orval-generated) fixtures used by the mock MSW worker in
 * development and by tests. Each persona is tuned to plausibly satisfy a
 * distinct rule from catalog/behaviors.yaml and a distinct set of
 * catalog/achievements.yaml codes, so the story/achievement UI has real
 * variety to render against before the backend is available.
 */
export interface Persona {
  profile: ProfileSummary;
  recap: Recap;
}

const icon = (key: string) =>
  `https://cdn.avito-recap.example/achievement-icons/${key}/v1.webp`;

const explorer: Persona = {
  profile: {
    id: "11111111-1111-4111-8111-111111111111",
    displayName: "Алексей",
    region: "Омск",
    registeredAt: "2019-03-12T00:00:00Z",
    availableYears: [2025],
    teaser: "132 объявления, 9 категорий — но так и не выбрал одну",
  },
  recap: {
    id: "21111111-1111-4111-8111-111111111111",
    profile: {
      id: "11111111-1111-4111-8111-111111111111",
      displayName: "Алексей",
      region: "Омск",
      registeredAt: "2019-03-12T00:00:00Z",
      availableYears: [2025],
    },
    period: { year: 2025, startsAt: "2025-01-01", endsAt: "2026-01-01" },
    behavior: {
      primary: {
        code: "explorer",
        name: "Исследователь",
        description: "Много уникальных просмотров и категорий, мало контактов",
        explanation:
          "Вы посмотрели 132 уникальных объявления в 9 разных категориях и почти не выходили на контакт с продавцами — вам было интереснее изучать, чем покупать.",
        score: 0.82,
        evidence: [
          {
            metricCode: "activity.unique_listings_viewed",
            description: "Вы посмотрели 132 уникальных объявления",
          },
          {
            metricCode: "interests.distinct_categories",
            description: "Вы заглянули в 9 разных категорий",
          },
        ],
      },
      traits: [
        {
          code: "new_arrivals_hunter",
          name: "Охотник за новинками",
          description:
            "Часто открывает свежие объявления и переходит из уведомлений",
          explanation:
            "34 раза вы открывали объявления прямо из уведомлений — свежее почти никогда от вас не ускользало.",
          score: 0.58,
          evidence: [
            {
              metricCode: "features.notification_opens",
              description: "34 перехода из уведомлений",
            },
          ],
        },
      ],
    },
    metrics: {
      activity: {
        views: 1400,
        uniqueListingsViewed: 132,
        searches: 210,
        activeDays: 240,
        activeMonths: 11,
        longestActiveStreakDays: 18,
        mostActiveMonth: "2025-09",
        favoriteHour: 21,
      },
      interests: {
        topCategories: [
          {
            code: "electronics",
            name: "Электроника",
            actions: 320,
            activeMonths: 9,
            share: 0.28,
          },
          {
            code: "home",
            name: "Дом и дача",
            actions: 210,
            activeMonths: 7,
            share: 0.18,
          },
          {
            code: "hobby",
            name: "Хобби и отдых",
            actions: 180,
            activeMonths: 6,
            share: 0.16,
          },
        ],
        newCategories: [
          { code: "auto_parts", name: "Автозапчасти", actions: 24, activeMonths: 2 },
        ],
      },
      intent: {
        repeatViews: 46,
        favoritesAdded: 18,
        activeFavorites: 9,
        contacts: 7,
        completedDeals: 1,
        contactToDealConversion: 0.14,
      },
      marketplace: {
        purchases: 1,
        sales: 0,
        deliveryDeals: 1,
        publishedListings: 0,
        closedListings: 0,
        listingViews: 0,
        listingContacts: 0,
      },
      community: {
        reviewsLeft: 1,
        reviewsReceived: 0,
        fiveStarRatings: 0,
      },
      features: {
        notificationOpens: 34,
        promotionUses: 0,
      },
    },
    achievements: [
      {
        code: "explorer_of_the_year",
        name: "Исследователь года",
        description: "Вы посмотрели множество разных объявлений",
        explanation: "132 уникальных объявления за год — заметно больше, чем у большинства.",
        image: { url: icon("explorer-of-the-year"), alt: "Исследователь года" },
        shareable: true,
      },
      {
        code: "always_informed",
        name: "Всегда в курсе",
        description: "Вы регулярно открывали объявления из уведомлений",
        explanation: "34 перехода из уведомлений — вы одними из первых узнавали о новом.",
        image: { url: icon("always-informed"), alt: "Всегда в курсе" },
        shareable: true,
      },
    ],
    story: {
      headline: "Ваш тип года — «Исследователь»",
      summary:
        "Вы посмотрели 132 объявления в 9 категориях и почти не связывались с продавцами — этот год был про изучение, а не про сделки.",
      cards: [
        {
          id: "card-top-category",
          kind: "top_category",
          title: "Электроника — ваша главная тема",
          text: "320 действий и 9 активных месяцев подряд в этой категории.",
          metricCodes: ["interests.top_categories"],
          shareable: true,
        },
        {
          id: "card-activity-streak",
          kind: "activity_streak",
          title: "18 дней подряд",
          text: "Ваша самая длинная серия активности за год — почти три недели без перерыва.",
          metricCodes: ["activity.longest_active_streak_days"],
          shareable: true,
        },
        {
          id: "card-behavior-reveal",
          kind: "behavior_reveal",
          title: "Исследователь",
          text: "Много уникальных просмотров и категорий, мало контактов.",
          shareable: false,
        },
        {
          id: "card-achievements",
          kind: "achievement_spotlight",
          title: "2 ачивки открыты",
          text: "«Исследователь года» и «Всегда в курсе».",
          shareable: true,
        },
      ],
    },
    nextAction: {
      code: "open_frequent_interest_feed",
      title: "Посмотреть предложения по частым интересам",
      text: "Продолжите изучать электронику — с сентября вы возвращались к ней чаще всего.",
      href: "/search?category=electronics",
      target: { type: "search", categoryCode: "electronics" },
    },
    shareCard: {
      title: "Алексей — Исследователь года",
      subtitle: "132 объявления, 9 категорий, 11 активных месяцев",
    },
    generatedAt: "2026-01-05T09:00:00Z",
  },
};

const determinedBuyer: Persona = {
  profile: {
    id: "11111111-1111-4111-8111-222222222222",
    displayName: "Марина",
    region: "Санкт-Петербург",
    registeredAt: "2021-06-01T00:00:00Z",
    availableYears: [2025],
    teaser: "Знала, чего хочет, и не отступала от цели",
  },
  recap: {
    id: "21111111-1111-4111-8111-222222222222",
    profile: {
      id: "11111111-1111-4111-8111-222222222222",
      displayName: "Марина",
      region: "Санкт-Петербург",
      registeredAt: "2021-06-01T00:00:00Z",
      availableYears: [2025],
    },
    period: { year: 2025, startsAt: "2025-01-01", endsAt: "2026-01-01" },
    behavior: {
      primary: {
        code: "determined_buyer",
        name: "Целеустремлённый покупатель",
        description: "Узкий интерес, повторные просмотры, избранное и контакты",
        explanation:
          "68% ваших действий пришлось на одну категорию, вы возвращались к объявлениям снова и снова и в итоге вышли на контакт с 9 продавцами.",
        score: 0.9,
        evidence: [
          {
            metricCode: "interests.top_category_share",
            description: "68% действий — в одной категории",
          },
          {
            metricCode: "intent.repeat_viewed_listings",
            description: "14 объявлений вы просматривали повторно",
          },
        ],
      },
      traits: [],
    },
    metrics: {
      activity: {
        views: 640,
        uniqueListingsViewed: 58,
        searches: 95,
        activeDays: 130,
        activeMonths: 8,
        longestActiveStreakDays: 9,
        mostActiveMonth: "2025-04",
        favoriteHour: 20,
      },
      interests: {
        topCategories: [
          {
            code: "furniture",
            name: "Мебель",
            actions: 185,
            activeMonths: 6,
            share: 0.68,
          },
        ],
        newCategories: [],
        mostConsistentCategory: {
          code: "furniture",
          name: "Мебель",
          actions: 185,
          activeMonths: 6,
          share: 0.68,
        },
      },
      intent: {
        repeatViews: 14,
        favoritesAdded: 34,
        activeFavorites: 21,
        contacts: 9,
        completedDeals: 4,
        contactToDealConversion: 0.44,
      },
      marketplace: {
        purchases: 2,
        sales: 0,
        deliveryDeals: 3,
        publishedListings: 0,
        closedListings: 0,
        listingViews: 0,
        listingContacts: 0,
      },
      community: {
        reviewsLeft: 2,
        reviewsReceived: 0,
        fiveStarRatings: 0,
      },
      features: {
        notificationOpens: 6,
        promotionUses: 0,
      },
    },
    achievements: [
      {
        code: "know_exactly_what_i_want",
        name: "Точно знаю, чего хочу",
        description: "Вы проводили большую часть времени в одной категории",
        explanation: "68% действий и 185 обращений — «Мебель» весь год оставалась вашим фокусом.",
        image: { url: icon("know-exactly-what-i-want"), alt: "Точно знаю, чего хочу" },
        shareable: true,
      },
      {
        code: "collection_of_the_year",
        name: "Коллекция года",
        description: "Вы собрали большую коллекцию объявлений в избранном",
        explanation: "34 объявления в избранном — внушительная подборка вариантов на выбор.",
        image: { url: icon("collection-of-the-year"), alt: "Коллекция года" },
        shareable: true,
      },
      {
        code: "delivery_is_easier",
        name: "С доставкой удобнее",
        description: "Вы выбирали доставку для значительной части сделок",
        explanation: "3 из 4 сделок прошли с доставкой — вы явно оценили удобство.",
        image: { url: icon("delivery-is-easier"), alt: "С доставкой удобнее" },
        shareable: true,
      },
    ],
    story: {
      headline: "Ваш тип года — «Целеустремлённый покупатель»",
      summary:
        "68% ваших действий — в мебели. Вы возвращались к избранным объявлениям и уверенно доводили дело до сделки.",
      cards: [
        {
          id: "card-top-category",
          kind: "top_category",
          title: "Мебель — 68% всех действий",
          text: "Ни одна другая категория не была вам так интересна.",
          metricCodes: ["interests.top_categories"],
          shareable: true,
        },
        {
          id: "card-favorites",
          kind: "activity_streak",
          title: "34 объявления в избранном",
          text: "Вы собрали внушительную коллекцию вариантов на выбор.",
          metricCodes: ["intent.favorites_added"],
          shareable: true,
        },
        {
          id: "card-behavior-reveal",
          kind: "behavior_reveal",
          title: "Целеустремлённый покупатель",
          text: "Узкий интерес, повторные просмотры, избранное и контакты.",
          shareable: false,
        },
        {
          id: "card-achievements",
          kind: "achievement_spotlight",
          title: "3 ачивки открыты",
          text: "«Точно знаю, чего хочу», «Коллекция года», «С доставкой удобнее».",
          shareable: true,
        },
      ],
    },
    nextAction: {
      code: "return_to_current_options",
      title: "Вернуться к актуальным вариантам",
      text: "В избранном ещё есть подходящая мебель — часть объявлений всё ещё активна.",
      href: "/favorites?category=furniture",
      target: { type: "favorites", categoryCode: "furniture" },
    },
    shareCard: {
      title: "Марина — Целеустремлённый покупатель",
      subtitle: "68% действий в одной категории, 4 сделки с доставкой",
    },
    generatedAt: "2026-01-05T09:00:00Z",
  },
};

const effectiveSeller: Persona = {
  profile: {
    id: "11111111-1111-4111-8111-333333333333",
    displayName: "Игорь",
    region: "Новосибирск",
    registeredAt: "2017-11-20T00:00:00Z",
    availableYears: [2025],
    teaser: "И продавал, и покупал — почти без осечек",
  },
  recap: {
    id: "21111111-1111-4111-8111-333333333333",
    profile: {
      id: "11111111-1111-4111-8111-333333333333",
      displayName: "Игорь",
      region: "Новосибирск",
      registeredAt: "2017-11-20T00:00:00Z",
      availableYears: [2025],
    },
    period: { year: 2025, startsAt: "2025-01-01", endsAt: "2026-01-01" },
    behavior: {
      primary: {
        code: "effective_seller",
        name: "Эффективный продавец",
        description: "Высокая доля объявлений завершается продажей",
        explanation:
          "6 из 9 объявлений закрылись продажей, а средняя оценка 4.8 — покупатели вам явно доверяли.",
        score: 0.87,
        evidence: [
          {
            metricCode: "marketplace.sales",
            description: "6 завершённых продаж",
          },
          {
            metricCode: "community.average_rating",
            description: "Средняя оценка 4.8 из 5",
          },
        ],
      },
      traits: [
        {
          code: "universal_user",
          name: "Универсальный пользователь",
          description: "Активен и как покупатель, и как продавец",
          explanation: "Вы не только продавали, но и купили 2 вещи — активны с обеих сторон.",
          score: 0.51,
          evidence: [
            {
              metricCode: "marketplace.purchases",
              description: "2 покупки за год",
            },
          ],
        },
      ],
    },
    metrics: {
      activity: {
        views: 2100,
        uniqueListingsViewed: 340,
        searches: 180,
        activeDays: 260,
        activeMonths: 12,
        longestActiveStreakDays: 24,
        mostActiveMonth: "2025-11",
        favoriteHour: 19,
      },
      interests: {
        topCategories: [
          {
            code: "tools",
            name: "Инструменты",
            actions: 410,
            activeMonths: 10,
            share: 0.35,
          },
        ],
        newCategories: [],
      },
      intent: {
        repeatViews: 12,
        favoritesAdded: 8,
        activeFavorites: 3,
        contacts: 63,
        completedDeals: 8,
        contactToDealConversion: 0.13,
      },
      marketplace: {
        purchases: 2,
        sales: 6,
        deliveryDeals: 2,
        publishedListings: 9,
        closedListings: 6,
        listingViews: 1450,
        listingContacts: 63,
      },
      community: {
        reviewsLeft: 5,
        reviewsReceived: 7,
        fiveStarRatings: 6,
        averageRating: 4.8,
      },
      features: {
        notificationOpens: 15,
        promotionUses: 3,
      },
    },
    achievements: [
      {
        code: "successful_seller",
        name: "Успешный продавец",
        description: "Вы успешно завершили несколько продаж",
        explanation: "6 завершённых продаж за год — уверенный результат.",
        image: { url: icon("successful-seller"), alt: "Успешный продавец" },
        shareable: true,
      },
      {
        code: "both_sides_of_deal",
        name: "Две стороны сделки",
        description: "Вы побывали и покупателем, и продавцом",
        explanation: "2 покупки и 6 продаж — вы знаете площадку с обеих сторон.",
        image: { url: icon("both-sides-of-deal"), alt: "Две стороны сделки" },
        shareable: true,
      },
      {
        code: "five_stars",
        name: "Пять звёзд",
        description: "Вы получили несколько оценок в пять звёзд",
        explanation: "6 оценок на пять звёзд — покупатели остались довольны.",
        image: { url: icon("five-stars"), alt: "Пять звёзд" },
        shareable: true,
      },
      {
        code: "reliable_seller",
        name: "Надёжный продавец",
        description: "Вы успешно продавали и получали положительные оценки",
        explanation: "Средняя оценка 4.8 при 6 продажах — стабильно надёжный продавец.",
        image: { url: icon("reliable-seller"), alt: "Надёжный продавец" },
        shareable: true,
      },
    ],
    story: {
      headline: "Ваш тип года — «Эффективный продавец»",
      summary:
        "6 из 9 объявлений закрылись продажей, а средняя оценка 4.8 — этот год вы провели по обе стороны сделки.",
      cards: [
        {
          id: "card-top-category",
          kind: "top_category",
          title: "Инструменты — ваша сильная сторона",
          text: "410 действий и 10 активных месяцев в этой категории.",
          metricCodes: ["interests.top_categories"],
          shareable: true,
        },
        {
          id: "card-sales",
          kind: "activity_streak",
          title: "6 продаж, 63 обращения",
          text: "Ваши объявления привлекали внимание и доводили дело до сделки.",
          metricCodes: ["marketplace.sales", "marketplace.listingContacts"],
          shareable: true,
        },
        {
          id: "card-behavior-reveal",
          kind: "behavior_reveal",
          title: "Эффективный продавец",
          text: "Высокая доля объявлений завершается продажей.",
          shareable: false,
        },
        {
          id: "card-achievements",
          kind: "achievement_spotlight",
          title: "4 ачивки открыты",
          text: "«Успешный продавец», «Две стороны сделки», «Пять звёзд», «Надёжный продавец».",
          shareable: true,
        },
      ],
    },
    nextAction: {
      code: "open_seller_statistics",
      title: "Посмотреть статистику и продать ещё",
      text: "У вас отличная конверсия в продажу — самое время выставить новый лот.",
      href: "/seller/statistics",
      target: { type: "seller_statistics" },
    },
    shareCard: {
      title: "Игорь — Эффективный продавец",
      subtitle: "6 продаж, оценка 4.8, обе стороны сделки",
    },
    generatedAt: "2026-01-05T09:00:00Z",
  },
};

export const personas: Persona[] = [explorer, determinedBuyer, effectiveSeller];
