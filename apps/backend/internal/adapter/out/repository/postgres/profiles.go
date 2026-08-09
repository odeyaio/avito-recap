package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"avito-recap/internal/model"
	"avito-recap/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

const listProfilesQuery = `
	SELECT
		u.id,
		u.display_name,
		u.registered_at,
		u.region,
		u.timezone,
		u.is_test_profile,
		u.created_at,
		COALESCE((
			SELECT array_agg(year ORDER BY year DESC)
			FROM (
				SELECT DISTINCT EXTRACT(YEAR FROM ua.occurred_at AT TIME ZONE u.timezone)::INTEGER AS year
				FROM user_activity AS ua
				WHERE ua.user_id = u.id
			) AS years
		), '{}'::INTEGER[]),
		latest_recap.id
	FROM users AS u
	LEFT JOIN LATERAL (
		SELECT snapshots.id
		FROM recap_snapshots AS snapshots
		WHERE snapshots.user_id = u.id
		ORDER BY snapshots.generated_at DESC, snapshots.id
		LIMIT 1
	) AS latest_recap ON TRUE
	WHERE u.is_test_profile
	ORDER BY u.display_name, u.id
`

const getProfileQuery = `
	SELECT
		u.id,
		u.display_name,
		u.registered_at,
		u.region,
		u.timezone,
		u.is_test_profile,
		u.created_at,
		COALESCE((
			SELECT array_agg(year ORDER BY year DESC)
			FROM (
				SELECT DISTINCT EXTRACT(YEAR FROM ua.occurred_at AT TIME ZONE u.timezone)::INTEGER AS year
				FROM user_activity AS ua
				WHERE ua.user_id = u.id
			) AS years
		), '{}'::INTEGER[]),
		latest_recap.id
	FROM users AS u
	LEFT JOIN LATERAL (
		SELECT snapshots.id
		FROM recap_snapshots AS snapshots
		WHERE snapshots.user_id = u.id
		ORDER BY snapshots.generated_at DESC, snapshots.id
		LIMIT 1
	) AS latest_recap ON TRUE
	WHERE u.id = $1 AND u.is_test_profile
`

type profileRow struct {
	ID             uuid.UUID
	DisplayName    string
	RegisteredAt   time.Time
	Region         string
	Timezone       string
	IsTestProfile  bool
	CreatedAt      time.Time
	AvailableYears []int32
	LatestRecapID  *uuid.UUID
}

func (row profileRow) profileSummary() model.ProfileSummary {
	return model.ProfileSummary{
		User: model.User{
			ID:            row.ID,
			DisplayName:   row.DisplayName,
			RegisteredAt:  row.RegisteredAt,
			Region:        row.Region,
			Timezone:      row.Timezone,
			IsTestProfile: row.IsTestProfile,
			CreatedAt:     row.CreatedAt,
		},
		AvailableYears: row.AvailableYears,
		LatestRecapID:  row.LatestRecapID,
	}
}

func (r *ProfileRepository) ListProfiles(ctx context.Context) ([]model.ProfileSummary, error) {
	const op = "postgres.ProfileRepository.ListProfiles"

	rows, err := r.pool.Query(ctx, listProfilesQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByPos[profileRow])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	profiles := make([]model.ProfileSummary, 0, len(values))
	for _, value := range values {
		profiles = append(profiles, value.profileSummary())
	}
	return profiles, nil
}

func (r *ProfileRepository) GetProfile(ctx context.Context, profileID uuid.UUID) (model.ProfileSummary, error) {
	const op = "postgres.ProfileRepository.GetProfile"

	rows, err := r.pool.Query(ctx, getProfileQuery, profileID)
	if err != nil {
		return model.ProfileSummary{}, fmt.Errorf("%s: %w", op, err)
	}
	value, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[profileRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ProfileSummary{}, fmt.Errorf("%s: %w", op, repository.ErrProfileNotFound)
	}
	if err != nil {
		return model.ProfileSummary{}, fmt.Errorf("%s: %w", op, err)
	}
	return value.profileSummary(), nil
}
