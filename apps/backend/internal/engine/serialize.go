package engine

import "encoding/json"

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
