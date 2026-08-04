package generator

import (
	"math/rand"
	"time"

	"avito-recap/internal/model"
)

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
		if !containsString(preferred, category.Name) {
			preferred = append(preferred, category.Name)
		}
	}

	for i := 0; i < 2 && i < len(categories); i++ {
		category := categories[rnd.Intn(len(categories))]
		if !containsString(unlikely, category.Name) && !containsString(preferred, category.Name) {
			unlikely = append(unlikely, category.Name)
		}
	}

	return preferred, unlikely
}

func selectIntentCategory(user model.User, categories []CategoryConfig) string {
	if len(categories) == 0 {
		return ""
	}
	for _, preferred := range user.PreferredCategories {
		if len(preferred) > 0 {
			return preferred
		}
	}
	return categories[0].Name
}

func estimateIntentPrice(user model.User, category string) float64 {
	switch user.PriceSegment {
	case "budget":
		return 15000
	case "premium":
		return 200000
	default:
		return 50000
	}
}

func filterListingsForIntent(listings []model.Listing, user model.User, category string, maxPrice float64, cfg UserConfig) []model.Listing {
	filtered := make([]model.Listing, 0)
	for _, listing := range listings {
		if listing.Category != category {
			continue
		}
		if listing.Price > maxPrice {
			continue
		}
		filtered = append(filtered, listing)
	}
	if len(filtered) == 0 {
		return listings
	}
	return filtered
}

func listingRelevance(listing model.Listing, user model.User, category string, maxPrice float64, cfg UserConfig) float64 {
	score := 0.25
	if listing.Category == category {
		score += 0.35
	}
	if listing.Price <= maxPrice {
		score += 0.2
	}
	if listing.DeliveryAvailable {
		score += 0.1
	}
	if containsString(user.PreferredCategories, listing.Category) {
		score += 0.15
	}
	if containsString(user.UnlikelyCategories, listing.Category) {
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

func selectTopCandidates(rnd *rand.Rand, listings []model.Listing, user model.User, cfg UserConfig) []model.Listing {
	if len(listings) == 0 {
		return nil
	}
	count := 3
	if cfg.Pattern == "high-intent" {
		count = 5
	} else if cfg.Pattern == "intentional" {
		count = 4
	}
	if count > len(listings) {
		count = len(listings)
	}
	return listings[:count]
}

func narrowPool(current []model.Listing, selected []model.Listing, user model.User, category string) []model.Listing {
	if len(current) <= 1 {
		return nil
	}
	filtered := make([]model.Listing, 0, len(current))
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

func containsListing(listings []model.Listing, target model.Listing) bool {
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

func randomDateAfter(base time.Time, rnd *rand.Rand, maxDays int) time.Time {
	if maxDays <= 0 {
		return base
	}
	return base.AddDate(0, 0, rnd.Intn(maxDays)+1)
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

func containsString(values []string, target string) bool {
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
