package config_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/config"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid_config",
			path:    "testdata/valid-config.json",
			wantErr: false,
		},
		{
			name:    "invalid_config_missing_cert",
			path:    "testdata/invalid-missing-cert-config.json",
			wantErr: true,
		},
		{
			name:    "invalid_config_missing_privacy_policy",
			path:    "testdata/invalid-missing-privacy-policy-config.json",
			wantErr: true,
		},
		{
			name:    "invalid_config_invalid_privacy_policy_url",
			path:    "testdata/invalid-privacy-policy-url.json",
			wantErr: true,
		},
		{
			name:    "nonexistent_config",
			path:    "testdata/nonexistent.json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := config.LoadConfig(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
