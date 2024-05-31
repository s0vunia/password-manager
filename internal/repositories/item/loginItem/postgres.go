package loginItem

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

type ItemSaver interface {
	CreateItem(ctx context.Context, item domain.Item) (uuid.UUID, error)
}

type ItemProvider interface {
	GetItem(ctx context.Context, itemId, userId uuid.UUID) (*domain.Item, error)
	GetItems(ctx context.Context, userId uuid.UUID) ([]*domain.Item, error)
}

type PostgresRepository struct {
	db           *sql.DB
	itemSaver    ItemSaver
	itemProvider ItemProvider
}

func NewPostgresRepository(dataSourceName string, itemSaver ItemSaver, provider ItemProvider) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db, itemSaver, provider}, nil
}
func (p *PostgresRepository) CreateLoginItem(ctx context.Context, item domain.LoginItem) (uuid.UUID, error) {
	const op = "repositories.item.loginItem.postgres.CreateLoginItem"

	itemId, err := p.itemSaver.CreateItem(ctx, item.Item)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := p.db.Prepare("INSERT INTO login_items (id, item_id, login, encrypt_password) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING ID")
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, repositories.ErrItemExists)
		}
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	var id uuid.UUID
	row := stmt.QueryRowContext(ctx, itemId, item.Login, item.EncryptPassword)
	err = row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *PostgresRepository) GetLoginItem(ctx context.Context, loginItemId, userId uuid.UUID) (*domain.LoginItem, error) {
	const op = "repositories.loginItem.loginItem.postgres.GetLoginItem"

	stmt, err := p.db.Prepare("SELECT login_items.id, login, encrypt_password, items.id FROM login_items JOIN items on items.id = login_items.item_id WHERE login_items.id = $1 AND user_id=$2")
	if err != nil {
		return &domain.LoginItem{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, loginItemId, userId)

	var loginItem domain.LoginItem
	var itemId uuid.UUID
	err = row.Scan(&loginItem.ID, &loginItem.Login, &loginItem.EncryptPassword, &itemId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repositories.ErrItemNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	item, err := p.itemProvider.GetItem(ctx, itemId, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	loginItem.Item = *item

	return &loginItem, nil
}

func (p *PostgresRepository) GetLoginItems(ctx context.Context, userId uuid.UUID) ([]*domain.LoginItem, error) {
	const op = "repositories.item.loginItem.postgres.GetLoginItems"

	stmt, err := p.db.Prepare("SELECT login_items.id, login, encrypt_password, items.id FROM login_items JOIN items on items.id = login_items.item_id WHERE user_id=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var items []*domain.LoginItem
	for rows.Next() {
		var loginItem domain.LoginItem
		var itemId uuid.UUID
		err = rows.Scan(&loginItem.ID, &loginItem.Login, &loginItem.EncryptPassword, &itemId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		item, err := p.itemProvider.GetItem(ctx, itemId, userId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		loginItem.Item = *item
		items = append(items, &loginItem)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return items, nil
}
func (p *PostgresRepository) DeleteLoginItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	const op = "repositories.item.loginItem.postgres.DeleteLoginItem"

	stmt, err := p.db.Prepare("DELETE FROM items WHERE id=$1 AND user_id=$2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	stmt.QueryRowContext(ctx, itemId, userId)
	return nil
}
