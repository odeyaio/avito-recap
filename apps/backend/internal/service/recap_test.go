package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"avito-recap/internal/catalog"
	"avito-recap/internal/engine"
	"avito-recap/internal/model"
	"avito-recap/internal/repository"

	"github.com/google/uuid"
)

func TestGenerateRecapBuildsAndPersistsDraft(t *testing.T) {
	t.Parallel()

	profileID := uuid.New()
	behaviorID := uuid.New()
	achievementID := uuid.New()
	now := time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC)
	user := model.User{ID: profileID, Timezone: "UTC", IsTestProfile: true}
	profile := model.ProfileSummary{User: user, AvailableYears: []int32{2024, 2025}}
	behavior := catalog.BehaviorDefinition{
		Code: "rare_user",
		DefaultAction: catalog.DefaultAction{
			Code:       "open_personal_collection",
			Title:      "Открыть подборку",
			TargetType: "search",
			Href:       "https://www.avito.ru/rossiya",
		},
	}
	achievement := catalog.AchievementDefinition{Code: "explorer", ShareableByDefault: true}
	metrics := engine.NewMetrics()
	metrics.Set("activity.total_actions", int64(3))

	repositories := &fakeRepositories{
		profile: profile,
		dataset: model.VersionedDataset{
			Value:   model.Dataset{User: user, DataCutoffAt: now},
			Version: "sha256:dataset",
		},
		catalogs: catalog.Set{
			Behaviors:      catalog.BehaviorCatalog{Version: "v3"},
			Achievements:   catalog.AchievementCatalog{Version: "v3"},
			BehaviorIDs:    map[string]uuid.UUID{"rare_user": behaviorID},
			AchievementIDs: map[string]uuid.UUID{"explorer": achievementID},
		},
	}
	ruleEngine := ruleEngineFunc(func(
		engine.Dataset,
		catalog.AchievementCatalog,
		catalog.BehaviorCatalog,
	) (engine.Result, error) {
		return engine.Result{
			Metrics: metrics,
			Behaviors: []engine.BehaviorMatch{{
				Definition: behavior, IsPrimary: true, Position: 1,
				Evidence: []engine.Evidence{{MetricCode: "activity.total_actions", Matched: true}},
			}},
			Achievements: []engine.AchievementMatch{{
				Definition: achievement, Position: 1,
			}},
		}, nil
	})

	recapService := NewRecapService(
		repositories,
		repositories,
		repositories,
		repositories,
		ruleEngine,
		DefaultActionResolver{},
		nil,
	)
	recapService.clock = fixedClock{value: now}

	recap, err := recapService.GenerateRecap(t.Context(), profileID, 2025)
	if err != nil {
		t.Fatalf("GenerateRecap() error = %v", err)
	}
	if repositories.saved == nil {
		t.Fatal("SaveRecap() was not called")
	}
	if recap.Snapshot.DatasetVersion != "sha256:dataset" {
		t.Fatalf("DatasetVersion = %q", recap.Snapshot.DatasetVersion)
	}
	if len(repositories.saved.Behaviors) != 1 ||
		repositories.saved.Behaviors[0].BehaviorTypeDefinitionID != behaviorID ||
		!repositories.saved.Behaviors[0].IsPrimary {
		t.Fatalf("saved behaviors = %#v", repositories.saved.Behaviors)
	}
	if len(repositories.saved.Achievements) != 1 ||
		repositories.saved.Achievements[0].AchievementDefinitionID != achievementID ||
		!repositories.saved.Achievements[0].IsShareable {
		t.Fatalf("saved achievements = %#v", repositories.saved.Achievements)
	}
	if repositories.saved.NextAction == nil || repositories.saved.NextAction.Code != "open_personal_collection" {
		t.Fatalf("saved next action = %#v", repositories.saved.NextAction)
	}
}

type fakeRepositories struct {
	profile  model.ProfileSummary
	dataset  model.VersionedDataset
	catalogs catalog.Set
	saved    *model.RecapDraft
}

func (f *fakeRepositories) ListProfiles(context.Context) ([]model.ProfileSummary, error) {
	return []model.ProfileSummary{f.profile}, nil
}

func (f *fakeRepositories) GetProfile(context.Context, uuid.UUID) (model.ProfileSummary, error) {
	return f.profile, nil
}

func (f *fakeRepositories) LoadDataset(
	_ context.Context,
	_ model.User,
	period model.Period,
	cutoff time.Time,
) (model.VersionedDataset, error) {
	f.dataset.Value.Period = period
	f.dataset.Value.DataCutoffAt = cutoff
	return f.dataset, nil
}

func (f *fakeRepositories) LoadCatalogs(context.Context) (catalog.Set, error) {
	return f.catalogs, nil
}

func (f *fakeRepositories) FindRecap(context.Context, model.RecapIdentity) (model.Recap, error) {
	return model.Recap{}, repository.ErrRecapNotFound
}

func (f *fakeRepositories) GetRecap(context.Context, uuid.UUID) (model.Recap, error) {
	return model.Recap{}, errors.New("unexpected GetRecap call")
}

func (f *fakeRepositories) SaveRecap(_ context.Context, draft model.RecapDraft) (model.Recap, error) {
	f.saved = &draft
	return model.Recap{Snapshot: draft.Snapshot}, nil
}

type ruleEngineFunc func(
	engine.Dataset,
	catalog.AchievementCatalog,
	catalog.BehaviorCatalog,
) (engine.Result, error)

func (f ruleEngineFunc) Run(
	dataset engine.Dataset,
	achievements catalog.AchievementCatalog,
	behaviors catalog.BehaviorCatalog,
) (engine.Result, error) {
	return f(dataset, achievements, behaviors)
}

type fixedClock struct{ value time.Time }

func (clock fixedClock) Now() time.Time { return clock.value }
