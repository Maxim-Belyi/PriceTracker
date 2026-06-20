package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pricetracker/worker/internal/parser"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) DeleteTask(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM items WHERE id = $1", id)
	return err
}

func (r *ItemRepository) SaveMultipleItems(ctx context.Context, items []parser.ParsedItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка старта транзакции: %v", err)
	}
	defer tx.Rollback()

	insertItemQuery := `
	INSERT INTO items (url, title, image_url, current_price, status, updated_at)
	VALUES ($1, $2, $3, $4, 'processed', CURRENT_TIMESTAMP)
	RETURNING id
	`
	insertHistoryQuery := `
	INSERT INTO price_history (item_id, price)
	VALUES($1, $2)
	`

	for _, item := range items {
		var newID int
		err := tx.QueryRowContext(ctx, insertItemQuery, item.ProductURL, item.Title, item.ImageURL, item.Price).Scan(&newID)
		if err != nil {
			return fmt.Errorf("ошибка вставки товара '%s': %v", item.Title, err)
		}

		if _, err := tx.ExecContext(ctx, insertHistoryQuery, newID, item.Price); err != nil {
			return fmt.Errorf("ошибка записи истории цены для '%s': %v", item.Title, err)
		}
	}

	return tx.Commit()
}
