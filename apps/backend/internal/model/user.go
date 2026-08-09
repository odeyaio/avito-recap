package model

import (
	"time"
)

type Search struct {
	ID          int64
	UserID      any
	SearchedAt  time.Time
	Topic       string
	Category    *string
	Filters     map[string]any
	ResultCount int
}

type View struct {
	ID              int64
	UserID          any
	ListingID       any
	Category        string
	ViewedAt        time.Time
	DurationSeconds int
	IsRepeat        bool
}

type Favorite struct {
	ID         int64
	UserID     any
	ListingID  any
	Category   string
	Action     string
	OccurredAt time.Time
}

type Contact struct {
	ID          int64
	UserID      any
	ListingID   any
	ContactType string
	OccurredAt  time.Time
}

type UserActivity struct {
	UserEvents []ActivityEvent
	Searches   []Search
	Views      []View
	Favorites  []Favorite
	Contacts   []Contact
	Deals      []Deal
	Reviews    []Review
	Listings   []Listing
	ViewedListings []Listing
}

type UserInterestPerMonth struct {
	Views         [12]int
	Searches      [12]int
	Favorites     [12]int
	Contacts      [12]int
	Deals         [12]int
	Reviews       [12]int
	FeatureUsages [12]int
}
