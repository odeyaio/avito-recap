package engine

import "encoding/json"

// categoryStatJSON is the on-disk shape for a CategoryStat inside the
// persisted RecapSnapshot.Metrics JSONB. Kept separate from CategoryStat
// itself so the Go-side struct (which also carries CategoryID, used only
// for in-process dedup/lookups) doesn't leak into the stored payload.
type categoryStatJSON struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Actions      int64   `json:"actions"`
	ActiveMonths int64   `json:"active_months"`
	Share        float64 `json:"share"`
}

func categoryStatsJSON(stats []CategoryStat) []categoryStatJSON {
	result := make([]categoryStatJSON, 0, len(stats))
	for _, stat := range stats {
		result = append(result, categoryStatJSON{
			Code:         stat.Code,
			Name:         stat.Name,
			Actions:      stat.Actions,
			ActiveMonths: stat.ActiveMonths,
			Share:        stat.Share,
		})
	}
	return result
}

// MetricsJSON renders the full payload stored in RecapSnapshot.Metrics: the
// flat rule metrics (activity.*, interests.*, ...) that Metrics.MarshalJSON
// already nests into {"group":{"name":value}}, plus the structured data that
// doesn't fit that flat scalar model — per-category breakdowns and the
// month-by-month activity count. adapter/in/http/mapper.go's
// snapshotMetrics struct reads exactly this shape back out.
func (r Result) MetricsJSON() (json.RawMessage, error) {
	flat, err := json.Marshal(r.Metrics)
	if err != nil {
		return nil, err
	}

	var nested map[string]any
	if err := json.Unmarshal(flat, &nested); err != nil {
		return nil, err
	}

	interests, _ := nested["interests"].(map[string]any)
	if interests == nil {
		interests = make(map[string]any)
	}
	interests["top_categories"] = categoryStatsJSON(r.TopCategories)
	interests["new_categories"] = categoryStatsJSON(r.NewCategories)
	if r.MostConsistentCategory != nil {
		consistent := categoryStatsJSON([]CategoryStat{*r.MostConsistentCategory})[0]
		interests["most_consistent_category"] = consistent
	}
	nested["interests"] = interests

	activity, _ := nested["activity"].(map[string]any)
	if activity == nil {
		activity = make(map[string]any)
	}
	activity["monthly_actions"] = r.MonthlyActivity
	nested["activity"] = activity

	return json.Marshal(nested)
}
