package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
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

func (p *PostgresRepository) Create(ctx context.Context, login string, passHash []byte) (uuid.UUID, error) {
	const op = "repositories.user.postgres.Create"
	var lastInsertId uuid.UUID
	stmt, err := p.db.Prepare("INSERT INTO users(id, login, pass_hash) VALUES (gen_random_uuid(), $1, $2) RETURNING id")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, login, string(passHash))
	err = row.Scan(&lastInsertId)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, repositories.ErrUserExists)
		}
		return uuid.UUID{}, fmt.Errorf("create user failure %e", err)
	}
	return lastInsertId, nil
}

func (s *PostgresRepository) Get(ctx context.Context, login string) (*domain.User, error) {
	const op = "repositories.user.postgres.Get"

	stmt, err := s.db.Prepare("SELECT id, login, pass_hash FROM users WHERE login = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, login)

	var user domain.User
	err = row.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repositories.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}
