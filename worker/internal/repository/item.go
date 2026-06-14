package repository

import (
	"database/sql"
	"fmt"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) UpdateItemData(id int, title string, imageUrl string, price float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Ошибка старта транзакции: %v", err)
	}
	defer tx.Rollback()

	updateQuery := `
	UPDATE items 
	SET current_price = $1,
		title = $2,
		image_url = $3,
	status = 'processed',
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $4
	`
	if _, err := tx.Exec(updateQuery, price, title, imageUrl, id); err != nil {
		return fmt.Errorf("Ошибка обновления items: %v", err)
	}

	historyQuery := `
	INSERT INTO price_history (item_id, price)
	VALUES($1, $2)
	`
	if _, err := tx.Exec(historyQuery, id, price); err != nil {
		return fmt.Errorf("Ошибка записи истории цены: %v", err)
	}
	return tx.Commit()
}
