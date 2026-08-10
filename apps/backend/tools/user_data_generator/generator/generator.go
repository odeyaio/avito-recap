package generator

import (
	"math/rand"
	"time"

	"avito-recap/internal/model"

	"github.com/google/uuid"
)

// ListingWithGenData wraps model.Listing with generation-only fields
type ListingWithGenData struct {
	model.Listing
	Category string
	Price    float64
}

// UserWithGenData wraps model.User with generation-only fields
type UserWithGenData struct {
	model.User
	search                  float64
	view                    float64
	favorite                float64
	contact                 float64
	deal                    float64
	review                  float64
	notificationOpen        float64
	PreferredCategories    []string
	UnlikelyCategories     []string
	PriceSegment           string
	SellFrequency          float64
	PreferredSellCategories []string
	UnlikelySellCategories []string
}

// GetConfig returns the user configuration as UserConfig
func (u UserWithGenData) GetConfig() UserConfig {
	return UserConfig{
		Search:                  u.search,
		View:                    u.view,
		Favorite:                u.favorite,
		Contact:                 u.contact,
		Deal:                    u.deal,
		Review:                  u.review,
		NotificationOpen:        u.notificationOpen,
		SellFrequency:           u.SellFrequency,
		PriceSegment:            u.PriceSegment,
		PreferredCategories:     u.PreferredCategories,
		UnlikelyCategories:      u.UnlikelyCategories,
		PreferredSellCategories: u.PreferredSellCategories,
		UnlikelySellCategories:  u.UnlikelySellCategories,
	}
}

func randomDateBetweenYearsAgo(rnd *rand.Rand, yearsAgoUpper, yearsAgoLower int) time.Time {
	now := time.Now().UTC()

	older := now.AddDate(-yearsAgoUpper, 0, 0)
	newer := now.AddDate(-yearsAgoLower, 0, 0)

	if older.After(newer) {
		older, newer = newer, older
	}

	diffDays := int(newer.Sub(older).Hours() / 24)
	if diffDays <= 0 {
		return newer
	}

	randomDays := rnd.Intn(diffDays + 1)
	return older.AddDate(0, 0, randomDays)
}

func weightedCategory(rnd *rand.Rand, categories []CategoryConfig) CategoryConfig {
	if len(categories) == 0 {
		return CategoryConfig{}
	}
	var totalWeight int
	for _, category := range categories {
		totalWeight += category.Weight
	}
	pick := rnd.Intn(totalWeight) + 1
	current := 0
	for _, category := range categories {
		current += category.Weight
		if pick <= current {
			return category
		}
	}
	return categories[len(categories)-1]
}

func randomPrice(rnd *rand.Rand, minPrice, maxPrice float64) float64 {
	if maxPrice <= minPrice {
		return minPrice
	}
	return minPrice + rnd.Float64()*(maxPrice-minPrice)
}

func buildCategoryPreferences(rnd *rand.Rand, categories []CategoryConfig) ([]string, []string) {
	if len(categories) == 0 {
		return nil, nil
	}

	preferred := make([]string, 0, 2)
	unlikely := make([]string, 0, 2)

	for i := 0; i < 2 && i < len(categories); i++ {
		category := categories[rnd.Intn(len(categories))]
		if !containsStringSlice(preferred, category.Name) {
			preferred = append(preferred, category.Name)
		}
	}

	for i := 0; i < 2 && i < len(categories); i++ {
		category := categories[rnd.Intn(len(categories))]
		if !containsStringSlice(unlikely, category.Name) && !containsStringSlice(preferred, category.Name) {
			unlikely = append(unlikely, category.Name)
		}
	}

	return preferred, unlikely
}

func selectIntentCategory(user UserConfig, categories []CategoryConfig) string {
	// Deprecated: use selectIntentCategoryForIndex instead
	return selectIntentCategoryForIndex(user, categories, 0)
}

func estimateIntentPrice(user UserConfig, category string, categories []CategoryConfig) float64 {
	switch user.PriceSegment {
	case "budget":
		return 15000
	case "premium":
		return 200000
	default:
		return 50000
	}
}

func filterListingsForIntent(listings []ListingWithGenData, user UserWithGenData, category string, maxPrice float64, cfg UserConfig) []ListingWithGenData {
	filtered := make([]ListingWithGenData, 0)
	for _, listing := range listings {
		if listing.SellerID == user.ID {
			continue
		}
		if listing.Category != category {
			continue
		}
		if listing.Price > maxPrice {
			continue
		}
		// Filter out unlikely categories
		if containsStringSlice(user.UnlikelyCategories, listing.Category) {
			continue
		}
		filtered = append(filtered, listing)
	}
	if len(filtered) == 0 {
		return listings
	}
	return filtered
}

func listingRelevance(listing ListingWithGenData, user UserWithGenData, category uuid.UUID, maxPrice float64, cfg UserConfig) float64 {
	score := 0.25
	if listing.CategoryID == category {
		score += 0.35
	}
	if listing.Price <= maxPrice {
		score += 0.2
	}
	if listing.DeliveryAvailable {
		score += 0.1
	}
	if containsStringSlice(user.PreferredCategories, listing.Category) {
		score += 0.15
	}
	if containsStringSlice(user.UnlikelyCategories, listing.Category) {
		score -= 0.1
	}
	if score > 1 {
		return 1
	}
	if score < 0 {
		return 0
	}
	return score
}

func selectTopCandidates(rnd *rand.Rand, listings []ListingWithGenData, user UserWithGenData, cfg UserConfig) []ListingWithGenData {
	if len(listings) == 0 {
		return nil
	}
	count := 3
	if count > len(listings) {
		count = len(listings)
	}
	return listings[:count]
}

func narrowPool(current []ListingWithGenData, selected []ListingWithGenData, user UserWithGenData, category string) []ListingWithGenData {
	if len(current) <= 1 {
		return nil
	}
	filtered := make([]ListingWithGenData, 0, len(current))
	for _, listing := range current {
		if containsListing(selected, listing) {
			continue
		}
		filtered = append(filtered, listing)
	}
	if len(filtered) == 0 {
		return current
	}
	return filtered
}

func containsListing(listings []ListingWithGenData, target ListingWithGenData) bool {
	for _, listing := range listings {
		if listing.ID == target.ID {
			return true
		}
	}
	return false
}

func findCategoryConfig(categories []CategoryConfig, name string) CategoryConfig {
	for _, category := range categories {
		if category.Name == name {
			return category
		}
	}
	if len(categories) > 0 {
		return categories[0]
	}
	return CategoryConfig{}
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func SellerCountForListings(listingCount int) int {
	count := listingCount / 20
	if count < 5 {
		return 5
	}
	if count > 20 {
		return 20
	}
	return count
}

func containsString(values []uuid.UUID, target uuid.UUID) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsStringSlice(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func randomDateBetweenRange(start, end time.Time, rnd *rand.Rand) time.Time {
	if end.Before(start) {
		return start
	}
	diff := end.Sub(start)
	if diff <= 0 {
		return start
	}
	return start.Add(time.Duration(rnd.Int63n(int64(diff))))
}

func randomTimeOffset(base time.Time, maxSeconds int, rnd *rand.Rand) time.Time {
	if maxSeconds <= 0 {
		return base
	}
	offset := rnd.Intn(maxSeconds + 1)
	return base.Add(time.Duration(offset) * time.Second)
}
