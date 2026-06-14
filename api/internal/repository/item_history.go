package repository

import (
	"database/sql"

)

type HistoryRepository struct {
	db *sql.DB
}

type HistoryItem struct {
	ID int `json:"id"`
	Price float64 `json:"price"`
}

func (h *HistoryRepository) GetHistory(itemID int) ([]HistoryItem, error) {
	query := `SELECT id, price
			  FROM price_history
			  WHERE item_id = $1
			  ORDER BY created_at DESC`
	
	rows, err := h.db.Query(query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []HistoryItem

	for rows.Next() {
		var i HistoryItem
		err := rows.Scan(&i.ID, &i.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}