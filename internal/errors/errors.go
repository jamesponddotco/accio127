// Package errors provides custom error types for the application.
package errors

import "git.sr.ht/~jamesponddotco/xstd-go/xerrors"

const (
	// ErrNilLogger is returned when a nil logger is passed to a function.
	ErrNilLogger xerrors.Error = "logger cannot be nil"

	// ErrEmptyDSN is returned when an empty DSN is passed to a function.
	ErrEmptyDSN xerrors.Error = "dsn cannot be empty"
)
