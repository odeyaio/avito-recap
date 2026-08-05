package catalog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

type BehaviorCatalog struct {
	Version   string               `yaml:"version"`
	Behaviors []BehaviorDefinition `yaml:"behaviors"`
}

type BehaviorDefinition struct {
	Code            string         `yaml:"code"`
	Name            string         `yaml:"name"`
	RuleDescription string         `yaml:"ruleDescription"`
	Enabled         *bool          `yaml:"enabled,omitempty"`
	SortOrder       int            `yaml:"sortOrder"`
	Rule            map[string]any `yaml:"rule"`
	DefaultAction   DefaultAction  `yaml:"defaultAction"`
}

type DefaultAction struct {
	Code     string `json:"code" yaml:"code"`
	Title    string `json:"title" yaml:"title"`
	Resolver string `json:"resolver" yaml:"resolver"`
}

func LoadBehaviors(path string) (BehaviorCatalog, error) {
	file, err := os.Open(path) // #nosec G304 -- the path is an explicit operator-provided CLI argument.
	if err != nil {
		return BehaviorCatalog{}, fmt.Errorf("open behavior catalog: %w", err)
	}
	defer func() { _ = file.Close() }()

	behaviorCatalog, err := DecodeBehaviors(file)
	if err != nil {
		return BehaviorCatalog{}, fmt.Errorf("decode behavior catalog: %w", err)
	}

	return behaviorCatalog, nil
}

func DecodeBehaviors(reader io.Reader) (BehaviorCatalog, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	var behaviorCatalog BehaviorCatalog
	if err := decoder.Decode(&behaviorCatalog); err != nil {
		return BehaviorCatalog{}, err
	}
	if err := ensureSingleDocument(decoder); err != nil {
		return BehaviorCatalog{}, err
	}
	if err := behaviorCatalog.Validate(); err != nil {
		return BehaviorCatalog{}, err
	}

	return behaviorCatalog, nil
}

func (behaviorCatalog BehaviorCatalog) Validate() error {
	if strings.TrimSpace(behaviorCatalog.Version) == "" {
		return errors.New("catalog version is required")
	}
	if len(behaviorCatalog.Behaviors) == 0 {
		return errors.New("catalog must contain at least one behavior")
	}

	codes := make(map[string]struct{}, len(behaviorCatalog.Behaviors))
	orders := make(map[int]struct{}, len(behaviorCatalog.Behaviors))
	for index, behavior := range behaviorCatalog.Behaviors {
		path := fmt.Sprintf("behaviors[%d]", index)
		if !codePattern.MatchString(behavior.Code) {
			return fmt.Errorf("%s.code must match %s", path, codePattern)
		}
		if _, exists := codes[behavior.Code]; exists {
			return fmt.Errorf("duplicate behavior code %q", behavior.Code)
		}
		codes[behavior.Code] = struct{}{}

		if strings.TrimSpace(behavior.Name) == "" {
			return fmt.Errorf("%s.name is required", path)
		}
		if strings.TrimSpace(behavior.RuleDescription) == "" {
			return fmt.Errorf("%s.ruleDescription is required", path)
		}
		if _, exists := orders[behavior.SortOrder]; exists {
			return fmt.Errorf("duplicate behavior sortOrder %d", behavior.SortOrder)
		}
		orders[behavior.SortOrder] = struct{}{}

		if err := validateRule(behavior.Rule, path+".rule"); err != nil {
			return err
		}
		if err := validateDefaultAction(behavior.DefaultAction, path+".defaultAction"); err != nil {
			return err
		}
	}

	return nil
}

func (behavior BehaviorDefinition) IsEnabled() bool {
	return behavior.Enabled == nil || *behavior.Enabled
}

func validateDefaultAction(action DefaultAction, path string) error {
	if !codePattern.MatchString(action.Code) {
		return fmt.Errorf("%s.code must match %s", path, codePattern)
	}
	if strings.TrimSpace(action.Title) == "" {
		return fmt.Errorf("%s.title is required", path)
	}
	if !codePattern.MatchString(action.Resolver) {
		return fmt.Errorf("%s.resolver must match %s", path, codePattern)
	}

	return nil
}
