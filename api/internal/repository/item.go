package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ItemRepository struct {
	db *sql.DB
}

type Item struct {
	ID           int     `json:"id"`
	Title        *string `json:"title"`
	ImageUrl     *string `json:"image_url"`
	CurrentPrice float64 `json:"current_price"`
	Status       string  `json:"status"`
}

type HistoryRepository struct {
	db *sql.DB
}

type HistoryItem struct {
	ID int `json:"id"`
	Price float64 `json:"price"`
}

func (r *ItemRepository) GetAllItems() ([]Item, error) {
    query := `SELECT id, title, image_url, current_price, status
              FROM items
              ORDER BY updated_at DESC`
     
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err 
    }
    defer rows.Close() 

    var items []Item 

    for rows.Next() {
        var i Item
        err := rows.Scan(&i.ID, &i.Title, &i.ImageUrl, &i.CurrentPrice, &i.Status)
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

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) Create(url string) (int, error) {
	var id int
	query := (`INSERT INTO items (url) VALUES ($1) RETURNING id`)
	if err := r.db.QueryRow(query, url).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
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