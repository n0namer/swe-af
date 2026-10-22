package dag

import "testing"

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestNorthstarContainsString(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		target string
		want   bool
	}{
		{
			name:   "positive match",
			values: []string{"source-edit", "verify-edit"},
			target: "source-edit",
			want:   true,
		},
		{
			name:   "negative missing",
			values: []string{"verify-edit"},
			target: "source-edit",
			want:   false,
		},
		{
			name:   "nil values",
			values: nil,
			target: "source-edit",
			want:   false,
		},
		{
			name:   "empty values",
			values: []string{},
			target: "source-edit",
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := containsString(test.values, test.target); got != test.want {
				t.Fatalf("containsString(%v, %q) = %t, want %t", test.values, test.target, got, test.want)
			}
		})
	}
}
