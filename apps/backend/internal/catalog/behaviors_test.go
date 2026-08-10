package catalog

import (
	"strings"
	"testing"
)

func TestDecodeBehaviors(t *testing.T) {
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
behaviors:
  - code: explorer
    name: Explorer
    description: Explores many listings
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: gte
          value: 10
    defaultAction:
      code: open_feed
      title: Open feed
      targetType: search
      href: https://www.avito.ru/rossiya
`,
		},
		{
			name: "invalid target type",
			input: `
version: v1
behaviors:
  - code: explorer
    name: Explorer
    description: Explores many listings
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: gte
          value: 10
    defaultAction:
      code: open_feed
      title: Open feed
      targetType: Invalid Target
      href: https://www.avito.ru/rossiya
`,
			wantErr: true,
		},
		{
			name: "duplicate code",
			input: `
version: v1
behaviors:
  - &behavior
    code: duplicate
    name: First
    description: First rule
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: gte
          value: 10
    defaultAction:
      code: open_feed
      title: Open feed
      targetType: search
      href: https://www.avito.ru/rossiya
  - <<: *behavior
    name: Second
    sortOrder: 20
`,
			wantErr: true,
		},
		{
			name: "unsupported rule operator",
			input: `
version: v1
behaviors:
  - code: explorer
    name: Explorer
    description: Explores many listings
    sortOrder: 10
    rule:
      all:
        - metric: activity.views
          operator: contains
          value: 10
    defaultAction:
      code: open_feed
      title: Open feed
      targetType: search
      href: https://www.avito.ru/rossiya
`,
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := DecodeBehaviors(strings.NewReader(testCase.input))
			assertDecodeResult(t, err, testCase.wantErr)
		})
	}
}
