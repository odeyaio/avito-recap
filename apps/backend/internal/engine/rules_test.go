package engine

import (
	"encoding/json"
	"testing"
)

func TestEvaluateRuleNestedAndNumericTypes(t *testing.T) {
	t.Parallel()

	metrics := NewMetrics()
	metrics.Set("activity.views", int64(120))
	metrics.Set("activity.returned", true)
	metrics.Set("profile.kind", "buyer")

	rule := map[string]any{
		"all": []any{
			map[string]any{"metric": "activity.views", "operator": "gte", "value": float64(100)},
			map[string]any{"any": []any{
				map[string]any{"metric": "activity.returned", "operator": "eq", "value": true},
				map[string]any{"metric": "profile.kind", "operator": "eq", "value": "seller"},
			}},
		},
	}

	matched, evidence, err := EvaluateRule(rule, metrics)
	if err != nil {
		t.Fatalf("EvaluateRule() error = %v", err)
	}
	if !matched {
		t.Fatal("EvaluateRule() matched = false, want true")
	}
	if len(evidence) != 3 {
		t.Fatalf("len(evidence) = %d, want 3", len(evidence))
	}
}

func TestEvaluateRuleMissingMetricIsFalseForEveryOperator(t *testing.T) {
	t.Parallel()

	for _, operator := range []string{"eq", "neq", "gt", "gte", "lt", "lte"} {
		t.Run(operator, func(t *testing.T) {
			t.Parallel()
			matched, evidence, err := EvaluateRule(map[string]any{
				"metric": "missing.value", "operator": operator, "value": 1,
			}, NewMetrics())
			if err != nil {
				t.Fatalf("EvaluateRule() error = %v", err)
			}
			if matched {
				t.Fatal("EvaluateRule() matched = true, want false")
			}
			if len(evidence) != 1 || evidence[0].Actual != nil {
				t.Fatalf("evidence = %#v, want one item with nil actual", evidence)
			}
		})
	}
}

func TestEvaluateRuleJSONNumber(t *testing.T) {
	t.Parallel()

	metrics := NewMetrics()
	metrics.Set("activity.views", json.Number("10"))
	matched, _, err := EvaluateRule(map[string]any{
		"metric": "activity.views", "operator": "eq", "value": int64(10),
	}, metrics)
	if err != nil {
		t.Fatalf("EvaluateRule() error = %v", err)
	}
	if !matched {
		t.Fatal("EvaluateRule() matched = false, want true")
	}
}

func TestEvaluateRuleRejectsMalformedRule(t *testing.T) {
	t.Parallel()

	_, _, err := EvaluateRule(map[string]any{
		"all":   []any{map[string]any{"metric": "activity.views", "operator": "gte", "value": 1}},
		"extra": true,
	}, NewMetrics())
	if err == nil {
		t.Fatal("EvaluateRule() error = nil, want malformed rule error")
	}
}
