package service

import (
	"context"
	"encoding/json"
	"log"
	"fmt"

	"pricetracker/api/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ItemService struct {
	repo *repository.ItemRepository
	ch   *amqp.Channel
	historyRepo *repository.HistoryRepository
}

func NewItemService(repo *repository.ItemRepository, historyRepo *repository.HistoryRepository, ch *amqp.Channel) *ItemService {
	return &ItemService{
		repo: repo,
		ch:   ch,
		historyRepo: historyRepo,
	}
}

func (s *ItemService) GetAllItems(ctx context.Context) ([]repository.Item, error) {
	items, err := s.repo.GetAllItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения элементов: %w",err)
	}
	return items, nil
} 

func (s *ItemService) ProcessItem(ctx context.Context, url string) (int, error) {
	
	type Task struct {
		Id  int    `json:"id"`
		Url string `json:"url"`
	}

	id, err := s.repo.Create(ctx, url)
	if err != nil {
		log.Printf("Ошибка: %v", err)
		return 0, err
	}

	t := Task {
		Id: id,
		Url: url,
	}

	bodyBytes, err := json.Marshal(t)
	if err != nil {
		log.Printf("Не удалось преобразовать структуру: %v", err)
		return 0, err
	}

	err = s.ch.PublishWithContext(
		ctx,
		"",
		"parsing_tasks",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: bodyBytes,
		},
	)
	if err != nil {
		log.Printf("Ошибка публикации: %v", err)
	}
	log.Println("Сообщение отправлено!")
	return id, nil

}

func (s *ItemService) GetHistory(ctx context.Context, id int) ([]repository.HistoryItem, error) {
	history, err := s.historyRepo.GetHistory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения истории: %w", err)
	}
	return history, nil
}