// Package database holds the database connection and functions for the service.
package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"sync"

	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	_ "github.com/mattn/go-sqlite3" //nolint:revive // SQLite3 driver
	"go.uber.org/zap"
)

//go:embed schema.sql
var schema string

// DB wraps the database connection and stores the access counter.
type DB struct {
	db     *sql.DB
	logger *zap.Logger
	count  uint64
	mu     sync.Mutex
}

// Open opens a database connection and returns a DB instance.
func Open(logger *zap.Logger, dsn string) (*DB, error) {
	if logger == nil {
		return nil, errors.ErrNilLogger
	}

	if dsn == "" {
		return nil, errors.ErrEmptyDSN
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	var count uint64

	err = db.QueryRow("SELECT count FROM counter WHERE id = 1").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get access counter: %w", err)
	}

	return &DB{
		db:     db,
		logger: logger,
		count:  count,
	}, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

// Ping checks if the PostgreSQL database is accessible by executing a simple query.
func (d *DB) Ping() error {
	var result int

	err := d.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil || result != 1 {
		return fmt.Errorf("failed to ping the database: %w", err)
	}

	return nil
}

// Count returns the current access counter.
func (d *DB) Count() uint64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.count
}

// Increment increments the access counter and stores the access in the database.
func (d *DB) Increment() (uint64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return d.count, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if err = tx.Rollback(); err != nil {
				d.logger.Error("failed to rollback after panic", zap.Error(err))
			}

			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				d.logger.Error("failed to rollback transaction", zap.Error(rbErr))
			}
		}

		err = tx.Commit()
		if err != nil {
			d.logger.Error("failed to commit transaction", zap.Error(err))
		}
	}()

	stmt, err := tx.Prepare("UPDATE counter SET count = ? WHERE id = 1")
	if err != nil {
		return d.count, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	d.count++

	_, err = stmt.Exec(d.count)
	if err != nil {
		return d.count, fmt.Errorf("failed to increment access counter: %w", err)
	}

	return d.count, nil
}
