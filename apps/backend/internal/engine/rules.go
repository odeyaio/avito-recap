package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
)

type ruleEvaluation struct {
	matched  bool
	evidence []Evidence
}

func EvaluateRule(rule map[string]any, metrics Metrics) (bool, []Evidence, error) {
	evaluation, err := evaluateRule(rule, metrics)
	if err != nil {
		return false, nil, err
	}

	return evaluation.matched, evaluation.evidence, nil
}

func evaluateRule(rule map[string]any, metrics Metrics) (ruleEvaluation, error) {
	if children, ok := rule["all"]; ok {
		if len(rule) != 1 {
			return ruleEvaluation{}, errorsAtRule("all group must not contain other fields")
		}
		return evaluateChildren(children, metrics, true)
	}
	if children, ok := rule["any"]; ok {
		if len(rule) != 1 {
			return ruleEvaluation{}, errorsAtRule("any group must not contain other fields")
		}
		return evaluateChildren(children, metrics, false)
	}

	metricCode, metricOK := rule["metric"].(string)
	operator, operatorOK := rule["operator"].(string)
	expected, valueOK := rule["value"]
	if len(rule) != 3 || !metricOK || metricCode == "" || !operatorOK || !valueOK {
		return ruleEvaluation{}, errorsAtRule("condition must contain metric, operator and value")
	}

	actual, exists := metrics.Get(metricCode)
	matched, err := compare(actual, exists, operator, expected)
	if err != nil {
		return ruleEvaluation{}, fmt.Errorf("evaluate metric %s: %w", metricCode, err)
	}

	return ruleEvaluation{
		matched: matched,
		evidence: []Evidence{{
			MetricCode: metricCode,
			Operator:   operator,
			Expected:   expected,
			Actual:     actual,
			Matched:    matched,
		}},
	}, nil
}

func evaluateChildren(value any, metrics Metrics, requireAll bool) (ruleEvaluation, error) {
	children, ok := value.([]any)
	if !ok || len(children) == 0 {
		return ruleEvaluation{}, errorsAtRule("logical group must be a non-empty list")
	}

	matched := requireAll
	evidence := make([]Evidence, 0, len(children))
	for index, child := range children {
		childRule, ok := child.(map[string]any)
		if !ok {
			return ruleEvaluation{}, fmt.Errorf("rule child %d must be an object", index)
		}

		result, err := evaluateRule(childRule, metrics)
		if err != nil {
			return ruleEvaluation{}, fmt.Errorf("rule child %d: %w", index, err)
		}
		evidence = append(evidence, result.evidence...)
		if requireAll {
			matched = matched && result.matched
		} else {
			matched = matched || result.matched
		}
	}

	return ruleEvaluation{matched: matched, evidence: evidence}, nil
}

func compare(actual any, exists bool, operator string, expected any) (bool, error) {
	if !exists {
		return false, nil
	}

	switch operator {
	case "eq", "neq":
		equal := scalarEqual(actual, expected)
		if operator == "neq" {
			return !equal, nil
		}
		return equal, nil
	case "gt", "gte", "lt", "lte":
		actualNumber, actualOK := number(actual)
		expectedNumber, expectedOK := number(expected)
		if !actualOK || !expectedOK {
			return false, nil
		}
		switch operator {
		case "gt":
			return actualNumber > expectedNumber, nil
		case "gte":
			return actualNumber >= expectedNumber, nil
		case "lt":
			return actualNumber < expectedNumber, nil
		default:
			return actualNumber <= expectedNumber, nil
		}
	default:
		return false, fmt.Errorf("unsupported operator %q", operator)
	}
}

func scalarEqual(left, right any) bool {
	leftNumber, leftNumeric := number(left)
	rightNumber, rightNumeric := number(right)
	if leftNumeric || rightNumeric {
		return leftNumeric && rightNumeric && leftNumber == rightNumber
	}

	return reflect.DeepEqual(left, right)
}

func number(value any) (float64, bool) {
	var result float64
	switch typed := value.(type) {
	case int:
		result = float64(typed)
	case int8:
		result = float64(typed)
	case int16:
		result = float64(typed)
	case int32:
		result = float64(typed)
	case int64:
		result = float64(typed)
	case uint:
		result = float64(typed)
	case uint8:
		result = float64(typed)
	case uint16:
		result = float64(typed)
	case uint32:
		result = float64(typed)
	case uint64:
		result = float64(typed)
	case float32:
		result = float64(typed)
	case float64:
		result = typed
	case json.Number:
		parsed, err := typed.Float64()
		if err != nil {
			return 0, false
		}
		result = parsed
	default:
		return 0, false
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, false
	}
	return result, true
}

func errorsAtRule(message string) error {
	return fmt.Errorf("invalid rule: %s", message)
}
