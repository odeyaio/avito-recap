package service

import (
	"context"
	"errors"
	"time"

	"avito-recap/internal/catalog"
	"avito-recap/internal/engine"
	"avito-recap/internal/model"

	"github.com/google/uuid"
)

var ErrBehaviorNotMatched = errors.New("behavior not matched")

type profileRepository interface {
	ListProfiles(context.Context) ([]model.ProfileSummary, error)
	GetProfile(context.Context, uuid.UUID) (model.ProfileSummary, error)
}

type datasetRepository interface {
	LoadDataset(context.Context, model.User, model.Period, time.Time) (model.VersionedDataset, error)
}

type catalogRepository interface {
	LoadCatalogs(context.Context) (catalog.Set, error)
}

type recapRepository interface {
	FindRecap(context.Context, model.RecapIdentity) (model.Recap, error)
	GetRecap(context.Context, uuid.UUID) (model.Recap, error)
	SaveRecap(context.Context, model.RecapDraft) (model.Recap, error)
}

type ruleEngine interface {
	Run(engine.Dataset, catalog.AchievementCatalog, catalog.BehaviorCatalog) (engine.Result, error)
}

type actionResolver interface {
	Resolve(context.Context, catalog.DefaultAction, model.Dataset, []engine.Evidence) (model.RecapNextAction, error)
}

type clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }
