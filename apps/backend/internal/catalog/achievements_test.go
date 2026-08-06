package catalog

import (
	"strings"
	"testing"
)

func TestDecodeAchievements(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "valid catalog",
			input: `
version: v1
achievements:
  - code: explorer
    name: Explorer
    description: Many unique listing views
    iconKey: achievement-icons/explorer.webp
    shareableByDefault: true
    sortOrder: 10
    rule:
      all:
        - metric: activity.unique_listings_viewed
          operator: gte
          value: 100
`,
		},
		{
			name: "unknown top-level field",
			input: `
version: v1
unknown: value
achievements: []
`,
			wantErr: true,
		},
		{
			name: "duplicate code",
			input: `
version: v1
achievements:
  - &achievement
    code: duplicate
    name: First
    description: First rule
    iconKey: achievement-icons/first.webp
    shareableByDefault: false
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: gte
          value: 1
  - <<: *achievement
    name: Second
    iconKey: achievement-icons/second.webp
    sortOrder: 20
`,
			wantErr: true,
		},
		{
			name: "unsupported operator",
			input: `
version: v1
achievements:
  - code: invalid_operator
    name: Invalid operator
    description: Invalid rule
    iconKey: achievement-icons/invalid.webp
    shareableByDefault: false
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: approximately
          value: 10
`,
			wantErr: true,
		},
		{
			name: "non-numeric ordered comparison",
			input: `
version: v1
achievements:
  - code: invalid_value
    name: Invalid value
    description: Invalid rule value
    iconKey: achievement-icons/invalid.webp
    shareableByDefault: false
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: gte
          value: many
`,
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := DecodeAchievements(strings.NewReader(testCase.input))
			assertDecodeResult(t, err, testCase.wantErr)
		})
	}
}

func assertDecodeResult(t *testing.T, err error, wantErr bool) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Fatalf("decode error = %v, wantErr %v", err, wantErr)
	}
}
