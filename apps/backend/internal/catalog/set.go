package catalog

import "github.com/google/uuid"

type Set struct {
	Achievements         AchievementCatalog
	Behaviors            BehaviorCatalog
	AchievementIDs       map[string]uuid.UUID
	BehaviorIDs          map[string]uuid.UUID
	AchievementVersionID uuid.UUID
	BehaviorVersionID    uuid.UUID
}

func (s Set) RulesetVersion() string {
	return "behavior:" + s.Behaviors.Version + "|achievement:" + s.Achievements.Version
}
