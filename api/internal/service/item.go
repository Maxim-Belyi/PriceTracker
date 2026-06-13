package service

import (
	"context"
	"encoding/json"
	"log"

	"pricetracker/api/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ItemService struct {
	repo *repository.ItemRepository
	ch   *amqp.Channel
}

func NewItemService(repo *repository.ItemRepository, ch *amqp.Channel) *ItemService {
	return &ItemService{
		repo: repo,
		ch:   ch,
	}
}

func (s *ItemService) GetAllItems() ([]repository.Item, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Ошибка получения данных из репозитория: %v",err)
		return nil, err
	}
	return items, nil
} 

func (s *ItemService) ProcessItem(ctx context.Context, url string) (int, error) {
	
	type Task struct {
		Id  int    `json:"id"`
		Url string `json:"url"`
	}

	id, err := s.repo.Create(url)
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
