// Package config provides the configuration of the application.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"git.sr.ht/~jamesponddotco/accio127/internal/jsonutil"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	// ErrInvalidConfigFile is returned when the configuration file is invalid
	// or cannot be read.
	ErrInvalidConfigFile xerrors.Error = "invalid configuration file"

	// ErrInvalidPrivacyPolicy is returned when the privacy policy is invalid.
	ErrInvalidPrivacyPolicy xerrors.Error = "invalid privacy policy"

	// ErrProxyRequired is returned when a Config is created without the IP
	// address of the reverse proxy.
	ErrProxyRequired xerrors.Error = "reverse proxy IP address is required"

	// ErrCertRequired is returned when a Config is created without
	// certification files.
	ErrCertRequired xerrors.Error = "certification files are required"

	// ErrPrivacyPolicyRequired is returned when a Config is created without a
	// privacy policy.
	ErrPrivacyPolicyRequired xerrors.Error = "privacy policy is required"
)

const (
	// DefaultAddress is the default address of the application.
	DefaultAddress string = ":1997"

	// DefaultPID is the default path to the PID file.
	DefaultPID string = "/var/run/accio127.pid"

	// DefaultDSN is the default data source name for the SQLite database.
	DefaultDSN string = "file:/var/share/accio127/sqlite.db?cache=shared&mode=rwc&_pragma_cache_size=-20000&_journal_mode=WAL&_synchronous=NORMAL"

	// DefaultMinTLSVersion is the default minimum TLS version supported by the
	// server.
	DefaultMinTLSVersion string = "TLS13"

	// DefaultReadTimeout is the default read timeout for the server.
	DefaultReadTimeout jsonutil.Duration = jsonutil.Duration(5 * time.Second)

	// DefaultWriteTimeout is the default write timeout for the server.
	DefaultWriteTimeout jsonutil.Duration = jsonutil.Duration(10 * time.Second)

	// DefaultIdleTimeout is the default idle timeout for the server.
	DefaultIdleTimeout jsonutil.Duration = jsonutil.Duration(60 * time.Second)
)

// Config holds shared configuration values for the application.
type Config struct {
	// Address is the address of the application.
	Address string `json:"address"`

	// Proxy is the IP address of the trusted reverse proxy.
	Proxy string `json:"proxy"`

	// PID is the path to the process ID file.
	PID string `json:"pid"`

	// DSN is the data source name for the SQLite database.
	DSN string `json:"dsn"`

	// CertFile is the path to the certificate file.
	CertFile string `json:"certFile"`

	// CertKey is the path to the certificate key file.
	CertKey string `json:"certKey"`

	// MinTLSVersion is the minimum TLS version supported by the server.
	MinTLSVersion string `json:"minTLSVersion"`

	// PrivacyPolicy is the link to the service's privacy policy.
	PrivacyPolicy string `json:"privacyPolicy"`

	// ReadTimeout is the read timeout for the server.
	ReadTimeout jsonutil.Duration `json:"readTimeout"`

	// WriteTimeout is the write timeout for the server.
	WriteTimeout jsonutil.Duration `json:"writeTimeout"`

	// IdleTimeout is the idle timeout for the server.
	IdleTimeout jsonutil.Duration `json:"idleTimeout"`
}

// LoadConfig loads the configuration from a file.
func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}
	defer configFile.Close()

	configData, err := io.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	var cfg *Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	if cfg.Address == "" {
		cfg.Address = DefaultAddress
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = DefaultReadTimeout
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = DefaultWriteTimeout
	}

	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = DefaultIdleTimeout
	}

	if cfg.PID == "" {
		cfg.PID = DefaultPID
	}

	if cfg.DSN == "" {
		cfg.DSN = DefaultDSN
	}

	if cfg.MinTLSVersion != "TLS12" && cfg.MinTLSVersion != "TLS13" {
		cfg.MinTLSVersion = DefaultMinTLSVersion
	}

	return cfg, nil
}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	if cfg.Proxy == "" {
		return ErrProxyRequired
	}

	if cfg.CertFile == "" || cfg.CertKey == "" {
		return ErrCertRequired
	}

	if cfg.PrivacyPolicy == "" {
		return ErrPrivacyPolicyRequired
	}

	if _, err := url.Parse(cfg.PrivacyPolicy); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidPrivacyPolicy, err)
	}

	return nil
}
