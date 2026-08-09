package model

import (
	"encoding/json"
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
	EventSourcePush        = "push"
	EventSourceEmail       = "email"
	EventSourceDirect      = "direct"
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
