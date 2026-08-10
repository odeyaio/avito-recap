package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PriceBand string

type EventType string

type EventSource string

type DealStatus string

const (
	EventTypeView       = "view"
	EventTypeSearch     = "search"
	EventTypeFavorite   = "favorite"
	EventTypeContact    = "contact"
	EventTypeDelete     = "delete"
	EventTypeBoost      = "boost"
	EventTypeAnother    = "another"
	EventTypeEdit       = "edit"
	EventTypeShare      = "share"
	EventTypeCancelDeal = "cancel_deal"
)

const (
	EventSourceNotification = "notification"
	EventSourcePush         = "push"
	EventSourceEmail        = "email"
	EventSourceDirect       = "direct"
)

const (
	EventListingView      EventType = "listing_view"
	EventSearch           EventType = "search"
	EventFavoriteAdd      EventType = "favorite_add"
	EventFavoriteRemove   EventType = "favorite_remove"
	EventContact          EventType = "contact"
	EventNotificationOpen EventType = "notification_open"
	EventPromotionUse     EventType = "promotion_use"
	EventListingEdit      EventType = "listing_edit"
	EventDeliveryEnable   EventType = "delivery_enable"
)

const (
	DealCompleted DealStatus = "completed"
	DealCancelled DealStatus = "cancelled"
)

type User struct {
	ID            uuid.UUID
	DisplayName   string
	RegisteredAt  time.Time
	Region        string
	Timezone      string
	IsTestProfile bool
	CreatedAt     time.Time
}

type Category struct {
	ID       uuid.UUID
	Code     string
	Name     string
	ParentID *uuid.UUID
}

type Listing struct {
	ID                  uuid.UUID
	SellerID            uuid.UUID
	CategoryID          uuid.UUID
	Region              string
	PriceBand           PriceBand
	PublishedAt         time.Time
	ClosedAt            *time.Time
	DeliveryAvailable   bool
	PhotoCount          int
	DescriptionComplete bool
}

type ActivityEvent struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Type            EventType
	OccurredAt      time.Time
	ListingID       *uuid.UUID
	CategoryID      *uuid.UUID
	DurationSeconds *int
	ResultCount     *int
	FilterCount     *int
	TopicKey        *string
	Source          *EventSource
	Properties      json.RawMessage
	IngestedAt      time.Time
}

type Deal struct {
	ID           uuid.UUID
	ListingID    uuid.UUID
	BuyerID      uuid.UUID
	CreatedAt    time.Time
	CompletedAt  *time.Time
	Status       DealStatus
	DeliveryUsed bool
	PriceBand    PriceBand
}

type Review struct {
	ID          uuid.UUID
	DealID      *uuid.UUID
	AuthorID    uuid.UUID
	RecipientID uuid.UUID
	Rating      int16
	CreatedAt   time.Time
}

type Period struct {
	Start time.Time
	End   time.Time
}

func (p Period) Validate() error {
	if p.Start.IsZero() || p.End.IsZero() {
		return errors.New("period boundaries are required")
	}
	if !p.End.After(p.Start) {
		return errors.New("period end must be after start")
	}
	return nil
}

func (p Period) Contains(value time.Time) bool {
	return !value.Before(p.Start) && value.Before(p.End)
}

type Dataset struct {
	User                  User
	Period                Period
	DataCutoffAt          time.Time
	Events                []ActivityEvent
	Categories            []Category
	Listings              []Listing
	Deals                 []Deal
	Reviews               []Review
	EngagementPercentiles map[uuid.UUID]float64
}

type VersionedDataset struct {
	Value   Dataset
	Version string
}
