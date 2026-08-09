package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"avito-recap/internal/engine"
	"avito-recap/internal/model"
	"avito-recap/internal/repository"

	"github.com/google/uuid"
)

type RecapService struct {
	profiles profileRepository
	datasets datasetRepository
	catalogs catalogRepository
	recaps   recapRepository
	engine   ruleEngine
	actions  actionResolver
	clock    clock
}

func NewRecapService(
	profiles profileRepository,
	datasets datasetRepository,
	catalogs catalogRepository,
	recaps recapRepository,
	ruleEngine ruleEngine,
	actions actionResolver,
) *RecapService {
	return &RecapService{
		profiles: profiles,
		datasets: datasets,
		catalogs: catalogs,
		recaps:   recaps,
		engine:   ruleEngine,
		actions:  actions,
		clock:    systemClock{},
	}
}

func (s *RecapService) GetRecap(ctx context.Context, recapID uuid.UUID) (model.Recap, error) {
	const op = "service.RecapService.GetRecap"

	recap, err := s.recaps.GetRecap(ctx, recapID)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	return recap, nil
}

func (s *RecapService) GenerateRecap(ctx context.Context, profileID uuid.UUID, year int) (model.Recap, error) {
	const op = "service.RecapService.GenerateRecap"

	profile, err := s.profiles.GetProfile(ctx, profileID)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: get profile: %w", op, err)
	}

	period, err := engine.PeriodForYear(year, profile.User.Timezone)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: build recap period: %w", op, err)
	}
	now := s.clock.Now()
	cutoff := now
	if cutoff.After(period.End) {
		cutoff = period.End
	}

	catalogs, err := s.catalogs.LoadCatalogs(ctx)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: load catalogs: %w", op, err)
	}
	dataset, err := s.datasets.LoadDataset(ctx, profile.User, period, cutoff)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: load dataset: %w", op, err)
	}

	identity := model.RecapIdentity{
		UserID:         profile.User.ID,
		PeriodStart:    period.Start,
		PeriodEnd:      period.End,
		RulesetVersion: catalogs.RulesetVersion(),
		DatasetVersion: dataset.Version,
	}
	existing, err := s.recaps.FindRecap(ctx, identity)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, repository.ErrRecapNotFound) {
		return model.Recap{}, fmt.Errorf("%s: find existing recap: %w", op, err)
	}

	result, err := s.engine.Run(dataset.Value, catalogs.Achievements, catalogs.Behaviors)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: run recap engine: %w", op, err)
	}
	primary, ok := primaryBehavior(result.Behaviors)
	if !ok {
		return model.Recap{}, fmt.Errorf("%s: %w", op, ErrBehaviorNotMatched)
	}

	metrics, err := result.MetricsJSON()
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: marshal recap metrics: %w", op, err)
	}

	recapID := uuid.New()
	draft := model.RecapDraft{
		Snapshot: model.RecapSnapshot{
			ID:             recapID,
			UserID:         profile.User.ID,
			PeriodStart:    period.Start,
			PeriodEnd:      period.End,
			RulesetVersion: identity.RulesetVersion,
			DatasetVersion: identity.DatasetVersion,
			DataCutoffAt:   cutoff,
			Metrics:        metrics,
			GeneratedAt:    now,
		},
	}

	for _, match := range result.Behaviors {
		definitionID, exists := catalogs.BehaviorIDs[match.Definition.Code]
		if !exists {
			return model.Recap{}, fmt.Errorf("%s: behavior definition %q has no database id", op, match.Definition.Code)
		}
		evidence, marshalErr := json.Marshal(match.Evidence)
		if marshalErr != nil {
			return model.Recap{}, fmt.Errorf("%s: marshal behavior %s evidence: %w", op, match.Definition.Code, marshalErr)
		}
		draft.Behaviors = append(draft.Behaviors, model.RecapBehaviorType{
			RecapID:                  recapID,
			BehaviorTypeDefinitionID: definitionID,
			IsPrimary:                match.IsPrimary,
			Position:                 match.Position,
			Evidence:                 evidence,
		})
	}

	for _, match := range result.Achievements {
		definitionID, exists := catalogs.AchievementIDs[match.Definition.Code]
		if !exists {
			return model.Recap{}, fmt.Errorf("%s: achievement definition %q has no database id", op, match.Definition.Code)
		}
		evidence, marshalErr := json.Marshal(match.Evidence)
		if marshalErr != nil {
			return model.Recap{}, fmt.Errorf("%s: marshal achievement %s evidence: %w", op, match.Definition.Code, marshalErr)
		}
		draft.Achievements = append(draft.Achievements, model.RecapAchievement{
			RecapID:                 recapID,
			AchievementDefinitionID: definitionID,
			Position:                match.Position,
			Evidence:                evidence,
			IsShareable:             match.Definition.ShareableByDefault,
		})
	}

	action, err := s.actions.Resolve(ctx, primary.Definition.DefaultAction, dataset.Value, primary.Evidence)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: resolve next action: %w", op, err)
	}
	action.RecapID = recapID
	action.ResolvedAt = now
	draft.NextAction = &action

	recap, err := s.recaps.SaveRecap(ctx, draft)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: save recap: %w", op, err)
	}
	return recap, nil
}

func primaryBehavior(matches []engine.BehaviorMatch) (engine.BehaviorMatch, bool) {
	for _, match := range matches {
		if match.IsPrimary {
			return match, true
		}
	}
	return engine.BehaviorMatch{}, false
}
