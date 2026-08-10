package generator

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-recap/internal/model"
)

func ConnectToDatabase(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func InsertUsers(ctx context.Context, pool *pgxpool.Pool, users []model.User) error {
	if len(users) == 0 {
		return nil
	}

	for _, user := range users {
		_, err := pool.Exec(ctx, `
			INSERT INTO users (id, display_name, registered_at, region, timezone, is_test_profile, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		`, user.ID, user.DisplayName, user.RegisteredAt, user.Region, user.Timezone, user.IsTestProfile)
		if err != nil {
			return fmt.Errorf("insert user %s: %w", user.DisplayName, err)
		}
	}

	return nil
}

func InsertListings(ctx context.Context, pool *pgxpool.Pool, listings []model.Listing) error {
	if len(listings) == 0 {
		return nil
	}

	for _, listing := range listings {
		_, err := pool.Exec(ctx, `
			INSERT INTO listings (id, seller_id, category_id, region, price_band, published_at, closed_at, delivery_available, photo_count, description_complete)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, listing.ID, listing.SellerID, listing.CategoryID, listing.Region, listing.PriceBand, listing.PublishedAt, listing.ClosedAt, listing.DeliveryAvailable, listing.PhotoCount, listing.DescriptionComplete)
		if err != nil {
			return fmt.Errorf("insert listing %s: %w", listing.ID, err)
		}
	}

	return nil
}

func InsertActivityEvents(ctx context.Context, pool *pgxpool.Pool, events []model.ActivityEvent) error {
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		_, err := pool.Exec(ctx, `
			INSERT INTO activity_events (id, user_id, event_type, occurred_at, listing_id, category_id, duration_seconds, result_count, filter_count, topic_key, source_type, properties, ingested_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`, event.ID, event.UserID, event.Type, event.OccurredAt, event.ListingID, event.CategoryID, event.DurationSeconds, event.ResultCount, event.FilterCount, event.TopicKey, event.Source, event.Properties, event.IngestedAt)
		if err != nil {
			return fmt.Errorf("insert activity event %s: %w", event.ID, err)
		}
	}

	return nil
}

func InsertDeals(ctx context.Context, pool *pgxpool.Pool, deals []model.Deal) error {
	if len(deals) == 0 {
		return nil
	}

	for _, deal := range deals {
		_, err := pool.Exec(ctx, `
			INSERT INTO deals (id, listing_id, buyer_id, created_at, completed_at, status, delivery_used, price_band)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, deal.ID, deal.ListingID, deal.BuyerID, deal.CreatedAt, deal.CompletedAt, deal.Status, deal.DeliveryUsed, deal.PriceBand)
		if err != nil {
			return fmt.Errorf("insert deal %s: %w", deal.ID, err)
		}
	}

	return nil
}

func InsertReviews(ctx context.Context, pool *pgxpool.Pool, reviews []model.Review) error {
	if len(reviews) == 0 {
		return nil
	}

	for _, review := range reviews {
		_, err := pool.Exec(ctx, `
			INSERT INTO reviews (id, deal_id, author_id, recipient_id, rating, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, review.ID, review.DealID, review.AuthorID, review.RecipientID, review.Rating, review.CreatedAt)
		if err != nil {
			return fmt.Errorf("insert review %s: %w", review.ID, err)
		}
	}

	return nil
}

func InsertUserEvents(ctx context.Context, pool *pgxpool.Pool, events []model.ActivityEvent) error {
	return InsertActivityEvents(ctx, pool, events)
}
