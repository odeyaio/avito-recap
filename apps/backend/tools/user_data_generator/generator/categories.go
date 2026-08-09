package generator

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-recap/internal/model"
)

var categoryNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

func CategoryID(name string) uuid.UUID {
	return uuid.NewSHA1(categoryNamespace, []byte(name))
}

func categoryCode(name string) string {
	id := uuid.NewSHA1(categoryNamespace, []byte("code:"+name))
	return "cat_" + strings.ReplaceAll(id.String()[:8], "-", "")
}

func BuildCategories(configs []CategoryConfig) []model.Category {
	categories := make([]model.Category, 0, len(configs))
	for _, cfg := range configs {
		categories = append(categories, model.Category{
			ID:   CategoryID(cfg.Name),
			Code: categoryCode(cfg.Name),
			Name: cfg.Name,
		})
	}
	return categories
}

func InsertCategories(ctx context.Context, pool *pgxpool.Pool, categories []model.Category) error {
	if len(categories) == 0 {
		return nil
	}

	for _, category := range categories {
		_, err := pool.Exec(ctx, `
			INSERT INTO categories (id, code, name, parent_id)
			VALUES ($1, $2, $3, NULL)
		`, category.ID, category.Code, category.Name)
		if err != nil {
			return fmt.Errorf("insert category %s: %w", category.Name, err)
		}
	}

	return nil
}
