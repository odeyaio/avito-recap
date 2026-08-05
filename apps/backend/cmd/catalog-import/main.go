package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"avito-recap/internal/catalog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const importTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var achievementCatalogPath string
	var behaviorCatalogPath string
	var databaseURL string
	flag.StringVar(
		&achievementCatalogPath,
		"achievements-file",
		"../../catalog/achievements.yaml",
		"path to the achievement catalog",
	)
	flag.StringVar(
		&behaviorCatalogPath,
		"behaviors-file",
		"../../catalog/behaviors.yaml",
		"path to the behavior catalog",
	)
	flag.StringVar(&databaseURL, "database-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection URL")
	flag.Parse()

	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	achievementCatalog, err := catalog.LoadAchievements(achievementCatalogPath)
	if err != nil {
		return err
	}
	behaviorCatalog, err := catalog.LoadBehaviors(behaviorCatalogPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), importTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	if err := importCatalogs(ctx, pool, achievementCatalog, behaviorCatalog); err != nil {
		return err
	}

	log.Printf(
		"imported %d achievement definitions from catalog %s and %d behavior definitions from catalog %s",
		len(achievementCatalog.Achievements),
		achievementCatalog.Version,
		len(behaviorCatalog.Behaviors),
		behaviorCatalog.Version,
	)

	return nil
}

func importCatalogs(
	ctx context.Context,
	pool *pgxpool.Pool,
	achievementCatalog catalog.AchievementCatalog,
	behaviorCatalog catalog.BehaviorCatalog,
) error {
	transaction, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin catalog import: %w", err)
	}
	defer func() { _ = transaction.Rollback(ctx) }()

	if err := upsertAchievements(ctx, transaction, achievementCatalog); err != nil {
		return err
	}
	if err := upsertBehaviors(ctx, transaction, behaviorCatalog); err != nil {
		return err
	}

	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("commit catalog import: %w", err)
	}

	return nil
}

func upsertAchievements(
	ctx context.Context,
	transaction pgx.Tx,
	achievementCatalog catalog.AchievementCatalog,
) error {
	for _, achievement := range achievementCatalog.Achievements {
		rule, err := json.Marshal(achievement.Rule)
		if err != nil {
			return fmt.Errorf("marshal rule for %s: %w", achievement.Code, err)
		}

		_, err = transaction.Exec(ctx, `
			INSERT INTO achievement_definitions (
				id,
				code,
				name,
				rule_description,
				rule,
				rule_version,
				icon_key,
				enabled,
				shareable_by_default,
				sort_order
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (code) DO UPDATE
			SET
				name = EXCLUDED.name,
				rule_description = EXCLUDED.rule_description,
				rule = EXCLUDED.rule,
				rule_version = EXCLUDED.rule_version,
				icon_key = EXCLUDED.icon_key,
				enabled = EXCLUDED.enabled,
				shareable_by_default = EXCLUDED.shareable_by_default,
				sort_order = EXCLUDED.sort_order,
				updated_at = CURRENT_TIMESTAMP
		`,
			uuid.New(),
			achievement.Code,
			achievement.Name,
			achievement.RuleDescription,
			rule,
			achievementCatalog.Version,
			achievement.IconKey,
			achievement.IsEnabled(),
			achievement.ShareableByDefault,
			achievement.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("upsert achievement %s: %w", achievement.Code, err)
		}
	}

	return nil
}

func upsertBehaviors(
	ctx context.Context,
	transaction pgx.Tx,
	behaviorCatalog catalog.BehaviorCatalog,
) error {
	for _, behavior := range behaviorCatalog.Behaviors {
		rule, err := json.Marshal(behavior.Rule)
		if err != nil {
			return fmt.Errorf("marshal rule for %s: %w", behavior.Code, err)
		}
		defaultAction, err := json.Marshal(behavior.DefaultAction)
		if err != nil {
			return fmt.Errorf("marshal default action for %s: %w", behavior.Code, err)
		}

		_, err = transaction.Exec(ctx, `
			INSERT INTO behavior_type_definitions (
				id,
				code,
				name,
				rule_description,
				rule,
				rule_version,
				default_action,
				enabled,
				sort_order
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (code) DO UPDATE
			SET
				name = EXCLUDED.name,
				rule_description = EXCLUDED.rule_description,
				rule = EXCLUDED.rule,
				rule_version = EXCLUDED.rule_version,
				default_action = EXCLUDED.default_action,
				enabled = EXCLUDED.enabled,
				sort_order = EXCLUDED.sort_order,
				updated_at = CURRENT_TIMESTAMP
		`,
			uuid.New(),
			behavior.Code,
			behavior.Name,
			behavior.RuleDescription,
			rule,
			behaviorCatalog.Version,
			defaultAction,
			behavior.IsEnabled(),
			behavior.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("upsert behavior %s: %w", behavior.Code, err)
		}
	}

	return nil
}
