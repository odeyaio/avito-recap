package engine

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"avito-recap/internal/model"

	"github.com/google/uuid"
)

type activityPoint struct {
	at         time.Time
	eventType  string
	listingID  *uuid.UUID
	categoryID *uuid.UUID
}

func (e *Engine) CalculateMetrics(dataset Dataset) (Metrics, error) {
	if err := dataset.Period.Validate(); err != nil {
		return Metrics{}, fmt.Errorf("validate period: %w", err)
	}

	location := dataset.Period.Start.Location()
	if dataset.User.Timezone != "" {
		var err error
		location, err = time.LoadLocation(dataset.User.Timezone)
		if err != nil {
			return Metrics{}, fmt.Errorf("load user timezone %q: %w", dataset.User.Timezone, err)
		}
	}
	cutoff := dataset.DataCutoffAt
	if cutoff.IsZero() || cutoff.After(dataset.Period.End) {
		cutoff = dataset.Period.End
	}

	listings := make(map[uuid.UUID]model.Listing, len(dataset.Listings))
	for _, listing := range dataset.Listings {
		listings[listing.ID] = listing
	}

	periodEvents := make([]model.ActivityEvent, 0, len(dataset.Events))
	userEvents := make([]model.ActivityEvent, 0, len(dataset.Events))
	for _, event := range dataset.Events {
		if dataset.Period.Contains(event.OccurredAt) && event.OccurredAt.Before(cutoff) {
			periodEvents = append(periodEvents, event)
			if event.UserID == dataset.User.ID {
				userEvents = append(userEvents, event)
			}
		}
	}
	slices.SortFunc(periodEvents, func(left, right model.ActivityEvent) int {
		return left.OccurredAt.Compare(right.OccurredAt)
	})
	slices.SortFunc(userEvents, func(left, right model.ActivityEvent) int {
		return left.OccurredAt.Compare(right.OccurredAt)
	})

	points := buildActivityPoints(dataset, userEvents, listings, cutoff)
	metrics := NewMetrics()
	e.calculateActivity(metrics, dataset, points, userEvents, listings, location, cutoff)
	e.calculateInterests(metrics, dataset, points, userEvents, listings, location)
	e.calculateIntent(metrics, dataset, userEvents, listings, cutoff)
	e.calculateMarketplace(metrics, dataset, periodEvents, listings, cutoff)
	e.calculateCommunity(metrics, dataset, cutoff)
	e.calculateFeatures(metrics, userEvents)

	return metrics, nil
}

func buildActivityPoints(
	dataset Dataset,
	periodEvents []model.ActivityEvent,
	listings map[uuid.UUID]model.Listing,
	cutoff time.Time,
) []activityPoint {
	points := make([]activityPoint, 0, len(periodEvents)+len(dataset.Listings)*2+len(dataset.Deals)+len(dataset.Reviews))
	for _, event := range periodEvents {
		points = append(points, activityPoint{
			at:         event.OccurredAt,
			eventType:  normalizeEventType(event.Type),
			listingID:  event.ListingID,
			categoryID: eventCategory(event, listings),
		})
	}

	for _, listing := range dataset.Listings {
		if listing.SellerID != dataset.User.ID {
			continue
		}
		categoryID := listing.CategoryID
		listingID := listing.ID
		if dataset.Period.Contains(listing.PublishedAt) && listing.PublishedAt.Before(cutoff) {
			points = append(points, activityPoint{listing.PublishedAt, "listing_published", &listingID, &categoryID})
		}
		if listing.ClosedAt != nil && dataset.Period.Contains(*listing.ClosedAt) && listing.ClosedAt.Before(cutoff) {
			points = append(points, activityPoint{*listing.ClosedAt, "listing_closed", &listingID, &categoryID})
		}
	}

	for _, deal := range dataset.Deals {
		if deal.Status != model.DealCompleted || deal.CompletedAt == nil ||
			!dataset.Period.Contains(*deal.CompletedAt) || !deal.CompletedAt.Before(cutoff) {
			continue
		}
		listing, exists := listings[deal.ListingID]
		if deal.BuyerID != dataset.User.ID && (!exists || listing.SellerID != dataset.User.ID) {
			continue
		}
		listingID := deal.ListingID
		var categoryID *uuid.UUID
		if exists {
			value := listing.CategoryID
			categoryID = &value
		}
		points = append(points, activityPoint{*deal.CompletedAt, "deal_completed", &listingID, categoryID})
	}

	for _, review := range dataset.Reviews {
		if review.AuthorID == dataset.User.ID && dataset.Period.Contains(review.CreatedAt) && review.CreatedAt.Before(cutoff) {
			points = append(points, activityPoint{at: review.CreatedAt, eventType: "review_left"})
		}
	}

	slices.SortFunc(points, func(left, right activityPoint) int { return left.at.Compare(right.at) })
	return points
}

func (e *Engine) calculateActivity(
	metrics Metrics,
	dataset Dataset,
	points []activityPoint,
	events []model.ActivityEvent,
	listings map[uuid.UUID]model.Listing,
	location *time.Location,
	cutoff time.Time,
) {
	views := int64(0)
	searches := int64(0)
	uniqueListings := make(map[uuid.UUID]struct{})
	freshViews := int64(0)
	days := make(map[string]time.Time)
	months := make(map[string]int64)
	hours := make(map[int]int64)
	quarters := [4]int64{}

	for _, point := range points {
		local := point.at.In(location)
		dayKey := local.Format("2006-01-02")
		days[dayKey] = time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
		months[local.Format("2006-01")]++
		hours[local.Hour()]++
		quarters[(int(local.Month())-1)/3]++
	}

	for _, event := range events {
		eventType := normalizeEventType(event.Type)
		switch eventType {
		case "listing_view":
			views++
			if event.ListingID != nil {
				uniqueListings[*event.ListingID] = struct{}{}
				if listing, exists := listings[*event.ListingID]; exists {
					age := event.OccurredAt.Sub(listing.PublishedAt)
					if age >= 0 && age <= e.config.FreshListingWindow {
						freshViews++
					}
				}
			}
		case "search":
			searches++
		}
	}

	month, mostActiveMonthActions := maxStringCount(months)
	favoriteHour, _ := maxIntCount(hours)
	peakQuarterActions := int64(0)
	for _, count := range quarters {
		if count > peakQuarterActions {
			peakQuarterActions = count
		}
	}
	totalActions := int64(len(points))

	metrics.Set("activity.views", views)
	metrics.Set("activity.unique_listings_viewed", int64(len(uniqueListings)))
	metrics.Set("activity.searches", searches)
	metrics.Set("activity.active_days", int64(len(days)))
	metrics.Set("activity.active_months", int64(len(months)))
	metrics.Set("activity.longest_active_streak_days", longestDayStreak(days))
	metrics.Set("activity.most_active_month_actions", mostActiveMonthActions)
	metrics.Set("activity.most_active_month", month)
	metrics.Set("activity.favorite_hour", int64(favoriteHour))
	metrics.Set("activity.fresh_listing_view_share", ratio(freshViews, views))
	metrics.Set("activity.total_actions", totalActions)
	metrics.Set("activity.peak_quarter_share", ratio(peakQuarterActions, totalActions))
	metrics.Set(
		"activity.returned_after_long_gap",
		returnedAfterLongGap(dataset, points, listings, cutoff, e.config.LongGap, location),
	)
}

func (e *Engine) calculateInterests(
	metrics Metrics,
	dataset Dataset,
	points []activityPoint,
	events []model.ActivityEvent,
	listings map[uuid.UUID]model.Listing,
	location *time.Location,
) {
	categoryActions := make(map[uuid.UUID]int64)
	categoryMonths := make(map[uuid.UUID]map[string]struct{})
	for _, point := range points {
		if point.categoryID == nil {
			continue
		}
		categoryActions[*point.categoryID]++
		if categoryMonths[*point.categoryID] == nil {
			categoryMonths[*point.categoryID] = make(map[string]struct{})
		}
		categoryMonths[*point.categoryID][point.at.In(location).Format("2006-01")] = struct{}{}
	}

	topCategoryActions := int64(0)
	consistentMonths := int64(0)
	for categoryID, actions := range categoryActions {
		if actions > topCategoryActions {
			topCategoryActions = actions
		}
		if count := int64(len(categoryMonths[categoryID])); count > consistentMonths {
			consistentMonths = count
		}
	}

	priorCategories := make(map[uuid.UUID]struct{})
	for _, event := range dataset.Events {
		if event.UserID != dataset.User.ID || !event.OccurredAt.Before(dataset.Period.Start) {
			continue
		}
		if categoryID := eventCategory(event, listings); categoryID != nil {
			priorCategories[*categoryID] = struct{}{}
		}
	}
	newSignificant := int64(0)
	for categoryID, actions := range categoryActions {
		_, seenBefore := priorCategories[categoryID]
		if !seenBefore && actions >= int64(e.config.SignificantCategoryEvents) {
			newSignificant++
		}
	}

	topicCounts := make(map[string]int64)
	lowSupplyActions := int64(0)
	for _, event := range events {
		if normalizeEventType(event.Type) != "search" {
			continue
		}
		if event.TopicKey != nil && *event.TopicKey != "" {
			topicCounts[*event.TopicKey]++
		}
		if event.ResultCount != nil && *event.ResultCount <= e.config.LowSupplyResultCount {
			lowSupplyActions++
		}
	}
	_, topTopicActions := maxStringCount(topicCounts)
	totalTopics := sumCounts(topicCounts)

	metrics.Set("interests.distinct_categories", int64(len(categoryActions)))
	metrics.Set("interests.top_category_actions", topCategoryActions)
	metrics.Set("interests.top_category_share", ratio(topCategoryActions, sumCounts(categoryActions)))
	metrics.Set("interests.new_significant_categories", newSignificant)
	metrics.Set("interests.most_consistent_category_months", consistentMonths)
	metrics.Set("interests.low_supply_actions", lowSupplyActions)
	metrics.Set("interests.top_search_topic_share", ratio(topTopicActions, totalTopics))
}

func (e *Engine) calculateIntent(
	metrics Metrics,
	dataset Dataset,
	events []model.ActivityEvent,
	listings map[uuid.UUID]model.Listing,
	cutoff time.Time,
) {
	viewTimes := make(map[uuid.UUID][]time.Time)
	contactTimes := make(map[uuid.UUID][]time.Time)
	periodFavorites := make(map[uuid.UUID]struct{})
	favoritesAdded := int64(0)
	topicSearches := make(map[string]int64)
	type favoriteState struct {
		active bool
		at     time.Time
	}
	activeFavorites := make(map[uuid.UUID]favoriteState)

	for _, event := range dataset.Events {
		if event.UserID != dataset.User.ID || !event.OccurredAt.Before(cutoff) || event.ListingID == nil {
			continue
		}
		switch favoriteAction(event) {
		case "add":
			state := activeFavorites[*event.ListingID]
			if state.at.IsZero() || event.OccurredAt.After(state.at) {
				activeFavorites[*event.ListingID] = favoriteState{active: true, at: event.OccurredAt}
			}
		case "remove":
			state := activeFavorites[*event.ListingID]
			if state.at.IsZero() || event.OccurredAt.After(state.at) {
				activeFavorites[*event.ListingID] = favoriteState{at: event.OccurredAt}
			}
		}
	}

	for _, event := range events {
		typeName := normalizeEventType(event.Type)
		if typeName == "search" {
			topic := ""
			if event.TopicKey != nil {
				topic = *event.TopicKey
			}
			if topic != "" {
				topicSearches[topic]++
			}
		}
		if event.ListingID == nil {
			continue
		}
		switch typeName {
		case "listing_view":
			viewTimes[*event.ListingID] = append(viewTimes[*event.ListingID], event.OccurredAt)
		case "contact":
			contactTimes[*event.ListingID] = append(contactTimes[*event.ListingID], event.OccurredAt)
		}
		if favoriteAction(event) == "add" {
			favoritesAdded++
			periodFavorites[*event.ListingID] = struct{}{}
		}
	}

	repeatViewedListings := int64(0)
	for _, times := range viewTimes {
		if len(times) >= 2 {
			repeatViewedListings++
		}
	}

	contactedListings := int64(len(contactTimes))
	carefulPaths := int64(0)
	fastPaths := int64(0)
	maxViewsBeforeContact := int64(0)
	for listingID, contacts := range contactTimes {
		slices.SortFunc(contacts, time.Time.Compare)
		firstContact := contacts[0]
		viewsBefore := int64(0)
		var firstView time.Time
		for _, viewedAt := range viewTimes[listingID] {
			if viewedAt.After(firstContact) {
				continue
			}
			viewsBefore++
			if firstView.IsZero() || viewedAt.Before(firstView) {
				firstView = viewedAt
			}
		}
		if viewsBefore >= 2 {
			carefulPaths++
		}
		if !firstView.IsZero() && firstContact.Sub(firstView) <= e.config.FastContactWindow {
			fastPaths++
		}
		if viewsBefore > maxViewsBeforeContact {
			maxViewsBeforeContact = viewsBefore
		}
	}

	completedDeals := int64(0)
	completedSearchPaths := int64(0)
	searchesByCategory := make(map[uuid.UUID][]time.Time)
	for _, event := range events {
		if normalizeEventType(event.Type) == "search" {
			if categoryID := eventCategory(event, listings); categoryID != nil {
				searchesByCategory[*categoryID] = append(searchesByCategory[*categoryID], event.OccurredAt)
			}
		}
	}
	for _, deal := range dataset.Deals {
		if deal.BuyerID != dataset.User.ID || deal.Status != model.DealCompleted || deal.CompletedAt == nil ||
			!dataset.Period.Contains(*deal.CompletedAt) || !deal.CompletedAt.Before(cutoff) {
			continue
		}
		completedDeals++
		listing, exists := listings[deal.ListingID]
		if !exists {
			continue
		}
		for _, searchedAt := range searchesByCategory[listing.CategoryID] {
			if searchedAt.Before(*deal.CompletedAt) {
				completedSearchPaths++
				break
			}
		}
	}

	repeatedSearches := int64(0)
	for _, count := range topicSearches {
		if count > 1 {
			repeatedSearches += count - 1
		}
	}
	contactedFavorites := int64(0)
	for listingID := range periodFavorites {
		if len(contactTimes[listingID]) > 0 {
			contactedFavorites++
		}
	}
	activeFavoriteCount := int64(0)
	for _, state := range activeFavorites {
		if state.active {
			activeFavoriteCount++
		}
	}

	metrics.Set("intent.repeat_viewed_listings", repeatViewedListings)
	metrics.Set("intent.favorites_added", favoritesAdded)
	metrics.Set("intent.active_favorites", activeFavoriteCount)
	metrics.Set("intent.contacts", int64(totalEventsOfType(events, "contact")))
	metrics.Set("intent.contacted_unique_listings", contactedListings)
	metrics.Set("intent.completed_deals", completedDeals)
	metrics.Set("intent.contact_to_deal_conversion", ratio(completedDeals, contactedListings))
	metrics.Set("intent.careful_contact_paths", carefulPaths)
	metrics.Set("intent.fast_view_to_contact_paths", fastPaths)
	metrics.Set("intent.repeated_searches", repeatedSearches)
	metrics.Set("intent.max_views_before_contact", maxViewsBeforeContact)
	metrics.Set("intent.favorite_to_contact_conversion", ratio(contactedFavorites, int64(len(periodFavorites))))
	metrics.Set("intent.completed_search_to_deal_paths", completedSearchPaths)
}

func (e *Engine) calculateMarketplace(
	metrics Metrics,
	dataset Dataset,
	events []model.ActivityEvent,
	listings map[uuid.UUID]model.Listing,
	cutoff time.Time,
) {
	purchases := int64(0)
	sales := int64(0)
	cancelledPurchases := int64(0)
	deliveryDeals := int64(0)
	deliveryRegions := make(map[string]struct{})
	completedBeforePeriod := false

	for _, deal := range dataset.Deals {
		listing, listingExists := listings[deal.ListingID]
		isBuyer := deal.BuyerID == dataset.User.ID
		isSeller := listingExists && listing.SellerID == dataset.User.ID
		if !isBuyer && !isSeller {
			continue
		}
		if deal.Status == model.DealCompleted && deal.CompletedAt != nil && deal.CompletedAt.Before(dataset.Period.Start) {
			completedBeforePeriod = true
		}
		if deal.Status == model.DealCancelled && isBuyer && dataset.Period.Contains(deal.CreatedAt) && deal.CreatedAt.Before(cutoff) {
			cancelledPurchases++
		}
		if deal.Status != model.DealCompleted || deal.CompletedAt == nil ||
			!dataset.Period.Contains(*deal.CompletedAt) || !deal.CompletedAt.Before(cutoff) {
			continue
		}
		if isBuyer {
			purchases++
		}
		if isSeller {
			sales++
		}
		if deal.DeliveryUsed {
			deliveryDeals++
			if listingExists && listing.Region != "" {
				deliveryRegions[listing.Region] = struct{}{}
			}
		}
	}

	publishedListings := int64(0)
	closedListings := int64(0)
	completeListings := int64(0)
	for _, listing := range dataset.Listings {
		if listing.SellerID != dataset.User.ID {
			continue
		}
		if dataset.Period.Contains(listing.PublishedAt) && listing.PublishedAt.Before(cutoff) {
			publishedListings++
			if listing.PhotoCount >= 3 && listing.DescriptionComplete {
				completeListings++
			}
		}
		if listing.ClosedAt != nil && dataset.Period.Contains(*listing.ClosedAt) && listing.ClosedAt.Before(cutoff) {
			closedListings++
		}
	}

	listingViews := int64(0)
	listingContacts := int64(0)
	listingEdits := int64(0)
	deliveryEnabled := make(map[uuid.UUID]struct{})
	fastResponses := make(map[uuid.UUID]struct{})
	for _, event := range events {
		if event.ListingID == nil {
			continue
		}
		listing, exists := listings[*event.ListingID]
		if !exists || listing.SellerID != dataset.User.ID {
			continue
		}
		switch normalizeEventType(event.Type) {
		case "listing_view":
			listingViews++
		case "contact":
			listingContacts++
			age := event.OccurredAt.Sub(listing.PublishedAt)
			if age >= 0 && age <= e.config.FastListingResponseWindow {
				fastResponses[listing.ID] = struct{}{}
			}
		case "listing_edit":
			listingEdits++
		case "delivery_enable":
			deliveryEnabled[listing.ID] = struct{}{}
		}
	}

	totalCompleted := purchases + sales
	metrics.Set("marketplace.purchases", purchases)
	metrics.Set("marketplace.sales", sales)
	metrics.Set("marketplace.cancelled_purchases", cancelledPurchases)
	metrics.Set("marketplace.delivery_deals", deliveryDeals)
	metrics.Set("marketplace.delivery_share", ratio(deliveryDeals, totalCompleted))
	metrics.Set("marketplace.delivery_regions", int64(len(deliveryRegions)))
	metrics.Set("marketplace.published_listings", publishedListings)
	metrics.Set("marketplace.closed_listings", closedListings)
	metrics.Set("marketplace.listing_views", listingViews)
	metrics.Set("marketplace.listing_contacts", listingContacts)
	metrics.Set("marketplace.sale_completion_rate", ratio(sales, publishedListings))
	metrics.Set("marketplace.first_completed_deal_in_period", totalCompleted > 0 && !completedBeforePeriod)
	metrics.Set("marketplace.fast_listing_responses", int64(len(fastResponses)))
	metrics.Set("marketplace.complete_listings", completeListings)
	metrics.Set("marketplace.listing_edits", listingEdits)
	metrics.Set("marketplace.delivery_enabled_listings", int64(len(deliveryEnabled)))

	bestPercentile := -1.0
	for listingID, percentile := range dataset.EngagementPercentiles {
		listing, exists := listings[listingID]
		if exists && listing.SellerID == dataset.User.ID && percentile > bestPercentile {
			bestPercentile = percentile
		}
	}
	if bestPercentile >= 0 {
		metrics.Set("marketplace.best_listing_engagement_percentile", clamp01(bestPercentile))
	}
}

func (e *Engine) calculateCommunity(metrics Metrics, dataset Dataset, cutoff time.Time) {
	reviewsLeft := int64(0)
	reviewsReceived := int64(0)
	fiveStarRatings := int64(0)
	ratingSum := int64(0)
	reviewsAfterDeal := int64(0)
	deals := make(map[uuid.UUID]model.Deal, len(dataset.Deals))
	for _, deal := range dataset.Deals {
		deals[deal.ID] = deal
	}

	for _, review := range dataset.Reviews {
		if !dataset.Period.Contains(review.CreatedAt) || !review.CreatedAt.Before(cutoff) {
			continue
		}
		if review.AuthorID == dataset.User.ID {
			reviewsLeft++
			if review.DealID != nil {
				deal, exists := deals[*review.DealID]
				if exists && deal.CompletedAt != nil && !review.CreatedAt.Before(*deal.CompletedAt) {
					reviewsAfterDeal++
				}
			}
		}
		if review.RecipientID == dataset.User.ID {
			reviewsReceived++
			ratingSum += int64(review.Rating)
			if review.Rating == 5 {
				fiveStarRatings++
			}
		}
	}

	metrics.Set("community.reviews_left", reviewsLeft)
	metrics.Set("community.reviews_received", reviewsReceived)
	metrics.Set("community.five_star_ratings", fiveStarRatings)
	metrics.Set("community.review_after_deal_share", ratio(reviewsAfterDeal, reviewsLeft))
	if reviewsReceived > 0 {
		metrics.Set("community.average_rating", float64(ratingSum)/float64(reviewsReceived))
	}
}

func (e *Engine) calculateFeatures(metrics Metrics, events []model.ActivityEvent) {
	notificationOpens := int64(0)
	promotionUses := int64(0)
	searchesWithFilters := int64(0)
	for _, event := range events {
		switch normalizeEventType(event.Type) {
		case "notification_open":
			notificationOpens++
		case "promotion_use":
			promotionUses++
		case "search":
			if event.FilterCount != nil && *event.FilterCount > 0 {
				searchesWithFilters++
			}
		}
		if event.Source != nil && strings.EqualFold(string(*event.Source), "notification") {
			notificationOpens++
		}
	}
	metrics.Set("features.notification_opens", notificationOpens)
	metrics.Set("features.promotion_uses", promotionUses)
	metrics.Set("features.searches_with_filters", searchesWithFilters)
}

func normalizeEventType(eventType model.EventType) string {
	value := strings.ToLower(string(eventType))
	switch value {
	case "view":
		return "listing_view"
	case "favorite", "favorite_added":
		return "favorite_add"
	case "favorite_removed":
		return "favorite_remove"
	case "seller_contact", "listing_contact":
		return "contact"
	case "notification":
		return "notification_open"
	case "promotion":
		return "promotion_use"
	default:
		return value
	}
}

func favoriteAction(event model.ActivityEvent) string {
	switch normalizeEventType(event.Type) {
	case "favorite_add":
		return "add"
	case "favorite_remove":
		return "remove"
	}
	if len(event.Properties) == 0 {
		return ""
	}
	var properties struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(event.Properties, &properties); err != nil {
		return ""
	}
	return strings.ToLower(properties.Action)
}

func eventCategory(event model.ActivityEvent, listings map[uuid.UUID]model.Listing) *uuid.UUID {
	if event.CategoryID != nil {
		value := *event.CategoryID
		return &value
	}
	if event.ListingID != nil {
		if listing, exists := listings[*event.ListingID]; exists {
			value := listing.CategoryID
			return &value
		}
	}
	return nil
}

func totalEventsOfType(events []model.ActivityEvent, eventType string) int {
	count := 0
	for _, event := range events {
		if normalizeEventType(event.Type) == eventType {
			count++
		}
	}
	return count
}

func ratio(numerator, denominator int64) float64 {
	if denominator <= 0 {
		return 0
	}
	return clamp01(float64(numerator) / float64(denominator))
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func sumCounts[K comparable](values map[K]int64) int64 {
	total := int64(0)
	for _, value := range values {
		total += value
	}
	return total
}

func maxStringCount(values map[string]int64) (string, int64) {
	bestKey := ""
	bestValue := int64(0)
	for key, value := range values {
		if value > bestValue || (value == bestValue && (bestKey == "" || key < bestKey)) {
			bestKey = key
			bestValue = value
		}
	}
	return bestKey, bestValue
}

func maxIntCount(values map[int]int64) (int, int64) {
	bestKey := 0
	bestValue := int64(0)
	for key, value := range values {
		if value > bestValue || (value == bestValue && key < bestKey) {
			bestKey = key
			bestValue = value
		}
	}
	return bestKey, bestValue
}

func longestDayStreak(days map[string]time.Time) int64 {
	ordered := make([]time.Time, 0, len(days))
	for _, day := range days {
		ordered = append(ordered, day)
	}
	slices.SortFunc(ordered, time.Time.Compare)
	longest := int64(0)
	current := int64(0)
	var previous time.Time
	for _, day := range ordered {
		if !previous.IsZero() && day.AddDate(0, 0, -1).Equal(previous) {
			current++
		} else {
			current = 1
		}
		if current > longest {
			longest = current
		}
		previous = day
	}
	return longest
}

func hasLongGap(points []activityPoint, gap time.Duration, location *time.Location) bool {
	if len(points) < 2 {
		return false
	}
	previous := points[0].at.In(location)
	for _, point := range points[1:] {
		current := point.at.In(location)
		if current.Sub(previous) >= gap {
			return true
		}
		previous = current
	}
	return false
}

func returnedAfterLongGap(
	dataset Dataset,
	periodPoints []activityPoint,
	listings map[uuid.UUID]model.Listing,
	cutoff time.Time,
	gap time.Duration,
	location *time.Location,
) bool {
	if hasLongGap(periodPoints, gap, location) {
		return true
	}
	if len(periodPoints) == 0 {
		return false
	}

	var previous time.Time
	consider := func(at time.Time) {
		if at.Before(dataset.Period.Start) && at.Before(cutoff) && (previous.IsZero() || at.After(previous)) {
			previous = at
		}
	}
	for _, event := range dataset.Events {
		if event.UserID == dataset.User.ID {
			consider(event.OccurredAt)
		}
	}
	for _, listing := range dataset.Listings {
		if listing.SellerID != dataset.User.ID {
			continue
		}
		consider(listing.PublishedAt)
		if listing.ClosedAt != nil {
			consider(*listing.ClosedAt)
		}
	}
	for _, deal := range dataset.Deals {
		listing, exists := listings[deal.ListingID]
		if deal.Status == model.DealCompleted && deal.CompletedAt != nil &&
			(deal.BuyerID == dataset.User.ID || (exists && listing.SellerID == dataset.User.ID)) {
			consider(*deal.CompletedAt)
		}
	}
	for _, review := range dataset.Reviews {
		if review.AuthorID == dataset.User.ID {
			consider(review.CreatedAt)
		}
	}

	return !previous.IsZero() && periodPoints[0].at.Sub(previous) >= gap
}
