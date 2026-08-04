package model

import "time"

type User struct {
	ID                  int64
	Username            string
	RegisterDate        time.Time
	Region              string
	PriceSegment        string
	PreferredCategories []string
	UnlikelyCategories  []string
}

type Listing struct {
	ID                int64
	SellerID          int64
	Category          string
	Price             float64
	Region            string
	PublishedAt       time.Time
	ClosedAt          *time.Time
	DeliveryAvailable bool
}

type View struct {
	ID              int64
	UserID          int64
	ListingID       int64
	Category        string
	ViewedAt        time.Time
	DurationSeconds int
	IsRepeat        bool
}

type Search struct {
	ID          int64
	UserID      int64
	SearchedAt  time.Time
	Topic       string
	Category    *string
	Filters     map[string]any
	ResultCount int
}

type Favorite struct {
	ID         int64
	UserID     int64
	ListingID  int64
	Category   string
	Action     string
	OccurredAt time.Time
}

type Contact struct {
	ID          int64
	UserID      int64
	ListingID   int64
	ContactType string
	OccurredAt  time.Time
}

type Deal struct {
	ID          int64
	BuyerID     int64
	SellerID    int64
	ListingID   int64
	Category    string
	Price       float64
	Delivery    bool
	CreatedAt   time.Time
	CompletedAt time.Time
	Status      string
}

type Review struct {
	ID             int64
	ReviewerID     int64
	ReviewedUserID int64
	DealID         *int64
	Rating         int16
	CreatedAt      time.Time
}

type FeatureUsage struct {
	ID         int64
	UserID     int64
	Feature    string
	Action     string
	OccurredAt time.Time
}

type UserEvent struct {
	ID         int64
	UserID     int64
	EventType  string
	OccurredAt time.Time
}
