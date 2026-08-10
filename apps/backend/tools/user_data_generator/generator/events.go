package generator

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"avito-recap/internal/model"
)

func GenerateUserEvents(seed int64, users []model.User, config []UserConfig, listings []model.Listing, categories []CategoryConfig) ([]model.ActivityEvent, []model.Deal, []model.Review) {
	source := rand.NewSource(seed)
	rnd := rand.New(source)

	events := make([]model.ActivityEvent, 0)
	deals := make([]model.Deal, 0)
	reviews := make([]model.Review, 0)

	if len(users) == 0 || len(listings) == 0 {
		return events, deals, reviews
	}

	categoryMap := make(map[string]uuid.UUID)
	for _, cat := range categories {
		categoryMap[cat.Name] = CategoryID(cat.Name)
	}

	usersWithGen := make([]UserWithGenData, len(users))
	for idx, user := range users {
		var cfg UserConfig
		var preferred, unlikely []string
		var preferredSell, unlikelySell []string

		if idx < len(config) {
			cfg = config[idx]
			preferred = cfg.PreferredCategories
			unlikely = cfg.UnlikelyCategories
			preferredSell = cfg.PreferredSellCategories
			unlikelySell = cfg.UnlikelySellCategories
		} else {
			cfg = RandomUserConfig(rnd, categories)
			preferred = cfg.PreferredCategories
			unlikely = cfg.UnlikelyCategories
			preferredSell = cfg.PreferredSellCategories
			unlikelySell = cfg.UnlikelySellCategories
		}

		usersWithGen[idx] = UserWithGenData{
			User:                    user,
			search:                  cfg.Search,
			view:                    cfg.View,
			favorite:                cfg.Favorite,
			contact:                 cfg.Contact,
			deal:                    cfg.Deal,
			review:                  cfg.Review,
			notificationOpen:        cfg.NotificationOpen,
			PreferredCategories:     preferred,
			UnlikelyCategories:      unlikely,
			PriceSegment:            cfg.PriceSegment,
			SellFrequency:           cfg.SellFrequency,
			PreferredSellCategories: preferredSell,
			UnlikelySellCategories:  unlikelySell,
		}
	}

	listingsWithGen := make([]ListingWithGenData, len(listings))
	for idx, listing := range listings {
		categoryName := ""
		for name, id := range categoryMap {
			if id == listing.CategoryID {
				categoryName = name
				break
			}
		}
		price := 50000.0
		switch listing.PriceBand {
		case "budget":
			price = 15000.0
		case "premium":
			price = 200000.0
		}
		listingsWithGen[idx] = ListingWithGenData{
			Listing:  listing,
			Category: categoryName,
			Price:    price,
		}
	}

	generatorConfig := DefaultGeneratorConfig()

	soldListings := make(map[uuid.UUID]bool)
	cancelledListingsPerUser := make(map[uuid.UUID]map[uuid.UUID]bool) // userID -> listingID -> cancelled

	for _, user := range usersWithGen {
		if cancelledListingsPerUser[user.ID] == nil {
			cancelledListingsPerUser[user.ID] = make(map[uuid.UUID]bool)
		}
		cfg := user.GetConfig()

		baseTime := user.RegisteredAt
		if baseTime.IsZero() {
			baseTime = time.Now().UTC().Add(-time.Duration(generatorConfig.EventBaseTimeOffsetHours) * time.Hour)
		}
		now := time.Now().UTC()
		recentWindowStart := now.AddDate(0, -generatorConfig.RecentActivityWindowMonths, 0)
		if recentWindowStart.Before(baseTime) {
			recentWindowStart = baseTime
		}

		intentCount := cfg.IntentCount
		if intentCount <= 0 {
			intentCount = calculateIntentCount(user.PreferredCategories, categories, generatorConfig, rnd)
		}

		for intentIdx := 0; intentIdx < intentCount; intentIdx++ {
			// Select category for this intent
			categoryName := selectIntentCategoryForIndex(cfg, categories, intentIdx)
			if categoryName == "" {
				categoryName = generatorConfig.EventDefaultCategory
				if len(categories) > 0 {
					categoryName = categories[rnd.Intn(len(categories))].Name
				}
			}

			var intentBaseTime time.Time
			if rnd.Float64() < generatorConfig.RecentActivityBias {
				intentBaseTime = randomDateBetweenRange(recentWindowStart, now, rnd)
			} else {
				intentBaseTime = randomDateBetweenRange(baseTime, now, rnd)
			}

			// Get initial candidates from search
			maxPrice := estimateIntentPrice(cfg, categoryName, categories)
			filtered := filterListingsForIntent(listingsWithGen, user, categoryName, maxPrice, cfg)
			initialCandidates := selectSearchCandidates(rnd, filtered, user, cfg)

			categoryUUID := CategoryID(categoryName)
			topic := categoryName
			resultCount := len(filtered)
			filterCount := 0
			if rnd.Float64() < 0.7 {
				filterCount = 1 + rnd.Intn(2)
			}
			properties, _ := json.Marshal(map[string]interface{}{
				"topic":   categoryName,
				"filters": map[string]interface{}{"category": categoryName},
			})
			events = append(events, model.ActivityEvent{
				ID:          seededUUID(rnd),
				UserID:      user.ID,
				Type:        (model.EventTypeSearch),
				OccurredAt:  intentBaseTime,
				CategoryID:  &categoryUUID,
				TopicKey:    &topic,
				ResultCount: &resultCount,
				FilterCount: &filterCount,
				Properties:  properties,
				IngestedAt:  time.Now().UTC(),
			})

			// VIEW stage - filter through view probability
			viewedListings := make([]ListingWithGenData, 0)
			for _, listing := range initialCandidates {
				if rnd.Float64() < clamp(cfg.View, 0, 1) {
					viewedListings = append(viewedListings, listing)

					listingID := listing.ID
					categoryID := listing.CategoryID

					// Determine if this view starts from a notification
					var source *model.EventSource
					if rnd.Float64() < clamp(cfg.NotificationOpen, 0, 1) {
						notificationSource := model.EventSource(model.EventSourceNotification)
						source = &notificationSource
					}

					viewTime := randomTimeOffset(intentBaseTime, generatorConfig.EventRandomSpreadSeconds, rnd)
					properties, _ := json.Marshal(map[string]interface{}{"price_band": listing.PriceBand})
					events = append(events, model.ActivityEvent{
						ID:         seededUUID(rnd),
						UserID:     user.ID,
						Type:       model.EventType(model.EventTypeView),
						OccurredAt: viewTime,
						ListingID:  &listingID,
						CategoryID: &categoryID,
						Source:     source,
						Properties: properties,
						IngestedAt: time.Now().UTC(),
					})

					// SHARE event - after view, independent probability
					if rnd.Float64() < clamp(generatorConfig.ShareProbability, 0, 1) {
						shareTime := randomTimeOffset(viewTime.Add(time.Duration(5)*time.Minute), generatorConfig.EventRandomSpreadSeconds, rnd)
						shareProperties, _ := json.Marshal(map[string]interface{}{"price_band": listing.PriceBand})
						events = append(events, model.ActivityEvent{
							ID:         seededUUID(rnd),
							UserID:     user.ID,
							Type:       model.EventType(model.EventTypeShare),
							OccurredAt: shareTime,
							ListingID:  &listingID,
							CategoryID: &categoryID,
							Properties: shareProperties,
							IngestedAt: time.Now().UTC(),
						})
					}
				}
			}

			// FAVORITE stage - only from viewed listings
			favoritedListings := make([]ListingWithGenData, 0)
			for _, listing := range viewedListings {
				if rnd.Float64() < clamp(cfg.Favorite, 0, 1) {
					favoritedListings = append(favoritedListings, listing)

					listingID := listing.ID
					categoryID := listing.CategoryID
					favoriteTime := randomTimeOffset(intentBaseTime.Add(time.Duration(generatorConfig.EventFavoriteOffsetMinutes)*time.Minute), generatorConfig.EventRandomSpreadSeconds, rnd)
					properties, _ := json.Marshal(map[string]interface{}{"price_band": listing.PriceBand})
					events = append(events, model.ActivityEvent{
						ID:         seededUUID(rnd),
						UserID:     user.ID,
						Type:       model.EventType(model.EventTypeFavorite),
						OccurredAt: favoriteTime,
						ListingID:  &listingID,
						CategoryID: &categoryID,
						Properties: properties,
						IngestedAt: time.Now().UTC(),
					})
				}
			}

			// CONTACT stage - only from viewed listings
			contactedListings := make([]ListingWithGenData, 0)
			cancelledListings := make(map[uuid.UUID]bool)
			for _, listing := range viewedListings {
				if rnd.Float64() < clamp(cfg.Contact, 0, 1) {
					listingID := listing.ID
					categoryID := listing.CategoryID
					contactTime := randomTimeOffset(intentBaseTime.Add(time.Duration(generatorConfig.EventContactOffsetMinutes)*time.Minute), generatorConfig.EventRandomSpreadSeconds, rnd)
					properties, _ := json.Marshal(map[string]interface{}{"price_band": listing.PriceBand})
					events = append(events, model.ActivityEvent{
						ID:         seededUUID(rnd),
						UserID:     user.ID,
						Type:       model.EventType(model.EventTypeContact),
						OccurredAt: contactTime,
						ListingID:  &listingID,
						CategoryID: &categoryID,
						Properties: properties,
						IngestedAt: time.Now().UTC(),
					})

					// CANCEL_DEAL event - after contact, before deal, doesn't lead to actual deal
					if rnd.Float64() < clamp(generatorConfig.CancelDealProbability, 0, 1) {
						cancelTime := randomTimeOffset(contactTime.Add(time.Duration(10)*time.Minute), generatorConfig.EventRandomSpreadSeconds, rnd)
						events = append(events, model.ActivityEvent{
							ID:         seededUUID(rnd),
							UserID:     user.ID,
							Type:       model.EventType(model.EventTypeCancelDeal),
							OccurredAt: cancelTime,
							ListingID:  &listingID,
							CategoryID: &categoryID,
							Properties: properties,
							IngestedAt: time.Now().UTC(),
						})
						cancelledListings[listing.ID] = true
						cancelledListingsPerUser[user.ID][listing.ID] = true
					} else {
						contactedListings = append(contactedListings, listing)
					}
				}
			}

			// DEAL stage - only from contacted listings (not cancelled) and not already sold
			for _, listing := range contactedListings {
				// Skip if listing is already cancelled by this user or already sold globally
				if cancelledListingsPerUser[user.ID][listing.ID] || soldListings[listing.ID] {
					continue
				}

				if rnd.Float64() < clamp(cfg.Deal, 0, 1) {
					dealTime := intentBaseTime.Add(time.Duration(generatorConfig.EventDealOffsetMinutes) * time.Minute)

					// Check if deal should be cancelled after creation
					isCancelled := rnd.Float64() < generatorConfig.CancelAfterDealProbability

					var completedAt *time.Time
					var status model.DealStatus

					if isCancelled {
						// Cancelled deals have no completion time
						status = model.DealStatus("cancelled")
						completedAt = nil
					} else {
						// Random completion time between min and max hours
						completionHours := generatorConfig.DealCompletionMinHours + rnd.Intn(generatorConfig.DealCompletionMaxHours-generatorConfig.DealCompletionMinHours+1)
						completedAtTime := dealTime.Add(time.Duration(completionHours) * time.Hour)
						completedAt = &completedAtTime
						status = model.DealStatus("completed")
					}

					deal := model.Deal{
						ID:           uuid.New(),
						ListingID:    listing.ID,
						BuyerID:      user.ID,
						CreatedAt:    dealTime,
						CompletedAt:  completedAt,
						Status:       status,
						DeliveryUsed: listing.DeliveryAvailable,
						PriceBand:    listing.PriceBand,
					}
					deals = append(deals, deal)

					// Only mark listing as sold if deal is not cancelled
					if !isCancelled {
						soldListings[listing.ID] = true
					}

					// REVIEW stage - only for completed deals
					if !isCancelled && rnd.Float64() < clamp(cfg.Review, 0, 1) && user.ID != listing.SellerID {
						reviewTime := completedAt.Add(time.Duration(generatorConfig.EventReviewOffsetHours) * time.Hour)
						review := model.Review{
							ID:          uuid.New(),
							DealID:      &deal.ID,
							AuthorID:    user.ID,
							RecipientID: listing.SellerID,
							Rating:      int16(3 + rnd.Intn(3)),
							CreatedAt:   reviewTime,
						}
						reviews = append(reviews, review)
					}
				}
			}
		}
	}

	return events, deals, reviews
}

// calculateIntentCount determines how many intents a user should have based on their preferred categories
func calculateIntentCount(preferredCategories []string, categories []CategoryConfig, config GeneratorConfig, rnd *rand.Rand) int {
	if len(preferredCategories) == 0 {
		return config.IntentCountMin + rnd.Intn(2)
	}

	totalVolume := 0.0
	for _, prefCat := range preferredCategories {
		for _, cat := range categories {
			if cat.Name == prefCat {
				totalVolume += cat.PurchaseVolume
				break
			}
		}
	}

	intentCount := int(totalVolume * config.IntentCountMultiplier)
	if intentCount < config.IntentCountMin {
		intentCount = config.IntentCountMin
	}
	if intentCount > config.IntentCountMax {
		intentCount = config.IntentCountMax
	}

	return intentCount
}

func selectIntentCategoryForIndex(cfg UserConfig, categories []CategoryConfig, intentIndex int) string {
	if len(cfg.PreferredCategories) == 0 {
		if len(categories) > 0 {
			return categories[intentIndex%len(categories)].Name
		}
		return ""
	}

	return cfg.PreferredCategories[intentIndex%len(cfg.PreferredCategories)]
}

func selectSearchCandidates(rnd *rand.Rand, listings []ListingWithGenData, user UserWithGenData, cfg UserConfig) []ListingWithGenData {
	if len(listings) == 0 {
		return nil
	}

	count := DefaultGeneratorConfig().SearchCandidatePoolSize
	if count > len(listings) {
		count = len(listings)
	}

	shuffled := make([]ListingWithGenData, len(listings))
	copy(shuffled, listings)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rnd.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}
