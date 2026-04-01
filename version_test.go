package pgextras

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"equal single", "1", "1", 0},
		{"equal multi", "1.8", "1.8", 0},
		{"equal three", "1.8.0", "1.8.0", 0},
		{"less major", "1.7", "1.8", -1},
		{"greater major", "1.9", "1.8", 1},
		{"less minor", "1.8.0", "1.8.1", -1},
		{"greater minor", "1.8.2", "1.8.1", 1},
		{"different lengths a shorter", "1", "1.1", -1},
		{"different lengths b shorter", "1.1", "1", 1},
		{"different lengths equal", "1.0", "1", 0},
		{"major version diff", "2.0", "1.9", 1},
		{"zero vs zero", "0", "0", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareVersions(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
