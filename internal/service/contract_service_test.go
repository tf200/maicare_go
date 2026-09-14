package service

import (
	"strings"
	"testing"
)

func TestNormalizeContractTypeName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{name: "trims whitespace", input: "  Supported living  ", want: "Supported living"},
		{name: "rejects empty", input: " \t ", wantError: true},
		{name: "accepts 100 characters", input: strings.Repeat("a", 100), want: strings.Repeat("a", 100)},
		{name: "rejects more than 100 characters", input: strings.Repeat("a", 101), wantError: true},
		{name: "counts unicode characters", input: strings.Repeat("é", 100), want: strings.Repeat("é", 100)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizeContractTypeName(test.input)
			if test.wantError {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
