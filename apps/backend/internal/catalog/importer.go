package catalog

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type catalogKind string

const (
	achievementCatalogKind catalogKind = "achievement"
	behaviorCatalogKind    catalogKind = "behavior"
)

func ImportCatalogs(
	ctx context.Context,
	pool *pgxpool.Pool,
	achievementCatalog AchievementCatalog,
	behaviorCatalog BehaviorCatalog,
) error {
	achievementHash, err := fingerprintCatalog(achievementCatalog)
	if err != nil {
		return fmt.Errorf("fingerprint achievement catalog: %w", err)
	}
	behaviorHash, err := fingerprintCatalog(behaviorCatalog)
	if err != nil {
		return fmt.Errorf("fingerprint behavior catalog: %w", err)
	}

	err = pgx.BeginFunc(ctx, pool, func(transaction pgx.Tx) error {
		achievementVersionID, err := ensureCatalogVersion(
			ctx,
			transaction,
			achievementCatalogKind,
			achievementCatalog.Version,
			achievementHash,
		)
		if err != nil {
			return err
		}

		behaviorVersionID, err := ensureCatalogVersion(
			ctx,
			transaction,
			behaviorCatalogKind,
			behaviorCatalog.Version,
			behaviorHash,
		)
		if err != nil {
			return err
		}

		if err := insertAchievements(
			ctx,
			transaction,
			achievementVersionID,
			achievementCatalog.Achievements,
		); err != nil {
			return err
		}

		return insertBehaviors(
			ctx,
			transaction,
			behaviorVersionID,
			behaviorCatalog.Behaviors,
		)
	})
	if err != nil {
		return fmt.Errorf("import catalogs: %w", err)
	}

	return nil
}

func fingerprintCatalog(catalogValue any) (string, error) {
	content, err := json.Marshal(catalogValue)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(content)), nil
}
