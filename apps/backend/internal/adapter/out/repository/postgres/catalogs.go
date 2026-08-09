package postgres

import (
	"context"
	"errors"
	"fmt"

	"avito-recap/internal/catalog"
	"avito-recap/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
	pool *pgxpool.Pool
}

func NewCatalogRepository(pool *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{pool: pool}
}

const latestCatalogVersionQuery = `
	SELECT id, version
	FROM catalog_versions
	WHERE kind = $1
	ORDER BY imported_at DESC, id DESC
	LIMIT 1
`

const loadBehaviorDefinitionsQuery = `
	SELECT id, code, name, description, rule, default_action, enabled, sort_order
	FROM behavior_type_definitions
	WHERE catalog_version_id = $1
	ORDER BY sort_order, code
`

const loadAchievementDefinitionsQuery = `
	SELECT id, code, name, description, rule, icon_key, enabled, shareable_by_default, sort_order
	FROM achievement_definitions
	WHERE catalog_version_id = $1
	ORDER BY sort_order, code
`

func (r *CatalogRepository) LoadCatalogs(ctx context.Context) (catalog.Set, error) {
	const op = "postgres.CatalogRepository.LoadCatalogs"

	behaviorVersionID, behaviorVersion, err := r.latestCatalogVersion(ctx, "behavior")
	if err != nil {
		return catalog.Set{}, fmt.Errorf("%s: %w", op, err)
	}
	achievementVersionID, achievementVersion, err := r.latestCatalogVersion(ctx, "achievement")
	if err != nil {
		return catalog.Set{}, fmt.Errorf("%s: %w", op, err)
	}

	behaviors, behaviorIDs, err := r.loadBehaviors(ctx, behaviorVersionID)
	if err != nil {
		return catalog.Set{}, fmt.Errorf("%s: %w", op, err)
	}
	achievements, achievementIDs, err := r.loadAchievements(ctx, achievementVersionID)
	if err != nil {
		return catalog.Set{}, fmt.Errorf("%s: %w", op, err)
	}

	result := catalog.Set{
		Behaviors: catalog.BehaviorCatalog{
			Version:   behaviorVersion,
			Behaviors: behaviors,
		},
		Achievements: catalog.AchievementCatalog{
			Version:      achievementVersion,
			Achievements: achievements,
		},
		BehaviorIDs:          behaviorIDs,
		AchievementIDs:       achievementIDs,
		BehaviorVersionID:    behaviorVersionID,
		AchievementVersionID: achievementVersionID,
	}
	if err := result.Behaviors.Validate(); err != nil {
		return catalog.Set{}, fmt.Errorf("%s: %w", op, err)
	}
	if err := result.Achievements.Validate(); err != nil {
		return catalog.Set{}, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

type catalogVersionRow struct {
	ID      uuid.UUID
	Version string
}

func (r *CatalogRepository) latestCatalogVersion(ctx context.Context, kind string) (uuid.UUID, string, error) {
	rows, err := r.pool.Query(ctx, latestCatalogVersionQuery, kind)
	if err != nil {
		return uuid.Nil, "", err
	}
	value, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[catalogVersionRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, "", repository.ErrCatalogUnavailable
	}
	if err != nil {
		return uuid.Nil, "", err
	}
	return value.ID, value.Version, nil
}

type behaviorDefinitionRow struct {
	ID            uuid.UUID
	Code          string
	Name          string
	Description   string
	Rule          map[string]any
	DefaultAction catalog.DefaultAction
	Enabled       bool
	SortOrder     int
}

func (r *CatalogRepository) loadBehaviors(
	ctx context.Context,
	versionID uuid.UUID,
) ([]catalog.BehaviorDefinition, map[string]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, loadBehaviorDefinitionsQuery, versionID)
	if err != nil {
		return nil, nil, err
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByPos[behaviorDefinitionRow])
	if err != nil {
		return nil, nil, err
	}

	definitions := make([]catalog.BehaviorDefinition, 0, len(values))
	ids := make(map[string]uuid.UUID, len(values))
	for _, value := range values {
		enabled := value.Enabled
		definition := catalog.BehaviorDefinition{
			Code:          value.Code,
			Name:          value.Name,
			Description:   value.Description,
			Enabled:       &enabled,
			SortOrder:     value.SortOrder,
			Rule:          value.Rule,
			DefaultAction: value.DefaultAction,
		}
		definitions = append(definitions, definition)
		ids[value.Code] = value.ID
	}
	return definitions, ids, nil
}

type achievementDefinitionRow struct {
	ID                 uuid.UUID
	Code               string
	Name               string
	Description        string
	Rule               map[string]any
	IconKey            string
	Enabled            bool
	ShareableByDefault bool
	SortOrder          int
}

func (r *CatalogRepository) loadAchievements(
	ctx context.Context,
	versionID uuid.UUID,
) ([]catalog.AchievementDefinition, map[string]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, loadAchievementDefinitionsQuery, versionID)
	if err != nil {
		return nil, nil, err
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByPos[achievementDefinitionRow])
	if err != nil {
		return nil, nil, err
	}

	definitions := make([]catalog.AchievementDefinition, 0, len(values))
	ids := make(map[string]uuid.UUID, len(values))
	for _, value := range values {
		enabled := value.Enabled
		definition := catalog.AchievementDefinition{
			Code:               value.Code,
			Name:               value.Name,
			Description:        value.Description,
			IconKey:            value.IconKey,
			Enabled:            &enabled,
			ShareableByDefault: value.ShareableByDefault,
			SortOrder:          value.SortOrder,
			Rule:               value.Rule,
		}
		definitions = append(definitions, definition)
		ids[value.Code] = value.ID
	}
	return definitions, ids, nil
}
