package generator

type UserConfig struct {
	search              float64
	view                float64
	favorite            float64
	contact             float64
	deal                float64
	review              float64
	PriceSegment        string
	PreferredCategories []string
	UnlikelyCategories  []string
	Pattern             string
}

type CategoryConfig struct {
	Name           string
	MinPrice       float64
	MaxPrice       float64
	Weight         int
	ReturnRate     float64
	PurchaseVolume float64
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
			search:              0.90,
			view:                0.80,
			favorite:            0.20,
			contact:             0.15,
			deal:                0.10,
			review:              0.20,
			PriceSegment:        "budget",
			PreferredCategories: []string{"Электроника", "Телефоны"},
			UnlikelyCategories:  []string{"Мебель", "Авто"},
			Pattern:             "curious",
		},
		{
			search:              0.60,
			view:                0.80,
			favorite:            0.50,
			contact:             0.30,
			deal:                0.15,
			review:              0.10,
			PriceSegment:        "mid",
			PreferredCategories: []string{"Одежда", "Обувь"},
			UnlikelyCategories:  []string{"Авто", "Электроника"},
			Pattern:             "intentional",
		},
		{
			search:              0.70,
			view:                0.90,
			favorite:            0.60,
			contact:             0.50,
			deal:                0.35,
			review:              0.25,
			PriceSegment:        "premium",
			PreferredCategories: []string{"Ноутбуки", "Мебель"},
			UnlikelyCategories:  []string{"Детские товары", "Обувь"},
			Pattern:             "high-intent",
		},
	}
}
