package generator

import (
	"math/rand"
	"time"

	"avito-recap/internal/model"
)

func GenerateListingData(seed int64, numListings int, categories []CategoryConfig, sellers []model.User, openUntil time.Time) []model.Listing {
	if numListings <= 0 {
		return nil
	}

	source := rand.NewSource(seed)
	rnd := rand.New(source)

	listings := make([]model.Listing, 0, numListings)
	for i := 0; i < numListings; i++ {
		category := weightedCategory(rnd, categories)
		price := randomPrice(rnd, category.MinPrice, category.MaxPrice)
		publishedWindowStart := openUntil.AddDate(0, 0, -30)
		if publishedWindowStart.After(openUntil) {
			publishedWindowStart = openUntil.AddDate(0, 0, -7)
		}
		publishedAt := randomDateBetweenRange(publishedWindowStart, openUntil, rnd)
		var closedAt *time.Time
		if openUntil.After(publishedAt) {
			closedAtValue := randomDateBetweenRange(publishedAt, openUntil, rnd)
			closedAt = &closedAtValue
		}

		sellerID := sellers[0].ID
		if len(sellers) > 1 {
			sellerID = sellers[rnd.Intn(len(sellers))].ID
		}

		listingRegion := ""
		for _, seller := range sellers {
			if seller.ID == sellerID {
				listingRegion = seller.Region
				break
			}
		}

		listings = append(listings, model.Listing{
			ID:                int64(i + 1),
			SellerID:          sellerID,
			Category:          category.Name,
			Price:             price,
			Region:            listingRegion,
			PublishedAt:       publishedAt,
			ClosedAt:          closedAt,
			DeliveryAvailable: rnd.Intn(100) < 30,
		})
	}

	return listings
}
