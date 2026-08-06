package catalog

import "testing"

func TestFingerprintCatalog(t *testing.T) {
	t.Parallel()

	first := AchievementCatalog{
		Version: "v1",
		Achievements: []AchievementDefinition{
			{
				Code:        "explorer",
				Name:        "Explorer",
				Description: "Views many listings",
				IconKey:     "achievement-icons/explorer.webp",
				SortOrder:   10,
				Rule: map[string]any{
					"operator": "gte",
					"metric":   "activity.views",
					"value":    100,
				},
			},
		},
	}
	second := first
	second.Achievements = append([]AchievementDefinition(nil), first.Achievements...)
	second.Achievements[0].Rule = map[string]any{
		"value":    100,
		"metric":   "activity.views",
		"operator": "gte",
	}

	firstHash, err := fingerprintCatalog(first)
	if err != nil {
		t.Fatalf("fingerprint first catalog: %v", err)
	}
	secondHash, err := fingerprintCatalog(second)
	if err != nil {
		t.Fatalf("fingerprint second catalog: %v", err)
	}
	if firstHash != secondHash {
		t.Fatalf("equivalent catalogs have different fingerprints: %s != %s", firstHash, secondHash)
	}

	second.Achievements[0].Name = "Changed Explorer"
	changedHash, err := fingerprintCatalog(second)
	if err != nil {
		t.Fatalf("fingerprint changed catalog: %v", err)
	}
	if changedHash == firstHash {
		t.Fatal("changed catalog has the same fingerprint")
	}
}
