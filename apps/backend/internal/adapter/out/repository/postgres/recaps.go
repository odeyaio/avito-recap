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

type RecapRepository struct {
	pool *pgxpool.Pool
}

func NewRecapRepository(pool *pgxpool.Pool) *RecapRepository {
	return &RecapRepository{pool: pool}
}

const findRecapQuery = `
	SELECT id
	FROM recap_snapshots
	WHERE user_id = $1
	  AND period_start = $2::DATE
	  AND period_end = $3::DATE
	  AND ruleset_version = $4
	  AND dataset_version = $5
`

const getRecapQuery = `
	SELECT
		s.id, s.user_id, s.period_start, s.period_end, s.ruleset_version,
		s.dataset_version, s.data_cutoff_at, s.metrics, s.generated_at,
		u.id, u.display_name, u.registered_at, u.region, u.timezone,
		u.is_test_profile, u.created_at,
		COALESCE((
			SELECT array_agg(year ORDER BY year DESC)
			FROM (
				SELECT DISTINCT EXTRACT(YEAR FROM ua.occurred_at AT TIME ZONE u.timezone)::INTEGER AS year
				FROM user_activity AS ua
				WHERE ua.user_id = u.id
			) AS years
		), '{}'::INTEGER[]),
		latest_recap.id
	FROM recap_snapshots AS s
	JOIN users AS u ON u.id = s.user_id
	LEFT JOIN LATERAL (
		SELECT snapshots.id
		FROM recap_snapshots AS snapshots
		WHERE snapshots.user_id = u.id
		ORDER BY snapshots.generated_at DESC, snapshots.id
		LIMIT 1
	) AS latest_recap ON TRUE
	WHERE s.id = $1
`

const getRecapBehaviorsQuery = `
	SELECT
		r.recap_id, r.behavior_type_definition_id, r.is_primary, r.position,
		r.score, r.evidence,
		d.id, d.catalog_version_id, d.code, d.name, d.description, d.rule,
		d.default_action, d.enabled, d.sort_order, d.updated_at
	FROM recap_behavior_types AS r
	JOIN behavior_type_definitions AS d ON d.id = r.behavior_type_definition_id
	WHERE r.recap_id = $1
	ORDER BY r.position, d.code
`

const getRecapAchievementsQuery = `
	SELECT
		r.recap_id, r.achievement_definition_id, r.position, r.achieved_at,
		r.evidence, r.is_shareable,
		d.id, d.catalog_version_id, d.code, d.name, d.description, d.rule,
		d.icon_key, d.enabled, d.shareable_by_default, d.sort_order, d.updated_at
	FROM recap_achievements AS r
	JOIN achievement_definitions AS d ON d.id = r.achievement_definition_id
	WHERE r.recap_id = $1
	ORDER BY r.position, d.code
`

const getRecapNextActionQuery = `
	SELECT recap_id, code, href, target, evidence, resolved_at
	FROM recap_next_action
	WHERE recap_id = $1
`

const getRecapPresentationQuery = `
	SELECT id, recap_id, locale, prompt_version, model_name, input_hash, content, generated_at
	FROM recap_presentations
	WHERE recap_id = $1
	ORDER BY generated_at DESC, id DESC
	LIMIT 1
`

const saveRecapQuery = `
	INSERT INTO recap_snapshots (
		id, user_id, period_start, period_end, ruleset_version, dataset_version,
		data_cutoff_at, metrics, generated_at
	)
	VALUES ($1, $2, $3::DATE, $4::DATE, $5, $6, $7, $8, $9)
	ON CONFLICT (user_id, period_start, period_end, ruleset_version, dataset_version)
	DO UPDATE SET
		data_cutoff_at = EXCLUDED.data_cutoff_at,
		metrics = EXCLUDED.metrics,
		generated_at = EXCLUDED.generated_at
	RETURNING id, generated_at
`

const saveBehaviorQuery = `
	INSERT INTO recap_behavior_types (
		recap_id, behavior_type_definition_id, is_primary, position, score, evidence
	)
	VALUES ($1, $2, $3, $4, $5, $6)
`

const saveAchievementQuery = `
	INSERT INTO recap_achievements (
		recap_id, achievement_definition_id, position, achieved_at, evidence, is_shareable
	)
	VALUES ($1, $2, $3, $4, $5, $6)
`

const saveNextActionQuery = `
	INSERT INTO recap_next_action (recap_id, code, href, target, evidence, resolved_at)
	VALUES ($1, $2, $3, $4, $5, $6)
`

type querier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

type savedRecapRow struct {
	ID          uuid.UUID
	GeneratedAt time.Time
}

type recapRow struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	PeriodStart    time.Time
	PeriodEnd      time.Time
	RulesetVersion string
	DatasetVersion string
	DataCutoffAt   time.Time
	Metrics        []byte
	GeneratedAt    time.Time
	ProfileID      uuid.UUID
	DisplayName    string
	RegisteredAt   time.Time
	Region         string
	Timezone       string
	IsTestProfile  bool
	CreatedAt      time.Time
	AvailableYears []int32
	LatestRecapID  *uuid.UUID
}

func (row recapRow) recap() model.Recap {
	return model.Recap{
		Snapshot: model.RecapSnapshot{
			ID: row.ID, UserID: row.UserID, PeriodStart: row.PeriodStart, PeriodEnd: row.PeriodEnd,
			RulesetVersion: row.RulesetVersion, DatasetVersion: row.DatasetVersion,
			DataCutoffAt: row.DataCutoffAt, Metrics: row.Metrics, GeneratedAt: row.GeneratedAt,
		},
		Profile: model.ProfileSummary{
			User: model.User{
				ID: row.ProfileID, DisplayName: row.DisplayName, RegisteredAt: row.RegisteredAt,
				Region: row.Region, Timezone: row.Timezone,
				IsTestProfile: row.IsTestProfile, CreatedAt: row.CreatedAt,
			},
			AvailableYears: row.AvailableYears,
			LatestRecapID:  row.LatestRecapID,
		},
	}
}

type recapBehaviorRow struct {
	RecapID                  uuid.UUID
	BehaviorTypeDefinitionID uuid.UUID
	IsPrimary                bool
	Position                 int
	Score                    *float64
	Evidence                 []byte
	DefinitionID             uuid.UUID
	CatalogVersionID         uuid.UUID
	Code                     string
	Name                     string
	Description              string
	Rule                     []byte
	DefaultAction            []byte
	Enabled                  bool
	SortOrder                int
	UpdatedAt                time.Time
}

func (row recapBehaviorRow) storedBehavior() model.StoredBehavior {
	return model.StoredBehavior{
		Match: model.RecapBehaviorType{
			RecapID: row.RecapID, BehaviorTypeDefinitionID: row.BehaviorTypeDefinitionID,
			IsPrimary: row.IsPrimary, Position: row.Position, Score: row.Score, Evidence: row.Evidence,
		},
		Definition: model.BehaviorTypeDefinition{
			ID: row.DefinitionID, CatalogVersionID: row.CatalogVersionID, Code: row.Code,
			Name: row.Name, Description: row.Description, Rule: row.Rule, DefaultAction: row.DefaultAction,
			Enabled: row.Enabled, SortOrder: row.SortOrder, UpdatedAt: row.UpdatedAt,
		},
	}
}

type recapAchievementRow struct {
	RecapID                 uuid.UUID
	AchievementDefinitionID uuid.UUID
	Position                int
	AchievedAt              *time.Time
	Evidence                []byte
	IsShareable             bool
	DefinitionID            uuid.UUID
	CatalogVersionID        uuid.UUID
	Code                    string
	Name                    string
	Description             string
	Rule                    []byte
	IconKey                 string
	Enabled                 bool
	ShareableByDefault      bool
	SortOrder               int
	UpdatedAt               time.Time
}

func (row recapAchievementRow) storedAchievement() model.StoredAchievement {
	return model.StoredAchievement{
		Match: model.RecapAchievement{
			RecapID: row.RecapID, AchievementDefinitionID: row.AchievementDefinitionID,
			Position: row.Position, AchievedAt: row.AchievedAt, Evidence: row.Evidence, IsShareable: row.IsShareable,
		},
		Definition: model.AchievementDefinition{
			ID: row.DefinitionID, CatalogVersionID: row.CatalogVersionID, Code: row.Code,
			Name: row.Name, Description: row.Description, Rule: row.Rule, IconKey: row.IconKey,
			Enabled: row.Enabled, ShareableByDefault: row.ShareableByDefault,
			SortOrder: row.SortOrder, UpdatedAt: row.UpdatedAt,
		},
	}
}

func (r *RecapRepository) FindRecap(ctx context.Context, identity model.RecapIdentity) (model.Recap, error) {
	const op = "postgres.RecapRepository.FindRecap"

	rows, err := r.pool.Query(
		ctx,
		findRecapQuery,
		identity.UserID,
		identity.PeriodStart.Format(time.DateOnly),
		identity.PeriodEnd.Format(time.DateOnly),
		identity.RulesetVersion,
		identity.DatasetVersion,
	)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	recapID, err := pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Recap{}, fmt.Errorf("%s: %w", op, repository.ErrRecapNotFound)
	}
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	recap, err := loadRecap(ctx, r.pool, recapID)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	return recap, nil
}

func (r *RecapRepository) GetRecap(ctx context.Context, recapID uuid.UUID) (model.Recap, error) {
	const op = "postgres.RecapRepository.GetRecap"

	recap, err := loadRecap(ctx, r.pool, recapID)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	return recap, nil
}

func (r *RecapRepository) SaveRecap(ctx context.Context, draft model.RecapDraft) (model.Recap, error) {
	const op = "postgres.RecapRepository.SaveRecap"

	var recapID uuid.UUID
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		rows, err := tx.Query(
			ctx,
			saveRecapQuery,
			draft.Snapshot.ID,
			draft.Snapshot.UserID,
			draft.Snapshot.PeriodStart.Format(time.DateOnly),
			draft.Snapshot.PeriodEnd.Format(time.DateOnly),
			draft.Snapshot.RulesetVersion,
			draft.Snapshot.DatasetVersion,
			draft.Snapshot.DataCutoffAt,
			draft.Snapshot.Metrics,
			draft.Snapshot.GeneratedAt,
		)
		if err != nil {
			return err
		}
		saved, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[savedRecapRow])
		if err != nil {
			return err
		}
		recapID = saved.ID
		draft.Snapshot.GeneratedAt = saved.GeneratedAt

		if _, err := tx.Exec(ctx, `DELETE FROM recap_behavior_types WHERE recap_id = $1`, recapID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `DELETE FROM recap_achievements WHERE recap_id = $1`, recapID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `DELETE FROM recap_next_action WHERE recap_id = $1`, recapID); err != nil {
			return err
		}

		for _, behavior := range draft.Behaviors {
			if _, err := tx.Exec(
				ctx,
				saveBehaviorQuery,
				recapID,
				behavior.BehaviorTypeDefinitionID,
				behavior.IsPrimary,
				behavior.Position,
				behavior.Score,
				behavior.Evidence,
			); err != nil {
				return err
			}
		}

		for _, achievement := range draft.Achievements {
			if _, err := tx.Exec(
				ctx,
				saveAchievementQuery,
				recapID,
				achievement.AchievementDefinitionID,
				achievement.Position,
				achievement.AchievedAt,
				achievement.Evidence,
				achievement.IsShareable,
			); err != nil {
				return err
			}
		}

		if draft.NextAction != nil {
			if _, err := tx.Exec(
				ctx,
				saveNextActionQuery,
				recapID,
				draft.NextAction.Code,
				draft.NextAction.Href,
				draft.NextAction.Target,
				draft.NextAction.Evidence,
				draft.NextAction.ResolvedAt,
			); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	recap, err := loadRecap(ctx, r.pool, recapID)
	if err != nil {
		return model.Recap{}, fmt.Errorf("%s: %w", op, err)
	}
	return recap, nil
}

func loadRecap(ctx context.Context, database querier, recapID uuid.UUID) (model.Recap, error) {
	rows, err := database.Query(ctx, getRecapQuery, recapID)
	if err != nil {
		return model.Recap{}, err
	}
	value, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[recapRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Recap{}, repository.ErrRecapNotFound
	}
	if err != nil {
		return model.Recap{}, err
	}
	recap := value.recap()

	if recap.Behaviors, err = loadRecapBehaviors(ctx, database, recapID); err != nil {
		return model.Recap{}, err
	}
	if recap.Achievements, err = loadRecapAchievements(ctx, database, recapID); err != nil {
		return model.Recap{}, err
	}

	rows, err = database.Query(ctx, getRecapNextActionQuery, recapID)
	if err != nil {
		return model.Recap{}, err
	}
	action, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.RecapNextAction])
	if err == nil {
		recap.NextAction = &action
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return model.Recap{}, err
	}

	rows, err = database.Query(ctx, getRecapPresentationQuery, recapID)
	if err != nil {
		return model.Recap{}, err
	}
	presentation, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.RecapPresentation])
	if err == nil {
		recap.Presentation = &presentation
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return model.Recap{}, err
	}

	return recap, nil
}

func loadRecapBehaviors(ctx context.Context, database querier, recapID uuid.UUID) ([]model.StoredBehavior, error) {
	rows, err := database.Query(ctx, getRecapBehaviorsQuery, recapID)
	if err != nil {
		return nil, err
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByPos[recapBehaviorRow])
	if err != nil {
		return nil, err
	}

	behaviors := make([]model.StoredBehavior, 0, len(values))
	for _, value := range values {
		behaviors = append(behaviors, value.storedBehavior())
	}
	return behaviors, nil
}

func loadRecapAchievements(
	ctx context.Context,
	database querier,
	recapID uuid.UUID,
) ([]model.StoredAchievement, error) {
	rows, err := database.Query(ctx, getRecapAchievementsQuery, recapID)
	if err != nil {
		return nil, err
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByPos[recapAchievementRow])
	if err != nil {
		return nil, err
	}

	achievements := make([]model.StoredAchievement, 0, len(values))
	for _, value := range values {
		achievements = append(achievements, value.storedAchievement())
	}
	return achievements, nil
}
