package jsonutil_test

import (
	"encoding/json"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/accio127/internal/jsonutil"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		json         string
		expected     jsonutil.Duration
		expectingErr bool
	}{
		{
			name:         "valid duration",
			json:         `"1h"`,
			expected:     jsonutil.Duration(time.Hour),
			expectingErr: false,
		},
		{
			name:         "valid zero duration",
			json:         `"0s"`,
			expected:     jsonutil.Duration(0),
			expectingErr: false,
		},
		{
			name:         "invalid duration",
			json:         `"invalid"`,
			expectingErr: true,
		},
		{
			name:         "non-string",
			json:         `1`,
			expectingErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var d jsonutil.Duration

			err := json.Unmarshal([]byte(tt.json), &d)
			if (err != nil) != tt.expectingErr {
				t.Errorf("UnmarshalJSON() error = %v, expectingErr %v", err, tt.expectingErr)

				return
			}

			if err == nil && d != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", d, tt.expected)
			}
		})
	}
}
