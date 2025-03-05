package postgres

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/robertfischer3/scrutiny_cnapp/internal/pkg/database"
	"github.com/robertfischer3/scrutiny_cnapp/internal/pkg/logger"
)

// Provider implements the database.Provider interface for PostgreSQL
type Provider struct {
	logger logger.Logger
}

// NewProvider creates a new PostgreSQL provider
func NewProvider(logger logger.Logger) *Provider {
	return &Provider{
		logger: logger,
	}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "postgres"
}

// Connect establishes a connection to PostgreSQL
func (p *Provider) Connect(config database.Config) (database.Connection, error) {
	connStr := config.ConnectionString
	
	// If SSL is enabled, configure TLS
	if config.SSLMode != "disable" && config.SSLMode != "" {
		// Setup TLS if certificates are provided
		if config.SSLCert != "" && config.SSLKey != "" {
			tlsConfig, err := setupTLS(config)
			if err != nil {
				return nil, fmt.Errorf("failed to setup TLS: %w", err)
			}
			sql.Register("postgres+tls", &wrappedDriver{tlsConfig: tlsConfig})
			connStr += " sslmode=require"
		} else {
			connStr += fmt.Sprintf(" sslmode=%s", config.SSLMode)
		}
	}

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	p.logger.Info("Successfully connected to PostgreSQL database")
	
	return &Connection{
		db:     db,
		logger: p.logger,
	}, nil
}

// Connection implements the database.Connection interface for PostgreSQL
type Connection struct {
	db     *sql.DB
	logger logger.Logger
}

// Execute runs a query without returning any rows
func (c *Connection) Execute(ctx context.Context, query string, args ...interface{}) (database.Result, error) {
	result, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		c.logger.WithField("query", query).WithError(err).Error("Failed to execute query")
		return nil, err
	}
	return &Result{result: result}, nil
}

// Query runs a query that returns rows
func (c *Connection) Query(ctx context.Context, query string, args ...interface{}) (database.Rows, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		c.logger.WithField("query", query).WithError(err).Error("Failed to execute query")
		return nil, err
	}
	return &Rows{rows: rows}, nil
}

// QueryRow runs a query that returns a single row
func (c *Connection) QueryRow(ctx context.Context, query string, args ...interface{}) database.Row {
	row := c.db.QueryRowContext(ctx, query, args...)
	return &Row{row: row}
}

// Begin starts a transaction
func (c *Connection) Begin(ctx context.Context) (database.Transaction, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		c.logger.WithError(err).Error("Failed to begin transaction")
		return nil, err
	}
	return &Transaction{tx: tx, logger: c.logger}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}

// Health checks the database connection
func (c *Connection) Health(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Result implements the database.Result interface
type Result struct {
	result sql.Result
}

// LastInsertId implements the database.Result interface
func (r *Result) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

// RowsAffected implements the database.Result interface
func (r *Result) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}

// Row implements the database.Row interface
type Row struct {
	row *sql.Row
}

// Scan implements the database.Row interface
func (r *Row) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// Rows implements the database.Rows interface
type Rows struct {
	rows *sql.Rows
}

// Next implements the database.Rows interface
func (r *Rows) Next() bool {
	return r.rows.Next()
}

// Scan implements the database.Rows interface
func (r *Rows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

// Close implements the database.Rows interface
func (r *Rows) Close() error {
	return r.rows.Close()
}

// Err implements the database.Rows interface
func (r *Rows) Err() error {
	return r.rows.Err()
}

// Transaction implements the database.Transaction interface
type Transaction struct {
	tx     *sql.Tx
	logger logger.Logger
}

// Execute implements the database.Transaction interface
func (t *Transaction) Execute(ctx context.Context, query string, args ...interface{}) (database.Result, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		t.logger.WithField("query", query).WithError(err).Error("Failed to execute query in transaction")
		return nil, err
	}
	return &Result{result: result}, nil
}

// Query implements the database.Transaction interface
func (t *Transaction) Query(ctx context.Context, query string, args ...interface{}) (database.Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		t.logger.WithField("query", query).WithError(err).Error("Failed to execute query in transaction")
		return nil, err
	}
	return &Rows{rows: rows}, nil
}

// QueryRow implements the database.Transaction interface
func (t *Transaction) QueryRow(ctx context.Context, query string, args ...interface{}) database.Row {
	row := t.tx.QueryRowContext(ctx, query, args...)
	return &Row{row: row}
}

// Commit implements the database.Transaction interface
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback implements the database.Transaction interface
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

// TLS helpers
func setupTLS(config database.Config) (*tls.Config, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(config.SSLCert, config.SSLKey)
	if err != nil {
		return nil, fmt.Errorf("could not load client cert: %w", err)
	}

	// Load CA cert if provided
	var caCertPool *x509.CertPool
	if config.SSLRootCert != "" {
		caCertPool = x509.NewCertPool()
		caCert, err := os.ReadFile(config.SSLRootCert)
		if err != nil {
			return nil, fmt.Errorf("could not read CA cert: %w", err)
		}
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA cert")
		}
	}

	// Create TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// Custom driver for TLS
type wrappedDriver struct {
	tlsConfig *tls.Config
}

// Open implements the driver.Driver interface
func (d *wrappedDriver) Open(name string) (driver.Conn, error) {
	// This is a simplified implementation
	// In a real implementation, you would use the TLS config
	return nil, fmt.Errorf("TLS driver not fully implemented")
}