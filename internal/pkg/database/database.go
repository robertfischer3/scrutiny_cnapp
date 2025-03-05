package database

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	ErrNoRows = errors.New("no rows in result set")
	ErrTxDone = errors.New("transaction has already been committed or rolled back")
)

// Config holds database configuration parameters
type Config struct {
	Driver           string
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	SSLMode          string
	SSLCert          string
	SSLKey           string
	SSLRootCert      string
}

// Connection represents a database connection
type Connection interface {
	// Execute runs a query without returning any rows
	Execute(ctx context.Context, query string, args ...interface{}) (Result, error)
	
	// Query runs a query that returns rows
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	
	// QueryRow runs a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	
	// Begin starts a transaction
	Begin(ctx context.Context) (Transaction, error)
	
	// Close closes the database connection
	Close() error
	
	// Health checks the database connection
	Health(ctx context.Context) error
}

// Result represents a query result
type Result interface {
	// LastInsertId returns the id of the last inserted row
	LastInsertId() (int64, error)
	
	// RowsAffected returns the number of rows affected by the query
	RowsAffected() (int64, error)
}

// Row represents a single database row
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows represents multiple database rows
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// Transaction represents a database transaction
type Transaction interface {
	// Execute runs a query within the transaction
	Execute(ctx context.Context, query string, args ...interface{}) (Result, error)
	
	// Query runs a query within the transaction that returns rows
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	
	// QueryRow runs a query within the transaction that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
}

// Provider is an interface for database provider implementations
type Provider interface {
	// Name returns the name of the database provider
	Name() string
	
	// Connect establishes a connection to the database
	Connect(config Config) (Connection, error)
}