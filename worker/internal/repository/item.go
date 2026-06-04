package repository
import (
	"database/sql"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository {
		db: db,
	}
}

func (r *ItemRepository) UpdatePrice(price float64, id int) (error) {
	query := `
	UPDATE items SET current_price = $1,
	status = 'processed',
	updated_at = CURRENT_TIMESTAMP WHERE id = $2
	`

	_, err := r.db.Exec(query, price, id)
	if err != nil {
		return err
	}
	return nil
	
}