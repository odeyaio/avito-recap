package http

import (
	"context"

	"avito-recap/internal/model"

	"github.com/google/uuid"
)

type profileService interface {
	ListProfiles(context.Context) ([]model.ProfileSummary, error)
}

type recapService interface {
	GenerateRecap(context.Context, uuid.UUID, int) (model.Recap, error)
	GetRecap(context.Context, uuid.UUID) (model.Recap, error)
}
