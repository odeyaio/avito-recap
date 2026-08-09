package generator

import "math/rand"

type UserConfig struct {
	Search                  float64
	View                    float64
	Favorite                float64
	Contact                 float64
	Delete                  float64
	Deal                    float64
	Review                  float64
	NotificationOpen        float64
	SellFrequency           float64
	IntentCount             int
	PriceSegment            string
	PreferredCategories     []string
	UnlikelyCategories      []string
	PreferredSellCategories []string
	UnlikelySellCategories  []string
}

type CategoryConfig struct {
	Name           string
	MinPrice       float64
	MaxPrice       float64
	Weight         int
	ReturnRate     float64
	PurchaseVolume float64
}

type GeneratorConfig struct {
	UserRegistrationYearsAgo             int
	SellerRegistrationYearsAgo           int
	SellerSelectionRate                  float64
	ListingPublishWindowDays             int
	ListingPublishWindowFallbackDays     int
	PriceBandLowMax                      float64
	PriceBandHighMin                     float64
	ListingDeliveryProbabilityPercent    int
	ListingPhotosMax                     int
	ListingDescriptionProbabilityPercent int
	EventDefaultCategory                 string
	EventBaseTimeOffsetHours             int
	EventSpreadDays                      int
	IntentCountMultiplier                float64
	IntentCountMin                       int
	IntentCountMax                       int
	IntentSpreadDays                     int
	EventFavoriteOffsetMinutes           int
	EventContactOffsetMinutes            int
	EventDealOffsetMinutes               int
	EventReviewOffsetHours               int
	DealCompletionMinHours               int
	DealCompletionMaxHours               int
	EventRandomSpreadSeconds             int
	ListingEditProbability               float64
	ListingEditMinDays                   int
	ListingEditMaxDays                   int
	ShareProbability                     float64
	CancelDealProbability                float64
	CancelAfterDealProbability           float64
}

func DefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		UserRegistrationYearsAgo:             2,
		SellerRegistrationYearsAgo:           3,
		SellerSelectionRate:                  0.15,
		ListingPublishWindowDays:             30,
		ListingPublishWindowFallbackDays:     7,
		PriceBandLowMax:                      10000,
		PriceBandHighMin:                     100000,
		ListingDeliveryProbabilityPercent:    30,
		ListingPhotosMax:                     8,
		ListingDescriptionProbabilityPercent: 80,
		EventDefaultCategory:                 "Электроника",
		EventBaseTimeOffsetHours:             24,
		EventSpreadDays:                      3,
		IntentCountMultiplier:                5.0,
		IntentCountMin:                       1,
		IntentCountMax:                       15,
		IntentSpreadDays:                     7,
		EventFavoriteOffsetMinutes:           5,
		EventContactOffsetMinutes:            15,
		EventDealOffsetMinutes:               30,
		EventReviewOffsetHours:               2,
		DealCompletionMinHours:               1,
		DealCompletionMaxHours:               72,
		EventRandomSpreadSeconds:             30,
		ListingEditProbability:               0.30,
		ListingEditMinDays:                   1,
		ListingEditMaxDays:                   7,
		ShareProbability:                     0.10,
		CancelDealProbability:                0.20,
		CancelAfterDealProbability:           0.15,
	}
}

func DefaultRegions() []string {
	return []string{
		"Москва",
		"Санкт-Петербург",
		"Новосибирск",
		"Екатеринбург",
	}
}

func DefaultCategories() []CategoryConfig {
	return []CategoryConfig{
		{Name: "Электроника", MinPrice: 500, MaxPrice: 250000, Weight: 25, ReturnRate: 0.65, PurchaseVolume: 0.45},
		{Name: "Телефоны", MinPrice: 3000, MaxPrice: 180000, Weight: 20, ReturnRate: 0.80, PurchaseVolume: 0.55},
		{Name: "Ноутбуки", MinPrice: 15000, MaxPrice: 300000, Weight: 15, ReturnRate: 0.70, PurchaseVolume: 0.50},
		{Name: "Одежда", MinPrice: 300, MaxPrice: 50000, Weight: 15, ReturnRate: 0.95, PurchaseVolume: 0.90},
		{Name: "Обувь", MinPrice: 500, MaxPrice: 40000, Weight: 10, ReturnRate: 0.90, PurchaseVolume: 0.80},
		{Name: "Мебель", MinPrice: 1000, MaxPrice: 150000, Weight: 5, ReturnRate: 0.55, PurchaseVolume: 0.70},
		{Name: "Авто", MinPrice: 50000, MaxPrice: 10000000, Weight: 5, ReturnRate: 0.40, PurchaseVolume: 1.00},
		{Name: "Спорт и отдых", MinPrice: 500, MaxPrice: 100000, Weight: 3, ReturnRate: 0.75, PurchaseVolume: 0.70},
		{Name: "Детские товары", MinPrice: 300, MaxPrice: 50000, Weight: 2, ReturnRate: 0.85, PurchaseVolume: 0.75},
	}
}

func DefaultUserConfig() []UserConfig {
	return []UserConfig{
		{
			Search:                  0.90,
			View:                    0.80,
			Favorite:                0.20,
			Contact:                 0.15,
			Deal:                    0.10,
			Review:                  0.20,
			NotificationOpen:        0.15,
			SellFrequency:           0.08,
			IntentCount:             5,
			PriceSegment:            "budget",
			PreferredCategories:     []string{"Электроника", "Телефоны"},
			UnlikelyCategories:      []string{"Мебель", "Авто"},
			PreferredSellCategories: []string{"Одежда", "Обувь"},
			UnlikelySellCategories:  []string{"Авто", "Мебель"},
		},
		{
			Search:                  0.60,
			View:                    0.80,
			Favorite:                0.50,
			Contact:                 0.30,
			Deal:                    0.15,
			Review:                  0.10,
			NotificationOpen:        0.25,
			SellFrequency:           0.18,
			IntentCount:             8,
			PriceSegment:            "mid",
			PreferredCategories:     []string{"Одежда", "Обувь"},
			UnlikelyCategories:      []string{"Авто", "Электроника"},
			PreferredSellCategories: []string{"Телефоны", "Ноутбуки"},
			UnlikelySellCategories:  []string{"Авто", "Одежда"},
		},
		{
			Search:                  0.70,
			View:                    0.90,
			Favorite:                0.60,
			Contact:                 0.50,
			Deal:                    0.35,
			Review:                  0.25,
			NotificationOpen:        0.35,
			SellFrequency:           0.30,
			IntentCount:             6,
			PriceSegment:            "premium",
			PreferredCategories:     []string{"Ноутбуки", "Мебель"},
			UnlikelyCategories:      []string{"Детские товары", "Обувь"},
			PreferredSellCategories: []string{"Авто", "Мебель"},
			UnlikelySellCategories:  []string{"Детские товары", "Одежда"},
		},
	}
}

// RandomUserConfig generates a random user configuration based on the provided random seed
func RandomUserConfig(rnd *rand.Rand, categories []CategoryConfig) UserConfig {
	search := rnd.Float64()
	view := rnd.Float64()
	favorite := rnd.Float64() * 0.7
	contact := rnd.Float64() * 0.5
	deal := rnd.Float64() * 0.4
	review := rnd.Float64() * 0.5
	notificationOpen := rnd.Float64() * 0.4
	sellFrequency := rnd.Float64() * 0.4

	priceSegments := []string{"budget", "mid", "premium"}
	priceSegment := priceSegments[rnd.Intn(len(priceSegments))]

	preferred, unlikely := buildCategoryPreferences(rnd, categories)
	preferredSell, unlikelySell := buildCategoryPreferences(rnd, categories)

	return UserConfig{
		Search:                  search,
		View:                    view,
		Favorite:                favorite,
		Contact:                 contact,
		Deal:                    deal,
		Review:                  review,
		NotificationOpen:        notificationOpen,
		SellFrequency:           sellFrequency,
		IntentCount:             0, 
		PriceSegment:            priceSegment,
		PreferredCategories:     preferred,
		UnlikelyCategories:      unlikely,
		PreferredSellCategories: preferredSell,
		UnlikelySellCategories:  unlikelySell,
	}
}
