package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"avito-recap/internal/catalog"
	"avito-recap/internal/engine"
	"avito-recap/internal/model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type pingerFunc func(context.Context) error

func (f pingerFunc) Ping(ctx context.Context) error {
	return f(ctx)
}

type profileServiceFunc func(context.Context) ([]model.ProfileSummary, error)

func (f profileServiceFunc) ListProfiles(ctx context.Context) ([]model.ProfileSummary, error) {
	return f(ctx)
}

type recapServiceStub struct {
	generate func(context.Context, uuid.UUID, int) (model.Recap, error)
	get      func(context.Context, uuid.UUID) (model.Recap, error)
}

func (s recapServiceStub) GenerateRecap(ctx context.Context, profileID uuid.UUID, year int) (model.Recap, error) {
	return s.generate(ctx, profileID, year)
}

func (s recapServiceStub) GetRecap(ctx context.Context, recapID uuid.UUID) (model.Recap, error) {
	return s.get(ctx, recapID)
}

func TestRegisteredHandlers(t *testing.T) {
	t.Parallel()

	router := echo.New()
	recap := recapFixture(t)
	RegisterHandlers(
		router,
		pingerFunc(func(context.Context) error { return nil }),
		profileServiceFunc(func(context.Context) ([]model.ProfileSummary, error) {
			return []model.ProfileSummary{{
				User: model.User{
					ID: recap.Profile.User.ID, DisplayName: recap.Profile.User.DisplayName,
					Region: recap.Profile.User.Region, RegisteredAt: recap.Profile.User.RegisteredAt,
				},
				AvailableYears: recap.Profile.AvailableYears,
			}}, nil
		}),
		recapServiceStub{
			generate: func(context.Context, uuid.UUID, int) (model.Recap, error) { return recap, nil },
			get:      func(context.Context, uuid.UUID) (model.Recap, error) { return recap, nil },
		},
	)

	tests := []struct {
		name        string
		method      string
		target      string
		body        string
		wantStatus  int
		contentType string
	}{
		{
			name:        "health",
			method:      http.MethodGet,
			target:      "/api/v1/health",
			wantStatus:  http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "profiles",
			method:      http.MethodGet,
			target:      "/api/v1/profiles",
			wantStatus:  http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "generate recap",
			method:      http.MethodPost,
			target:      "/api/v1/profiles/7a9e06c2-42a2-4cb9-ae7d-3f187a51735c/recaps",
			body:        `{"year":2025}`,
			wantStatus:  http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "get recap",
			method:      http.MethodGet,
			target:      "/api/v1/recaps/410942ba-9544-49db-99fb-a02cc19e5b84",
			wantStatus:  http.StatusOK,
			contentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequestWithContext(t.Context(), tt.method, tt.target, strings.NewReader(tt.body))
			if tt.body != "" {
				request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.wantStatus, response.Body.String())
			}
			if got := response.Header().Get(echo.HeaderContentType); got != tt.contentType {
				t.Fatalf("Content-Type = %q, want %q", got, tt.contentType)
			}
		})
	}
}

func TestCheckHealthDatabaseUnavailable(t *testing.T) {
	t.Parallel()

	router := echo.New()
	RegisterHandlers(
		router,
		pingerFunc(func(context.Context) error { return errors.New("database is down") }),
		profileServiceFunc(func(context.Context) ([]model.ProfileSummary, error) { return nil, nil }),
		recapServiceStub{
			generate: func(context.Context, uuid.UUID, int) (model.Recap, error) { return model.Recap{}, nil },
			get:      func(context.Context, uuid.UUID) (model.Recap, error) { return model.Recap{}, nil },
		},
	)

	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusServiceUnavailable, response.Body.String())
	}
}

func recapFixture(t *testing.T) model.Recap {
	t.Helper()

	now := time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC)
	recapID := uuid.New()
	profileID := uuid.New()
	defaultAction, err := json.Marshal(catalog.DefaultAction{
		Code: "open_personal_collection", Title: "Открыть подборку",
		TargetType: "search", Href: "https://www.avito.ru/rossiya",
	})
	if err != nil {
		t.Fatal(err)
	}
	evidence, err := json.Marshal([]engine.Evidence{{
		MetricCode: "activity.total_actions", Actual: int64(3), Matched: true,
	}})
	if err != nil {
		t.Fatal(err)
	}

	return model.Recap{
		Snapshot: model.RecapSnapshot{
			ID: recapID, UserID: profileID,
			PeriodStart: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:   time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			Metrics:     json.RawMessage(`{"activity":{"views":3,"active_days":1,"active_months":1},"intent":{},"marketplace":{},"community":{},"features":{}}`),
			GeneratedAt: now,
		},
		Profile: model.ProfileSummary{
			User: model.User{
				ID: profileID, DisplayName: "Тестовый профиль", Region: "Омск",
				RegisteredAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			AvailableYears: []int32{2024, 2025},
			LatestRecapID:  &recapID,
		},
		Behaviors: []model.StoredBehavior{{
			Match: model.RecapBehaviorType{
				RecapID: recapID, IsPrimary: true, Position: 1, Evidence: evidence,
			},
			Definition: model.BehaviorTypeDefinition{
				Code: "rare_user", Name: "Редкий пользователь",
				Description: "Небольшое количество действий за год", DefaultAction: defaultAction,
			},
		}},
		Achievements: []model.StoredAchievement{},
		NextAction: &model.RecapNextAction{
			RecapID: recapID, Code: "open_personal_collection",
			Href: "https://www.avito.ru/rossiya", Target: json.RawMessage(`{"type":"search"}`),
		},
	}
}
