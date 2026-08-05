package catalog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

var codePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

var supportedOperators = map[string]struct{}{
	"eq":  {},
	"neq": {},
	"gt":  {},
	"gte": {},
	"lt":  {},
	"lte": {},
}

var numericOperators = map[string]struct{}{
	"gt":  {},
	"gte": {},
	"lt":  {},
	"lte": {},
}

type AchievementCatalog struct {
	Version      string                  `yaml:"version"`
	Achievements []AchievementDefinition `yaml:"achievements"`
}

type AchievementDefinition struct {
	Code               string         `yaml:"code"`
	Name               string         `yaml:"name"`
	RuleDescription    string         `yaml:"ruleDescription"`
	IconKey            string         `yaml:"iconKey"`
	Enabled            *bool          `yaml:"enabled,omitempty"`
	ShareableByDefault bool           `yaml:"shareableByDefault"`
	SortOrder          int            `yaml:"sortOrder"`
	Rule               map[string]any `yaml:"rule"`
}

func LoadAchievements(path string) (AchievementCatalog, error) {
	file, err := os.Open(path) // #nosec G304 -- the path is an explicit operator-provided CLI argument.
	if err != nil {
		return AchievementCatalog{}, fmt.Errorf("open achievement catalog: %w", err)
	}
	defer func() { _ = file.Close() }()

	catalog, err := DecodeAchievements(file)
	if err != nil {
		return AchievementCatalog{}, fmt.Errorf("decode achievement catalog: %w", err)
	}

	return catalog, nil
}

func DecodeAchievements(reader io.Reader) (AchievementCatalog, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	var catalog AchievementCatalog
	if err := decoder.Decode(&catalog); err != nil {
		return AchievementCatalog{}, err
	}
	if err := ensureSingleDocument(decoder); err != nil {
		return AchievementCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return AchievementCatalog{}, err
	}

	return catalog, nil
}

func (catalog AchievementCatalog) Validate() error {
	if strings.TrimSpace(catalog.Version) == "" {
		return errors.New("catalog version is required")
	}
	if len(catalog.Achievements) == 0 {
		return errors.New("catalog must contain at least one achievement")
	}

	codes := make(map[string]struct{}, len(catalog.Achievements))
	orders := make(map[int]struct{}, len(catalog.Achievements))
	for index, achievement := range catalog.Achievements {
		path := fmt.Sprintf("achievements[%d]", index)
		if !codePattern.MatchString(achievement.Code) {
			return fmt.Errorf("%s.code must match %s", path, codePattern)
		}
		if _, exists := codes[achievement.Code]; exists {
			return fmt.Errorf("duplicate achievement code %q", achievement.Code)
		}
		codes[achievement.Code] = struct{}{}

		if strings.TrimSpace(achievement.Name) == "" {
			return fmt.Errorf("%s.name is required", path)
		}
		if strings.TrimSpace(achievement.RuleDescription) == "" {
			return fmt.Errorf("%s.ruleDescription is required", path)
		}
		if strings.TrimSpace(achievement.IconKey) == "" || strings.HasPrefix(achievement.IconKey, "/") {
			return fmt.Errorf("%s.iconKey must be a relative object key", path)
		}
		if _, exists := orders[achievement.SortOrder]; exists {
			return fmt.Errorf("duplicate achievement sortOrder %d", achievement.SortOrder)
		}
		orders[achievement.SortOrder] = struct{}{}

		if err := validateRule(achievement.Rule, path+".rule"); err != nil {
			return err
		}
	}

	return nil
}

func (achievement AchievementDefinition) IsEnabled() bool {
	return achievement.Enabled == nil || *achievement.Enabled
}

func ensureSingleDocument(decoder *yaml.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err == nil {
		return errors.New("catalog must contain exactly one YAML document")
	} else if !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func validateRule(rule map[string]any, path string) error {
	if len(rule) == 0 {
		return fmt.Errorf("%s must not be empty", path)
	}

	if children, ok := rule["all"]; ok {
		if len(rule) != 1 {
			return fmt.Errorf("%s all rule must not contain other fields", path)
		}
		return validateRuleChildren(children, path+".all")
	}
	if children, ok := rule["any"]; ok {
		if len(rule) != 1 {
			return fmt.Errorf("%s any rule must not contain other fields", path)
		}
		return validateRuleChildren(children, path+".any")
	}

	metric, metricOK := rule["metric"].(string)
	operator, operatorOK := rule["operator"].(string)
	value, valueOK := rule["value"]
	if len(rule) != 3 || !metricOK || strings.TrimSpace(metric) == "" || !operatorOK || strings.TrimSpace(operator) == "" || !valueOK {
		return fmt.Errorf("%s must contain only metric, operator and value", path)
	}
	if _, supported := supportedOperators[operator]; !supported {
		return fmt.Errorf("%s.operator %q is not supported", path, operator)
	}
	if !isScalar(value) {
		return fmt.Errorf("%s.value must be a string, number or boolean", path)
	}
	if _, numeric := numericOperators[operator]; numeric && !isNumber(value) {
		return fmt.Errorf("%s.value must be a number for operator %q", path, operator)
	}

	return nil
}

func isScalar(value any) bool {
	if _, ok := value.(string); ok {
		return true
	}
	if _, ok := value.(bool); ok {
		return true
	}

	return isNumber(value)
}

func isNumber(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

func validateRuleChildren(value any, path string) error {
	children, ok := value.([]any)
	if !ok || len(children) == 0 {
		return fmt.Errorf("%s must be a non-empty list", path)
	}

	for index, child := range children {
		rule, ok := child.(map[string]any)
		if !ok {
			return fmt.Errorf("%s[%d] must be a rule object", path, index)
		}
		if err := validateRule(rule, fmt.Sprintf("%s[%d]", path, index)); err != nil {
			return err
		}
	}

	return nil
}
