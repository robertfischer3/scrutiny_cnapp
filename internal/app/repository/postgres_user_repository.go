package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/robertfischer3/scrutiny_cnapp/internal/app/service"
	"github.com/robertfischer3/scrutiny_cnapp/internal/pkg/database"
	appErrors "github.com/robertfischer3/scrutiny_cnapp/internal/pkg/errors"
	"github.com/robertfischer3/scrutiny_cnapp/internal/pkg/logger"
)

// PostgresUserRepository implements service.UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db     database.Connection
	logger logger.Logger
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db database.Connection, logger logger.Logger) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:     db,
		logger: logger,
	}
}

// FindByID retrieves a user by their ID
func (r *PostgresUserRepository) FindByID(id int) (service.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, role, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user service.User
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Active,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return service.User{}, appErrors.NewNotFoundError("user not found", nil)
		}
		return service.User{}, appErrors.NewDatabaseError("error retrieving user", err)
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	return user, nil
}

// FindAll retrieves all users
func (r *PostgresUserRepository) FindAll() ([]service.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, role, active, created_at, updated_at
		FROM users
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, appErrors.NewDatabaseError("error retrieving users", err)
	}
	defer rows.Close()

	var users []service.User

	for rows.Next() {
		var user service.User
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Active,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, appErrors.NewDatabaseError("error scanning user", err)
		}

		user.CreatedAt = createdAt.Format(time.RFC3339)
		user.UpdatedAt = updatedAt.Format(time.RFC3339)

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, appErrors.NewDatabaseError("error iterating users", err)
	}

	return users, nil
}

// Create creates a new user
func (r *PostgresUserRepository) Create(user service.User) (service.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return service.User{}, appErrors.NewDatabaseError("failed to begin transaction", err)
	}
	defer tx.Rollback() // Rollback if not committed

	now := time.Now().UTC()

	query := `
		INSERT INTO users (name, email, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.Role,
		user.Active,
		now,
		now,
	)

	err = row.Scan(&user.ID)
	if err != nil {
		return service.User{}, appErrors.NewDatabaseError("failed to create user", err)
	}

	user.CreatedAt = now.Format(time.RFC3339)
	user.UpdatedAt = now.Format(time.RFC3339)

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return service.User{}, appErrors.NewDatabaseError("failed to commit user creation", err)
	}

	return user, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(user service.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return appErrors.NewDatabaseError("failed to begin transaction", err)
	}
	defer tx.Rollback() // Rollback if not committed

	now := time.Now().UTC()

	query := `
		UPDATE users
		SET name = $1, email = $2, role = $3, active = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := tx.Execute(ctx, query,
		user.Name,
		user.Email,
		user.Role,
		user.Active,
		now,
		user.ID,
	)
	if err != nil {
		return appErrors.NewDatabaseError("failed to update user", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return appErrors.NewDatabaseError("failed to get affected rows", err)
	}
	if rowsAffected == 0 {
		return appErrors.NewNotFoundError(fmt.Sprintf("user with ID %d not found", user.ID), nil)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return appErrors.NewDatabaseError("failed to commit user update", err)
	}

	return nil
}

// Delete deletes a user by their ID
func (r *PostgresUserRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return appErrors.NewDatabaseError("failed to begin transaction", err)
	}
	defer tx.Rollback() // Rollback if not committed

	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := tx.Execute(ctx, query, id)
	if err != nil {
		return appErrors.NewDatabaseError("failed to delete user", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return appErrors.NewDatabaseError("failed to get affected rows", err)
	}
	if rowsAffected == 0 {
		return appErrors.NewNotFoundError(fmt.Sprintf("user with ID %d not found", id), nil)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return appErrors.NewDatabaseError("failed to commit user deletion", err)
	}

	return nil
}