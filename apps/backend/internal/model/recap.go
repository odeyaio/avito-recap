package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RecapSnapshot struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	PeriodStart    time.Time
	PeriodEnd      time.Time
	RulesetVersion string
	DatasetVersion string
	DataCutoffAt   time.Time
	Metrics        json.RawMessage
	GeneratedAt    time.Time
}

type BehaviorTypeDefinition struct {
	ID              uuid.UUID
	Code            string
	Name            string
	RuleDescription string
	Rule            json.RawMessage
	RuleVersion     string
	DefaultAction   json.RawMessage
	Enabled         bool
	SortOrder       int
	UpdatedAt       time.Time
}

type AchievementDefinition struct {
	ID                 uuid.UUID
	Code               string
	Name               string
	RuleDescription    string
	Rule               json.RawMessage
	RuleVersion        string
	IconKey            string
	Enabled            bool
	ShareableByDefault bool
	SortOrder          int
	UpdatedAt          time.Time
}

type RecapBehaviorType struct {
	RecapID                  uuid.UUID
	BehaviorTypeDefinitionID uuid.UUID
	IsPrimary                bool
	Position                 int
	Score                    *float64
	Evidence                 json.RawMessage
}

type RecapAchievement struct {
	RecapID                 uuid.UUID
	AchievementDefinitionID uuid.UUID
	Position                int
	AchievedAt              *time.Time
	Evidence                json.RawMessage
	IsShareable             bool
}

type RecapNextAction struct {
	RecapID    uuid.UUID
	Code       string
	Href       string
	Target     json.RawMessage
	Evidence   json.RawMessage
	ResolvedAt time.Time
}

type RecapPresentation struct {
	ID            uuid.UUID
	RecapID       uuid.UUID
	Locale        string
	PromptVersion string
	ModelName     string
	InputHash     string
	Content       json.RawMessage
	GeneratedAt   time.Time
}
