package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/repositories"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dataSourceName string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db}, nil
}

func (s *PostgresRepository) App(ctx context.Context, id int64) (domain.App, error) {
	const op = "repositories.app.postgres.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {

		return domain.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app domain.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return domain.App{}, fmt.Errorf("%s: %w", op, repositories.ErrAppNotFound)
		}

		return domain.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
