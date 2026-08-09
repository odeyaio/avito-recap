package engine

import (
	"cmp"
	"fmt"
	"slices"

	"avito-recap/internal/catalog"
)

type Engine struct {
	config Config
}

func New(config Config) (*Engine, error) {
	if config.FreshListingWindow <= 0 {
		return nil, fmt.Errorf("fresh listing window must be positive")
	}
	if config.LongGap <= 0 {
		return nil, fmt.Errorf("long gap must be positive")
	}
	if config.FastContactWindow <= 0 {
		return nil, fmt.Errorf("fast contact window must be positive")
	}
	if config.FastListingResponseWindow <= 0 {
		return nil, fmt.Errorf("fast listing response window must be positive")
	}
	if config.SignificantCategoryEvents <= 0 {
		return nil, fmt.Errorf("significant category event count must be positive")
	}
	if config.LowSupplyResultCount < 0 {
		return nil, fmt.Errorf("low supply result count must not be negative")
	}
	if config.MinimumActions < 0 {
		return nil, fmt.Errorf("minimum actions must not be negative")
	}

	return &Engine{config: config}, nil
}

func (e *Engine) Run(
	dataset Dataset,
	achievementCatalog catalog.AchievementCatalog,
	behaviorCatalog catalog.BehaviorCatalog,
) (Result, error) {
	if err := dataset.Period.Validate(); err != nil {
		return Result{}, fmt.Errorf("validate period: %w", err)
	}
	if err := achievementCatalog.Validate(); err != nil {
		return Result{}, fmt.Errorf("validate achievement catalog: %w", err)
	}
	if err := behaviorCatalog.Validate(); err != nil {
		return Result{}, fmt.Errorf("validate behavior catalog: %w", err)
	}

	metrics, err := e.CalculateMetrics(dataset)
	if err != nil {
		return Result{}, err
	}
	if actions, ok := intMetric(metrics, "activity.total_actions"); !ok || actions < int64(e.config.MinimumActions) {
		return Result{}, ErrNoActivity
	}

	behaviors, err := matchBehaviors(behaviorCatalog.Behaviors, metrics)
	if err != nil {
		return Result{}, err
	}
	achievements, err := matchAchievements(achievementCatalog.Achievements, metrics)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Metrics:      metrics,
		Behaviors:    behaviors,
		Achievements: achievements,
	}, nil
}

func matchBehaviors(definitions []catalog.BehaviorDefinition, metrics Metrics) ([]BehaviorMatch, error) {
	sorted := append([]catalog.BehaviorDefinition(nil), definitions...)
	slices.SortFunc(sorted, func(left, right catalog.BehaviorDefinition) int {
		return cmp.Or(cmp.Compare(left.SortOrder, right.SortOrder), cmp.Compare(left.Code, right.Code))
	})

	matches := make([]BehaviorMatch, 0)
	for _, definition := range sorted {
		if !definition.IsEnabled() {
			continue
		}

		matched, evidence, err := EvaluateRule(definition.Rule, metrics)
		if err != nil {
			return nil, fmt.Errorf("evaluate behavior %s: %w", definition.Code, err)
		}
		if !matched {
			continue
		}

		matches = append(matches, BehaviorMatch{
			Definition: definition,
			IsPrimary:  len(matches) == 0,
			Position:   len(matches) + 1,
			Evidence:   evidence,
		})
	}

	return matches, nil
}

func matchAchievements(definitions []catalog.AchievementDefinition, metrics Metrics) ([]AchievementMatch, error) {
	sorted := append([]catalog.AchievementDefinition(nil), definitions...)
	slices.SortFunc(sorted, func(left, right catalog.AchievementDefinition) int {
		return cmp.Or(cmp.Compare(left.SortOrder, right.SortOrder), cmp.Compare(left.Code, right.Code))
	})

	matches := make([]AchievementMatch, 0)
	for _, definition := range sorted {
		if !definition.IsEnabled() {
			continue
		}

		matched, evidence, err := EvaluateRule(definition.Rule, metrics)
		if err != nil {
			return nil, fmt.Errorf("evaluate achievement %s: %w", definition.Code, err)
		}
		if !matched {
			continue
		}

		matches = append(matches, AchievementMatch{
			Definition: definition,
			Position:   len(matches) + 1,
			Evidence:   evidence,
		})
	}

	return matches, nil
}

func intMetric(metrics Metrics, code string) (int64, bool) {
	value, exists := metrics.Get(code)
	if !exists {
		return 0, false
	}
	numeric, ok := number(value)
	return int64(numeric), ok
}
