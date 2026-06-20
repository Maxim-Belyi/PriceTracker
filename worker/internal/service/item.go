package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"pricetracker/worker/internal/parser"
	"pricetracker/worker/internal/repository"
)

type ItemService struct {
	repo   *repository.ItemRepository
	parser parser.ItemParser
}

func NewItemService(repo *repository.ItemRepository, p parser.ItemParser) *ItemService {
	return &ItemService{
		repo:   repo,
		parser: p,
	}
}

func (s *ItemService) ProcessTask(ctx context.Context, messageBody []byte) error {
	type Task struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	}

	var t Task
	if err := json.Unmarshal(messageBody, &t); err != nil {
		log.Printf("Ошибка декодирования Json: %v", err)
		return err
	}
	
	log.Printf("Начинаю парсинг каталога: %v", t.URL)

	parsedItems, err := s.parser.Parse(ctx, t.URL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга: %w", err)
	}
	log.Printf("спарсили %d товаров. Начинаем запись в БД...", len(parsedItems))
	
	err = s.repo.SaveMultipleItems(ctx, parsedItems)
	if err != nil {
		return fmt.Errorf("ошибка массового сохранения: %w", err)
	}

	if err := s.repo.DeleteTask(ctx, t.ID); err != nil {
		log.Printf("не удалось удалить задачу ID=%d: %v", t.ID, err)
	}
	
	log.Printf("Каталог обработан, все товары в базе")
	return nil
}