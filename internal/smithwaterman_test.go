package internal

import "testing"

func TestSmithWaterman(t *testing.T) {
	tests := []struct {
		query      string
		candidates []string
		expected   string
	}{
		{
			query:      "Excalibur Prime Blueprint",
			candidates: []string{"Excalibur Prime Blueprint", "Excalibur Prime Chassis", "Mag Prime Blueprint"},
			expected:   "Excalibur Prime Blueprint",
		},
		{
			query:      "Excalbur Prime Bluepnt",
			candidates: []string{"Excalibur Prime Blueprint", "Excalibur Prime Chassis", "Mag Prime Blueprint"},
			expected:   "Excalibur Prime Blueprint",
		},
		{
			query:      "Mag Prime",
			candidates: []string{"Excalibur Prime Blueprint", "Excalibur Prime Chassis", "Mag Prime Blueprint"},
			expected:   "Mag Prime Blueprint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			actual := smithWaterman(tt.query, tt.candidates)
			if actual != tt.expected {
				t.Errorf("smithWaterman(%s, %v) = %s; want %s", tt.query, tt.candidates, actual, tt.expected)
			}
		})
	}
}
