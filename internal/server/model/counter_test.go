package model_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
)

func TestNewCounter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give uint64
		want *model.Counter
	}{
		{
			name: "count_greater_than_1",
			give: 5,
			want: &model.Counter{Count: 5},
		},
		{
			name: "count_equals_to_1",
			give: 1,
			want: &model.Counter{Count: 1},
		},
		{
			name: "count_less_than_1",
			give: 0,
			want: &model.Counter{Count: 1},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := model.NewCounter(tt.give)
			if got.Count != tt.want.Count {
				t.Errorf("NewCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}
