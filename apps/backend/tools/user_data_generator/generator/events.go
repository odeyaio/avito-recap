package generator

import (
	"math/rand"

	"avito-recap/internal/model"
)

func GenerateUserEvents(seed int64, users []model.User, config []UserConfig, listings []model.Listing, categories []CategoryConfig) ([]model.Search, []model.View, []model.Favorite, []model.Contact, []model.Deal, []model.Review, []model.UserEvent) {
	source := rand.NewSource(seed)
	rnd := rand.New(source)

	searches := make([]model.Search, 0)
	views := make([]model.View, 0)
	favorites := make([]model.Favorite, 0)
	contacts := make([]model.Contact, 0)
	deals := make([]model.Deal, 0)
	reviews := make([]model.Review, 0)

	for idx, user := range users {
		if len(listings) == 0 {
			continue
		}

		cfg := config[idx%len(config)]
		baseTime := user.RegisterDate
		lastTime := baseTime

		intentCategory := selectIntentCategory(user, categories)
		categoryBehavior := findCategoryConfig(categories, intentCategory)
		intentPrice := estimateIntentPrice(user, intentCategory)
		searchTopic := intentCategory
		if intentCategory == "" {
			intentCategory = "Электроника"
		}

		candidateListings := filterListingsForIntent(listings, user, intentCategory, intentPrice, cfg)
		if len(candidateListings) == 0 {
			candidateListings = listings
		}

		searches = append(searches, model.Search{ID: int64(len(searches) + 1), UserID: user.ID, SearchedAt: baseTime, Topic: searchTopic, Category: &intentCategory, Filters: map[string]any{"category": intentCategory, "max_price": intentPrice}, ResultCount: len(candidateListings)})

		currentPool := candidateListings
		for step := 0; step < 3; step++ {
			if len(currentPool) == 0 {
				break
			}
			selected := selectTopCandidates(rnd, currentPool, user, cfg)
			if len(selected) == 0 {
				selected = currentPool[:min(3, len(currentPool))]
			}

			for _, listing := range selected {
				score := listingRelevance(listing, user, intentCategory, intentPrice, cfg)
				if score < 0.2 {
					continue
				}

				currentTime := randomDateAfter(lastTime, rnd, 2)
				if listing.PublishedAt.After(currentTime) {
					currentTime = listing.PublishedAt
				}
				if currentTime.After(lastTime) {
					lastTime = currentTime
				}

				viewProbability := clamp(cfg.view*score*(0.75+categoryBehavior.ReturnRate*0.35), 0, 1)
				favoriteProbability := clamp(cfg.favorite*score*(0.70+categoryBehavior.ReturnRate*0.40), 0, 1)
				contactProbability := clamp(cfg.contact*score*(0.75+categoryBehavior.ReturnRate*0.25), 0, 1)

				switch step {
				case 0:
					if rnd.Float64() < viewProbability {
						views = append(views, model.View{ID: int64(len(views) + 1), UserID: user.ID, ListingID: listing.ID, Category: listing.Category, ViewedAt: currentTime, DurationSeconds: int(score*600) + 30, IsRepeat: rnd.Intn(100) < 20})
					}
				case 1:
					if rnd.Float64() < favoriteProbability || (cfg.Pattern == "curious" && score > 0.7) {
						favorites = append(favorites, model.Favorite{ID: int64(len(favorites) + 1), UserID: user.ID, ListingID: listing.ID, Category: listing.Category, Action: "add", OccurredAt: currentTime})
					}
				case 2:
					if rnd.Float64() < contactProbability {
						contacts = append(contacts, model.Contact{ID: int64(len(contacts) + 1), UserID: user.ID, ListingID: listing.ID, ContactType: []string{"chat", "call", "offer"}[rnd.Intn(3)], OccurredAt: currentTime})

						dealProbability := clamp(cfg.deal*score*(0.75+categoryBehavior.PurchaseVolume*0.25), 0, 1)
						if rnd.Float64() < dealProbability {
							dealID := int64(len(deals) + 1)
							status := "cancelled"
							completedAt := currentTime
							if rnd.Float64() < clamp(0.65+cfg.deal*0.25, 0, 1) {
								status = "completed"
								completedAt = randomDateAfter(currentTime, rnd, 14)
							}

							deal := model.Deal{
								ID:          dealID,
								BuyerID:     user.ID,
								SellerID:    listing.SellerID,
								ListingID:   listing.ID,
								Category:    listing.Category,
								Price:       listing.Price,
								Delivery:    listing.DeliveryAvailable,
								CreatedAt:   currentTime,
								CompletedAt: completedAt,
								Status:      status,
							}
							deals = append(deals, deal)

							if status == "completed" && rnd.Float64() < clamp(cfg.review*0.8+0.2, 0, 1) {
								reviewID := int64(len(reviews) + 1)
								rating := int16(3 + rnd.Intn(3))
								reviewCreatedAt := randomDateAfter(completedAt, rnd, 7)
								reviews = append(reviews, model.Review{ID: reviewID, ReviewerID: user.ID, ReviewedUserID: listing.SellerID, DealID: &deal.ID, Rating: rating, CreatedAt: reviewCreatedAt})
							}
						}
					}
				}
			}

			currentPool = narrowPool(currentPool, selected, user, intentCategory)
		}
	}

	return searches, views, favorites, contacts, deals, reviews, buildUserEventsFromActions(searches, views, favorites, contacts, deals, reviews)
}

func buildUserEventsFromActions(searches []model.Search, views []model.View, favorites []model.Favorite, contacts []model.Contact, deals []model.Deal, reviews []model.Review) []model.UserEvent {
	capacity := len(searches) + len(views) + len(favorites) + len(contacts) + len(deals)*2 + len(reviews)
	events := make([]model.UserEvent, 0, capacity)
	nextID := int64(1)

	for _, search := range searches {
		events = append(events, model.UserEvent{ID: nextID, UserID: search.UserID, EventType: "search", OccurredAt: search.SearchedAt})
		nextID++
	}
	for _, view := range views {
		events = append(events, model.UserEvent{ID: nextID, UserID: view.UserID, EventType: "view", OccurredAt: view.ViewedAt})
		nextID++
	}
	for _, favorite := range favorites {
		events = append(events, model.UserEvent{ID: nextID, UserID: favorite.UserID, EventType: "favorite", OccurredAt: favorite.OccurredAt})
		nextID++
	}
	for _, contact := range contacts {
		events = append(events, model.UserEvent{ID: nextID, UserID: contact.UserID, EventType: "contact", OccurredAt: contact.OccurredAt})
		nextID++
	}
	for _, deal := range deals {
		events = append(events, model.UserEvent{ID: nextID, UserID: deal.BuyerID, EventType: "deal_created", OccurredAt: deal.CreatedAt})
		nextID++
		if deal.Status == "completed" {
			events = append(events, model.UserEvent{ID: nextID, UserID: deal.BuyerID, EventType: "deal_closed", OccurredAt: deal.CompletedAt})
			nextID++
		}
	}
	for _, review := range reviews {
		events = append(events, model.UserEvent{ID: nextID, UserID: review.ReviewerID, EventType: "review", OccurredAt: review.CreatedAt})
		nextID++
	}

	return events
}
