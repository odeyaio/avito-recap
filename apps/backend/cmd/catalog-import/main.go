package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"avito-recap/internal/catalog"

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

	if err := catalog.ImportCatalogs(ctx, pool, achievementCatalog, behaviorCatalog); err != nil {
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
