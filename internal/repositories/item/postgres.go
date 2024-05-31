package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (p *PostgresRepository) GetItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.Item, error) {
	const op = "repositories.item.postgres.GetItem"

	stmt, err := p.db.Prepare("SELECT id, type, name, folder_id, user_id, is_favorite FROM items WHERE id = $1 AND user_id=$2")
	if err != nil {
		return &domain.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, itemId.String(), userId.String())

	var item domain.Item
	err = row.Scan(&item.ID, &item.Type, &item.Name, &item.FolderId, &item.UserId, &item.IsFavorite)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &item, fmt.Errorf("%s: %w", op, repositories.ErrItemNotFound)
		}

		return &item, fmt.Errorf("%s: %w", op, err)
	}

	return &item, nil
}

func (p *PostgresRepository) GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error) {
	const op = "repositories.item.postgres.GetItems"

	stmt, err := p.db.Prepare("SELECT id, type, name, folder_id, user_id, is_favorite FROM items WHERE user_id=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userId.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var items []*domain.Item
	for rows.Next() {
		var item domain.Item
		err = rows.Scan(&item.ID, &item.Type, &item.Name, &item.FolderId, &item.UserId, &item.IsFavorite)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}

func (p *PostgresRepository) CreateItem(ctx context.Context, item domain.Item) (uuid.UUID, error) {
	const op = "repositories.item.postgres.CreateItem"

	stmt, err := p.db.Prepare("INSERT INTO items (id, type, name, folder_id, user_id, is_favorite) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5) RETURNING ID")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	var id uuid.UUID
	row := stmt.QueryRowContext(ctx, item.Type, item.Name, item.FolderId, item.UserId, item.IsFavorite)
	err = row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
