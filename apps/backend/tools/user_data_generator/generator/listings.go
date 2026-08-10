package generator

import (
	"encoding/json"
	"math/rand"
	"time"

	"avito-recap/internal/model"
)

func GenerateListingData(seed int64, numListings int, categories []CategoryConfig, sellers []model.User, referenceNow time.Time) []model.Listing {
	if numListings <= 0 || len(sellers) == 0 {
		return nil
	}

	source := rand.NewSource(seed)
	rnd := rand.New(source)
	config := DefaultGeneratorConfig()
	historyStart := referenceNow.AddDate(-config.ListingHistoryYears, 0, 0)

	listings := make([]model.Listing, 0, numListings)
	for i := 0; i < numListings; i++ {
		seller := sellers[rnd.Intn(len(sellers))]
		category := CategoryConfig{}
		if len(categories) > 0 {
			category = categories[rnd.Intn(len(categories))]
		}

		publishedAt := randomDateBetweenRange(historyStart, referenceNow, rnd)
		var closedAt *time.Time
		if rnd.Float64() < 0.5 {
			closedAtValue := randomDateBetweenRange(publishedAt, referenceNow, rnd)
			closedAt = &closedAtValue
		}

		priceBand := "mid"
		if category.MaxPrice > 0 && category.MaxPrice <= config.PriceBandLowMax {
			priceBand = "budget"
		} else if category.MaxPrice > config.PriceBandHighMin {
			priceBand = "premium"
		}

		listings = append(listings, model.Listing{
			ID:                  seededUUID(rnd),
			SellerID:            seller.ID,
			CategoryID:          CategoryID(category.Name),
			Region:              seller.Region,
			PriceBand:           model.PriceBand(priceBand),
			PublishedAt:         publishedAt,
			ClosedAt:            closedAt,
			DeliveryAvailable:   rnd.Intn(100) < config.ListingDeliveryProbabilityPercent,
			PhotoCount:          rnd.Intn(config.ListingPhotosMax) + 1,
			DescriptionComplete: rnd.Intn(100) < config.ListingDescriptionProbabilityPercent,
		})
	}

	return listings
}

// GenerateListingEditEvents generates edit events for listings
func GenerateListingEditEvents(seed int64, listings []model.Listing) []model.ActivityEvent {
	source := rand.NewSource(seed)
	rnd := rand.New(source)
	config := DefaultGeneratorConfig()

	events := make([]model.ActivityEvent, 0)
	for _, listing := range listings {
		// Randomly decide if this listing gets edited
		if rnd.Float64() > config.ListingEditProbability {
			continue
		}

		if listing.ClosedAt != nil {
			continue
		}

		minEditTime := listing.PublishedAt.AddDate(0, 0, config.ListingEditMinDays)
		maxEditTime := listing.PublishedAt.AddDate(0, 0, config.ListingEditMaxDays)
		editTime := randomDateBetweenRange(minEditTime, maxEditTime, rnd)

		listingID := listing.ID
		categoryID := listing.CategoryID
		editProperties, _ := json.Marshal(map[string]interface{}{"price_band": listing.PriceBand})

		events = append(events, model.ActivityEvent{
			ID:         seededUUID(rnd),
			UserID:     listing.SellerID,
			Type:       model.EventType(model.EventTypeEdit),
			OccurredAt: editTime,
			ListingID:  &listingID,
			CategoryID: &categoryID,
			Properties: editProperties,
			IngestedAt: time.Now().UTC(),
		})
	}

	return events
}
