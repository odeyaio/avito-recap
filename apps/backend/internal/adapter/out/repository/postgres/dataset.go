package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"avito-recap/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatasetRepository struct {
	pool *pgxpool.Pool
}

func NewDatasetRepository(pool *pgxpool.Pool) *DatasetRepository {
	return &DatasetRepository{pool: pool}
}

const loadEventsQuery = `
	SELECT DISTINCT
		e.id, e.user_id, e.event_type, e.occurred_at, e.listing_id, e.category_id,
		e.duration_seconds, e.result_count, e.filter_count, e.topic_key, e.source_type,
		e.properties, e.ingested_at
	FROM activity_events AS e
	LEFT JOIN listings AS l ON l.id = e.listing_id
	WHERE e.occurred_at < $2
	  AND (e.user_id = $1 OR l.seller_id = $1)
	ORDER BY e.occurred_at, e.id
`

const loadListingsQuery = `
	SELECT DISTINCT
		l.id, l.seller_id, l.category_id, l.region, l.price_band, l.published_at,
		l.closed_at, l.delivery_available, l.photo_count, l.description_complete
	FROM listings AS l
	WHERE l.published_at < $2
	  AND (
		l.seller_id = $1
		OR EXISTS (
			SELECT 1 FROM activity_events AS e
			WHERE e.user_id = $1 AND e.listing_id = l.id AND e.occurred_at < $2
		)
		OR EXISTS (
			SELECT 1 FROM deals AS d
			WHERE d.buyer_id = $1 AND d.listing_id = l.id AND d.created_at < $2
		)
	  )
	ORDER BY l.id
`

const loadDealsQuery = `
	SELECT d.id, d.listing_id, d.buyer_id, d.created_at, d.completed_at,
		d.status, d.delivery_used, d.price_band
	FROM deals AS d
	JOIN listings AS l ON l.id = d.listing_id
	WHERE d.created_at < $2
	  AND (d.buyer_id = $1 OR l.seller_id = $1)
	ORDER BY d.created_at, d.id
`

const loadReviewsQuery = `
	SELECT id, deal_id, author_id, recipient_id, rating, created_at
	FROM reviews
	WHERE created_at < $2
	  AND (author_id = $1 OR recipient_id = $1)
	ORDER BY created_at, id
`

const loadCategoriesQuery = `
	SELECT id, code, name, parent_id
	FROM categories
	WHERE id = ANY($1)
	ORDER BY id
`

const loadEngagementPercentilesQuery = `
	WITH engagement AS (
		SELECT
			l.id,
			l.category_id,
			COUNT(e.id)::DOUBLE PRECISION AS event_count
		FROM listings AS l
		LEFT JOIN activity_events AS e
			ON e.listing_id = l.id
			AND e.occurred_at >= $2
			AND e.occurred_at < $3
			AND e.event_type IN ('listing_view', 'contact')
		WHERE l.published_at < $3
		GROUP BY l.id, l.category_id
	), ranked AS (
		SELECT
			id,
			PERCENT_RANK() OVER (PARTITION BY category_id ORDER BY event_count) AS percentile
		FROM engagement
	)
	SELECT ranked.id, ranked.percentile
	FROM ranked
	JOIN listings AS l ON l.id = ranked.id
	WHERE l.seller_id = $1
	ORDER BY ranked.id
`

func (r *DatasetRepository) LoadDataset(
	ctx context.Context,
	user model.User,
	period model.Period,
	cutoff time.Time,
) (model.VersionedDataset, error) {
	const op = "postgres.DatasetRepository.LoadDataset"

	dataset := model.Dataset{User: user, Period: period, DataCutoffAt: cutoff}

	var err error
	if dataset.Events, err = r.loadEvents(ctx, user.ID, cutoff); err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}
	if dataset.Listings, err = r.loadListings(ctx, user.ID, cutoff); err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}
	if dataset.Deals, err = r.loadDeals(ctx, user.ID, cutoff); err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}
	if dataset.Reviews, err = r.loadReviews(ctx, user.ID, cutoff); err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}
	if dataset.Categories, err = r.loadCategories(ctx, dataset.Events, dataset.Listings); err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}
	if dataset.EngagementPercentiles, err = r.loadEngagementPercentiles(ctx, user.ID, period, cutoff); err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}

	encoded, err := json.Marshal(dataset)
	if err != nil {
		return model.VersionedDataset{}, fmt.Errorf("%s: %w", op, err)
	}
	return model.VersionedDataset{
		Value:   dataset,
		Version: fmt.Sprintf("sha256:%x", sha256.Sum256(encoded)),
	}, nil
}

func (r *DatasetRepository) loadEvents(ctx context.Context, userID uuid.UUID, cutoff time.Time) ([]model.ActivityEvent, error) {
	rows, err := r.pool.Query(ctx, loadEventsQuery, userID, cutoff)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.ActivityEvent])
}

func (r *DatasetRepository) loadListings(ctx context.Context, userID uuid.UUID, cutoff time.Time) ([]model.Listing, error) {
	rows, err := r.pool.Query(ctx, loadListingsQuery, userID, cutoff)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.Listing])
}

func (r *DatasetRepository) loadDeals(ctx context.Context, userID uuid.UUID, cutoff time.Time) ([]model.Deal, error) {
	rows, err := r.pool.Query(ctx, loadDealsQuery, userID, cutoff)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.Deal])
}

func (r *DatasetRepository) loadReviews(ctx context.Context, userID uuid.UUID, cutoff time.Time) ([]model.Review, error) {
	rows, err := r.pool.Query(ctx, loadReviewsQuery, userID, cutoff)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.Review])
}

func (r *DatasetRepository) loadCategories(
	ctx context.Context,
	events []model.ActivityEvent,
	listings []model.Listing,
) ([]model.Category, error) {
	categorySet := make(map[uuid.UUID]struct{})
	for _, event := range events {
		if event.CategoryID != nil {
			categorySet[*event.CategoryID] = struct{}{}
		}
	}
	for _, listing := range listings {
		categorySet[listing.CategoryID] = struct{}{}
	}
	if len(categorySet) == 0 {
		return []model.Category{}, nil
	}

	categoryIDs := make([]uuid.UUID, 0, len(categorySet))
	for categoryID := range categorySet {
		categoryIDs = append(categoryIDs, categoryID)
	}
	rows, err := r.pool.Query(ctx, loadCategoriesQuery, categoryIDs)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.Category])
}

type engagementPercentileRow struct {
	ListingID  uuid.UUID
	Percentile float64
}

func (r *DatasetRepository) loadEngagementPercentiles(
	ctx context.Context,
	userID uuid.UUID,
	period model.Period,
	cutoff time.Time,
) (map[uuid.UUID]float64, error) {
	rows, err := r.pool.Query(ctx, loadEngagementPercentilesQuery, userID, period.Start, cutoff)
	if err != nil {
		return nil, err
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByPos[engagementPercentileRow])
	if err != nil {
		return nil, err
	}

	percentiles := make(map[uuid.UUID]float64, len(values))
	for _, value := range values {
		percentiles[value.ListingID] = value.Percentile
	}
	return percentiles, nil
}
