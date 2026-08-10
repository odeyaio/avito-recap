package engine

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"avito-recap/internal/catalog"
	"avito-recap/internal/model"

	"github.com/google/uuid"
)

func TestEngineRunCalculatesMetricsAndClassifies(t *testing.T) {
	t.Parallel()

	dataset := fixtureDataset(t)
	behaviorCatalog := catalog.BehaviorCatalog{
		Version: "test-v1",
		Behaviors: []catalog.BehaviorDefinition{
			behaviorDefinition("secondary", 20, "activity.views", "gte", 2),
			behaviorDefinition("primary", 10, "intent.careful_contact_paths", "gte", 1),
			behaviorDefinition("not_matched", 5, "activity.views", "gt", 100),
		},
	}
	achievementCatalog := catalog.AchievementCatalog{
		Version: "test-v1",
		Achievements: []catalog.AchievementDefinition{
			achievementDefinition("contacted", 20, "intent.contacted_unique_listings", "gte", 1),
			achievementDefinition("fresh", 10, "activity.fresh_listing_view_share", "gte", 0.5),
		},
	}

	result, err := newTestEngine(t).Run(dataset, achievementCatalog, behaviorCatalog)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	assertMetric(t, result.Metrics, "activity.views", int64(3))
	assertMetric(t, result.Metrics, "activity.longest_active_streak_days", int64(5))
	assertMetric(t, result.Metrics, "activity.returned_after_long_gap", true)
	assertMetric(t, result.Metrics, "interests.distinct_categories", int64(2))
	assertMetric(t, result.Metrics, "intent.repeat_viewed_listings", int64(1))
	assertMetric(t, result.Metrics, "intent.active_favorites", int64(1))
	assertMetric(t, result.Metrics, "intent.careful_contact_paths", int64(1))
	assertMetric(t, result.Metrics, "intent.repeated_searches", int64(2))
	assertMetric(t, result.Metrics, "marketplace.purchases", int64(1))
	assertMetric(t, result.Metrics, "marketplace.sales", int64(1))
	assertMetric(t, result.Metrics, "marketplace.listing_views", int64(1))
	assertMetric(t, result.Metrics, "marketplace.listing_contacts", int64(1))
	assertMetric(t, result.Metrics, "community.reviews_left", int64(1))
	assertMetric(t, result.Metrics, "community.reviews_received", int64(1))
	assertMetric(t, result.Metrics, "features.searches_with_filters", int64(3))

	if len(result.Behaviors) != 2 {
		t.Fatalf("len(Behaviors) = %d, want 2", len(result.Behaviors))
	}
	if result.Behaviors[0].Definition.Code != "primary" || !result.Behaviors[0].IsPrimary {
		t.Fatalf("first behavior = %#v, want primary", result.Behaviors[0])
	}
	if result.Behaviors[1].IsPrimary {
		t.Fatal("second behavior unexpectedly marked primary")
	}
	if len(result.Achievements) != 2 || result.Achievements[0].Definition.Code != "fresh" {
		t.Fatalf("Achievements = %#v, want fresh then contacted", result.Achievements)
	}

	encoded, err := json.Marshal(result.Metrics)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !strings.Contains(string(encoded), `"activity":{"active_days"`) {
		t.Fatalf("metrics JSON is not nested: %s", encoded)
	}
	if strings.Contains(string(encoded), `"activity.views"`) {
		t.Fatalf("metrics JSON contains a flat key: %s", encoded)
	}
	var decoded Metrics
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	decodedViews, exists := decoded.Get("activity.views")
	if !exists || decodedViews != json.Number("3") {
		t.Fatalf("decoded activity.views = %#v, %v; want json.Number(3)", decodedViews, exists)
	}
}

func TestEngineRunRejectsEmptyActivity(t *testing.T) {
	t.Parallel()

	period, err := PeriodForYear(2025, "UTC")
	if err != nil {
		t.Fatal(err)
	}
	dataset := Dataset{User: model.User{ID: uuid.New(), Timezone: "UTC"}, Period: period}

	_, err = newTestEngine(t).Run(dataset,
		catalog.AchievementCatalog{Version: "v1", Achievements: []catalog.AchievementDefinition{
			achievementDefinition("something", 10, "activity.views", "gte", 1),
		}},
		catalog.BehaviorCatalog{Version: "v1", Behaviors: []catalog.BehaviorDefinition{
			behaviorDefinition("something", 10, "activity.views", "gte", 1),
		}},
	)
	if !errors.Is(err, ErrNoActivity) {
		t.Fatalf("Run() error = %v, want ErrNoActivity", err)
	}
}

func newTestEngine(t *testing.T) *Engine {
	t.Helper()

	engine, err := New(Config{
		FreshListingWindow:        48 * time.Hour,
		LongGap:                   90 * 24 * time.Hour,
		FastContactWindow:         time.Hour,
		FastListingResponseWindow: 24 * time.Hour,
		SignificantCategoryEvents: 3,
		LowSupplyResultCount:      10,
		MinimumActions:            1,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return engine
}

func TestPeriodForYearUsesUserTimezone(t *testing.T) {
	t.Parallel()

	period, err := PeriodForYear(2025, "Asia/Omsk")
	if err != nil {
		t.Fatalf("PeriodForYear() error = %v", err)
	}
	if got := period.Start.Format(time.RFC3339); got != "2025-01-01T00:00:00+06:00" {
		t.Fatalf("Start = %s", got)
	}
	if got := period.End.Format(time.RFC3339); got != "2026-01-01T00:00:00+06:00" {
		t.Fatalf("End = %s", got)
	}
}

func fixtureDataset(t *testing.T) Dataset {
	t.Helper()

	location, err := time.LoadLocation("Asia/Omsk")
	if err != nil {
		t.Fatal(err)
	}
	period, err := PeriodForYear(2025, "Asia/Omsk")
	if err != nil {
		t.Fatal(err)
	}
	at := func(month time.Month, day, hour int) time.Time {
		return time.Date(2025, month, day, hour, 0, 0, 0, location)
	}

	userID := uuid.New()
	otherUserID := uuid.New()
	categoryA := uuid.New()
	categoryB := uuid.New()
	browsedListingID := uuid.New()
	soldListingID := uuid.New()
	browsedListing := model.Listing{
		ID: browsedListingID, SellerID: otherUserID, CategoryID: categoryA,
		Region: "Омск", PublishedAt: at(time.January, 1, 7),
	}
	soldListing := model.Listing{
		ID: soldListingID, SellerID: userID, CategoryID: categoryB,
		Region: "Томск", PublishedAt: at(time.February, 1, 10),
		DeliveryAvailable: true, PhotoCount: 5, DescriptionComplete: true,
	}
	filterCount := 2
	lowResults := 5
	topic := "phone"
	categoryAPointer := categoryA
	view := func(occurredAt time.Time) model.ActivityEvent {
		return model.ActivityEvent{
			ID: uuid.New(), UserID: userID, Type: model.EventListingView, OccurredAt: occurredAt,
			ListingID: &browsedListingID,
		}
	}
	search := func(occurredAt time.Time) model.ActivityEvent {
		return model.ActivityEvent{
			ID: uuid.New(), UserID: userID, Type: model.EventSearch, OccurredAt: occurredAt,
			CategoryID: &categoryAPointer, FilterCount: &filterCount, ResultCount: &lowResults, TopicKey: &topic,
		}
	}
	events := []model.ActivityEvent{
		{ID: uuid.New(), UserID: userID, Type: model.EventListingView, OccurredAt: time.Date(2024, time.September, 1, 10, 0, 0, 0, location), ListingID: &browsedListingID},
		view(at(time.January, 1, 8)), view(at(time.January, 1, 9)), view(at(time.January, 2, 8)),
		search(at(time.January, 1, 10)), search(at(time.January, 2, 10)), search(at(time.January, 3, 10)),
		{ID: uuid.New(), UserID: userID, Type: model.EventContact, OccurredAt: at(time.January, 2, 9), ListingID: &browsedListingID},
		// Intentionally unordered: the latest action is add, so this listing is active.
		{ID: uuid.New(), UserID: userID, Type: model.EventFavoriteAdd, OccurredAt: at(time.January, 5, 9), ListingID: &browsedListingID},
		{ID: uuid.New(), UserID: userID, Type: model.EventFavoriteRemove, OccurredAt: at(time.January, 4, 9), ListingID: &browsedListingID},
		// Other users' events must count towards this user's seller metrics only.
		{ID: uuid.New(), UserID: otherUserID, Type: model.EventListingView, OccurredAt: at(time.February, 1, 11), ListingID: &soldListingID},
		{ID: uuid.New(), UserID: otherUserID, Type: model.EventContact, OccurredAt: at(time.February, 1, 12), ListingID: &soldListingID},
	}
	purchaseCompletedAt := at(time.March, 1, 10)
	saleCompletedAt := at(time.April, 1, 10)
	purchaseDealID := uuid.New()
	saleDealID := uuid.New()
	deals := []model.Deal{
		{ID: purchaseDealID, ListingID: browsedListingID, BuyerID: userID, CreatedAt: at(time.February, 20, 10), CompletedAt: &purchaseCompletedAt, Status: model.DealCompleted},
		{ID: saleDealID, ListingID: soldListingID, BuyerID: otherUserID, CreatedAt: at(time.March, 20, 10), CompletedAt: &saleCompletedAt, Status: model.DealCompleted, DeliveryUsed: true},
	}
	reviews := []model.Review{
		{ID: uuid.New(), DealID: &purchaseDealID, AuthorID: userID, RecipientID: otherUserID, Rating: 5, CreatedAt: at(time.March, 2, 10)},
		{ID: uuid.New(), DealID: &saleDealID, AuthorID: otherUserID, RecipientID: userID, Rating: 5, CreatedAt: at(time.April, 2, 10)},
	}

	return Dataset{
		User: model.User{ID: userID, Timezone: "UTC"}, Period: period,
		DataCutoffAt: period.End,
		Events:       events, Listings: []model.Listing{browsedListing, soldListing}, Deals: deals, Reviews: reviews,
		Categories: []model.Category{{ID: categoryA, Code: "phones"}, {ID: categoryB, Code: "cars"}},
	}
}

func behaviorDefinition(code string, order int, metric, operator string, value any) catalog.BehaviorDefinition {
	return catalog.BehaviorDefinition{
		Code: code, Name: code, Description: code, SortOrder: order,
		Rule: map[string]any{"metric": metric, "operator": operator, "value": value},
		DefaultAction: catalog.DefaultAction{
			Code: "open_something", Title: "Open", TargetType: "search", Href: "https://www.avito.ru/rossiya",
		},
	}
}

func achievementDefinition(code string, order int, metric, operator string, value any) catalog.AchievementDefinition {
	return catalog.AchievementDefinition{
		Code: code, Name: code, Description: code, IconKey: "achievements/" + code + ".webp", SortOrder: order,
		Rule: map[string]any{"metric": metric, "operator": operator, "value": value},
	}
}

func assertMetric(t *testing.T, metrics Metrics, code string, want any) {
	t.Helper()
	got, exists := metrics.Get(code)
	if !exists {
		t.Fatalf("metric %q does not exist", code)
	}
	if got != want {
		t.Fatalf("metric %q = %#v, want %#v", code, got, want)
	}
}
