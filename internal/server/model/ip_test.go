package model_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
)

func TestNewIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want *model.IP
	}{
		{
			name: "valid_ipv4",
			give: "192.0.2.1",
			want: &model.IP{V4: "192.0.2.1"},
		},
		{
			name: "valid_ipv6",
			give: "2001:db8::68",
			want: &model.IP{V6: "2001:db8::68"},
		},
		{
			name: "invalid_ip",
			give: "999.999.999.999",
			want: nil,
		},
		{
			name: "empty_string",
			give: "",
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := model.NewIP(tt.give)

			if got == nil && tt.want == nil {
				return
			}

			if got == nil || tt.want == nil {
				t.Fatalf("NewIP(%q) = %v, want %v", tt.give, got, tt.want)
			}

			if got.V4 != tt.want.V4 || got.V6 != tt.want.V6 {
				t.Fatalf("NewIP(%q) = %v, want %v", tt.give, got, tt.want)
			}
		})
	}
}
