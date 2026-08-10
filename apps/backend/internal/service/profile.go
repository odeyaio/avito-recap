package service

import (
	"context"
	"fmt"

	"avito-recap/internal/model"
)

type ProfileService struct {
	profiles profileRepository
}

func NewProfileService(profiles profileRepository) *ProfileService {
	return &ProfileService{profiles: profiles}
}

func (s *ProfileService) ListProfiles(ctx context.Context) ([]model.ProfileSummary, error) {
	const op = "service.ProfileService.ListProfiles"

	profiles, err := s.profiles.ListProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return profiles, nil
}
