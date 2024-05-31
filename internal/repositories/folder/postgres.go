package folder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
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

func (p *PostgresRepository) Create(ctx context.Context, folder domain.Folder) (uuid.UUID, error) {
	const op = "repositories.folder.postgres.Create"
	var lastInsertId uuid.UUID
	stmt, err := p.db.Prepare("INSERT INTO folders(user_id, name) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, folder.UserId, folder.Name)
	err = row.Scan(&lastInsertId)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, repositories.ErrFolderExists)
		}
		return uuid.UUID{}, fmt.Errorf("create user failure %e", err)
	}
	return lastInsertId, nil
}

func (p *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Folder, error) {
	const op = "repositories.folder.postgres.Get"

	stmt, err := p.db.Prepare("SELECT id, user_id, name FROM folders WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var folder domain.Folder
	err = row.Scan(&folder.ID, &folder.UserId, &folder.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repositories.ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &folder, nil
}
