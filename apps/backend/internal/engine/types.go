package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"avito-recap/internal/catalog"
	"avito-recap/internal/model"

	"github.com/google/uuid"
)

var ErrNoActivity = errors.New("no activity in recap period")

type Period = model.Period

func PeriodForYear(year int, timezone string) (Period, error) {
	if year < 1 {
		return Period{}, fmt.Errorf("year must be positive")
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return Period{}, fmt.Errorf("load timezone %q: %w", timezone, err)
	}

	return model.Period{
		Start: time.Date(year, time.January, 1, 0, 0, 0, 0, location),
		End:   time.Date(year+1, time.January, 1, 0, 0, 0, 0, location),
	}, nil
}

type Dataset = model.Dataset

type Config struct {
	FreshListingWindow        time.Duration `env:"ENGINE_FRESH_LISTING_WINDOW" env-default:"48h"`
	LongGap                   time.Duration `env:"ENGINE_LONG_GAP" env-default:"2160h"`
	FastContactWindow         time.Duration `env:"ENGINE_FAST_CONTACT_WINDOW" env-default:"1h"`
	FastListingResponseWindow time.Duration `env:"ENGINE_FAST_LISTING_RESPONSE_WINDOW" env-default:"24h"`
	SignificantCategoryEvents int           `env:"ENGINE_SIGNIFICANT_CATEGORY_EVENTS" env-default:"3"`
	LowSupplyResultCount      int           `env:"ENGINE_LOW_SUPPLY_RESULT_COUNT" env-default:"10"`
	MinimumActions            int           `env:"ENGINE_MINIMUM_ACTIONS" env-default:"1"`
}

type Metrics struct {
	values map[string]any
}

func NewMetrics() Metrics {
	return Metrics{values: make(map[string]any)}
}

func (m *Metrics) Set(code string, value any) {
	if m.values == nil {
		m.values = make(map[string]any)
	}
	m.values[code] = value
}

func (m Metrics) Get(code string) (any, bool) {
	value, ok := m.values[code]
	return value, ok
}

func (m Metrics) Len() int {
	return len(m.values)
}

func (m Metrics) Codes() []string {
	codes := make([]string, 0, len(m.values))
	for code := range m.values {
		codes = append(codes, code)
	}
	slices.Sort(codes)
	return codes
}

func (m Metrics) Values() map[string]any {
	return maps.Clone(m.values)
}

func (m Metrics) MarshalJSON() ([]byte, error) {
	nested := make(map[string]any)
	for code, value := range m.values {
		parts := strings.Split(code, ".")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid metric code %q", code)
		}

		group, ok := nested[parts[0]].(map[string]any)
		if !ok {
			group = make(map[string]any)
			nested[parts[0]] = group
		}
		group[parts[1]] = value
	}

	return json.Marshal(nested)
}

func (m *Metrics) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var nested map[string]map[string]any
	if err := decoder.Decode(&nested); err != nil {
		return fmt.Errorf("decode metrics: %w", err)
	}

	values := make(map[string]any)
	for group, groupValues := range nested {
		if group == "" {
			return errors.New("metric group must not be empty")
		}
		for name, value := range groupValues {
			if name == "" {
				return fmt.Errorf("metric name in group %q must not be empty", group)
			}
			values[group+"."+name] = value
		}
	}
	m.values = values
	return nil
}

type Evidence struct {
	MetricCode string `json:"metricCode"`
	Operator   string `json:"operator"`
	Expected   any    `json:"expected"`
	Actual     any    `json:"actual,omitempty"`
	Matched    bool   `json:"matched"`
}

type BehaviorMatch struct {
	Definition catalog.BehaviorDefinition
	IsPrimary  bool
	Position   int
	Evidence   []Evidence
}

type AchievementMatch struct {
	Definition catalog.AchievementDefinition
	Position   int
	Evidence   []Evidence
}

// CategoryStat is a per-category rollup used for the "interests" section of
// a recap (topCategories / newCategories / mostConsistentCategory). Unlike
// the flat Metrics map, this carries the category's identity (Code/Name) —
// previously calculateInterests only ever aggregated scalar counts and threw
// away which category they belonged to.
type CategoryStat struct {
	CategoryID   uuid.UUID
	Code         string
	Name         string
	Actions      int64
	ActiveMonths int64
	Share        float64
}

type Result struct {
	Metrics      Metrics
	Behaviors    []BehaviorMatch
	Achievements []AchievementMatch

	// TopCategories / NewCategories / MostConsistentCategory carry category
	// identity for interests.* — see CategoryStat.
	TopCategories          []CategoryStat
	NewCategories          []CategoryStat
	MostConsistentCategory *CategoryStat

	// MonthlyActivity is a 12-slot (Jan..Dec) count of activity points in the
	// recap period, indexed by month-1.
	MonthlyActivity [12]int64
}
