package catalog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const ensureCatalogVersionQuery = `
	INSERT INTO catalog_versions (id, kind, version, content_hash)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (kind, version) DO UPDATE
	SET content_hash = catalog_versions.content_hash
	RETURNING id, content_hash
`

const insertAchievementQuery = `
	INSERT INTO achievement_definitions (
		id,
		catalog_version_id,
		code,
		name,
		rule_description,
		rule,
		icon_key,
		enabled,
		shareable_by_default,
		sort_order
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT (catalog_version_id, code) DO NOTHING
`

const insertBehaviorQuery = `
	INSERT INTO behavior_type_definitions (
		id,
		catalog_version_id,
		code,
		name,
		rule_description,
		rule,
		default_action,
		enabled,
		sort_order
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (catalog_version_id, code) DO NOTHING
`

func ensureCatalogVersion(
	ctx context.Context,
	transaction pgx.Tx,
	kind catalogKind,
	version string,
	contentHash string,
) (uuid.UUID, error) {
	versionID := uuid.New()
	var storedHash string
	err := transaction.QueryRow(
		ctx,
		ensureCatalogVersionQuery,
		versionID,
		kind,
		version,
		contentHash,
	).Scan(&versionID, &storedHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("register %s catalog version %s: %w", kind, version, err)
	}
	if storedHash != contentHash {
		return uuid.Nil, fmt.Errorf(
			"%s catalog version %s already exists with different content; bump the version",
			kind,
			version,
		)
	}

	return versionID, nil
}

func insertAchievements(
	ctx context.Context,
	transaction pgx.Tx,
	catalogVersionID uuid.UUID,
	achievements []AchievementDefinition,
) error {
	for _, achievement := range achievements {
		rule, err := json.Marshal(achievement.Rule)
		if err != nil {
			return fmt.Errorf("marshal rule for %s: %w", achievement.Code, err)
		}

		_, err = transaction.Exec(
			ctx,
			insertAchievementQuery,
			uuid.New(),
			catalogVersionID,
			achievement.Code,
			achievement.Name,
			achievement.RuleDescription,
			rule,
			achievement.IconKey,
			achievement.IsEnabled(),
			achievement.ShareableByDefault,
			achievement.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert achievement %s: %w", achievement.Code, err)
		}
	}

	return nil
}

func insertBehaviors(
	ctx context.Context,
	transaction pgx.Tx,
	catalogVersionID uuid.UUID,
	behaviors []BehaviorDefinition,
) error {
	for _, behavior := range behaviors {
		rule, err := json.Marshal(behavior.Rule)
		if err != nil {
			return fmt.Errorf("marshal rule for %s: %w", behavior.Code, err)
		}
		defaultAction, err := json.Marshal(behavior.DefaultAction)
		if err != nil {
			return fmt.Errorf("marshal default action for %s: %w", behavior.Code, err)
		}

		_, err = transaction.Exec(
			ctx,
			insertBehaviorQuery,
			uuid.New(),
			catalogVersionID,
			behavior.Code,
			behavior.Name,
			behavior.RuleDescription,
			rule,
			defaultAction,
			behavior.IsEnabled(),
			behavior.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("insert behavior %s: %w", behavior.Code, err)
		}
	}

	return nil
}
