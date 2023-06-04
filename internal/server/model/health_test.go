package model_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
)

func TestNewHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveName    string
		giveVersion string
		giveDeps    []model.Dependency
		want        *model.Health
	}{
		{
			name:        "no_dependencies",
			giveName:    "TestService",
			giveVersion: "v1.0.0",
			giveDeps:    nil,
			want: &model.Health{
				Name:         "TestService",
				Version:      "v1.0.0",
				Dependencies: nil,
			},
		},
		{
			name:        "with_dependencies",
			giveName:    "TestService",
			giveVersion: "v1.0.0",
			giveDeps: []model.Dependency{
				{
					Service: "Database",
					Status:  "healthy",
				},
			},
			want: &model.Health{
				Name:    "TestService",
				Version: "v1.0.0",
				Dependencies: []model.Dependency{
					{
						Service: "Database",
						Status:  "healthy",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := model.NewHealth(tt.giveName, tt.giveVersion, tt.giveDeps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}
