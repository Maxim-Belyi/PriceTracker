package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ItemRepository struct {
	db *sql.DB
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
