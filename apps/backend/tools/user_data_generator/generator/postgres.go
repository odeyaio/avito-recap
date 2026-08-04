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
			INSERT INTO users (username, register_date, region)
			VALUES ($1, $2, $3)
		`, user.Username, user.RegisterDate, user.Region)
		if err != nil {
			return fmt.Errorf("insert user %s: %w", user.Username, err)
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
			INSERT INTO listings (seller_id, category, price, published_at, closed_at, delivery_available)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, listing.SellerID, listing.Category, listing.Price, listing.PublishedAt, listing.ClosedAt, listing.DeliveryAvailable)
		if err != nil {
			return fmt.Errorf("insert listing %d: %w", listing.ID, err)
		}
	}

	return nil
}

func InsertSearches(ctx context.Context, pool *pgxpool.Pool, searches []model.Search) error {
	if len(searches) == 0 {
		return nil
	}

	for _, search := range searches {
		_, err := pool.Exec(ctx, `
			INSERT INTO searches (user_id, searched_at, topic, category, filters, result_count)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, search.UserID, search.SearchedAt, search.Topic, search.Category, search.Filters, search.ResultCount)
		if err != nil {
			return fmt.Errorf("insert search %d: %w", search.ID, err)
		}
	}

	return nil
}

func InsertViews(ctx context.Context, pool *pgxpool.Pool, views []model.View) error {
	if len(views) == 0 {
		return nil
	}

	for _, view := range views {
		_, err := pool.Exec(ctx, `
			INSERT INTO views (user_id, listing_id, category, viewed_at, duration_seconds, is_repeat)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, view.UserID, view.ListingID, view.Category, view.ViewedAt, view.DurationSeconds, view.IsRepeat)
		if err != nil {
			return fmt.Errorf("insert view %d: %w", view.ID, err)
		}
	}

	return nil
}

func InsertFavorites(ctx context.Context, pool *pgxpool.Pool, favorites []model.Favorite) error {
	if len(favorites) == 0 {
		return nil
	}

	for _, favorite := range favorites {
		_, err := pool.Exec(ctx, `
			INSERT INTO favorites (user_id, listing_id, category, action, occurred_at)
			VALUES ($1, $2, $3, $4, $5)
		`, favorite.UserID, favorite.ListingID, favorite.Category, favorite.Action, favorite.OccurredAt)
		if err != nil {
			return fmt.Errorf("insert favorite %d: %w", favorite.ID, err)
		}
	}

	return nil
}

func InsertContacts(ctx context.Context, pool *pgxpool.Pool, contacts []model.Contact) error {
	if len(contacts) == 0 {
		return nil
	}

	for _, contact := range contacts {
		_, err := pool.Exec(ctx, `
			INSERT INTO contacts (user_id, listing_id, contact_type, occurred_at)
			VALUES ($1, $2, $3, $4)
		`, contact.UserID, contact.ListingID, contact.ContactType, contact.OccurredAt)
		if err != nil {
			return fmt.Errorf("insert contact %d: %w", contact.ID, err)
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
			INSERT INTO deals (buyer_id, seller_id, listing_id, category, price, delivery, completed_at, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, deal.BuyerID, deal.SellerID, deal.ListingID, deal.Category, deal.Price, deal.Delivery, deal.CompletedAt, deal.Status)
		if err != nil {
			return fmt.Errorf("insert deal %d: %w", deal.ID, err)
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
			INSERT INTO reviews (reviewer_id, reviewed_user_id, deal_id, rating, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, review.ReviewerID, review.ReviewedUserID, review.DealID, review.Rating, review.CreatedAt)
		if err != nil {
			return fmt.Errorf("insert review %d: %w", review.ID, err)
		}
	}

	return nil
}

func InsertUserEvents(ctx context.Context, pool *pgxpool.Pool, userEvents []model.UserEvent) error {
	if len(userEvents) == 0 {
		return nil
	}

	for _, userEvent := range userEvents {
		_, err := pool.Exec(ctx, `
			INSERT INTO user_events (user_id, event_type, occurred_at)
			VALUES ($1, $2, $3)
		`, userEvent.UserID, userEvent.EventType, userEvent.OccurredAt)
		if err != nil {
			return fmt.Errorf("insert user event %d: %w", userEvent.ID, err)
		}
	}

	return nil
}
