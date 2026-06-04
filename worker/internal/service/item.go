package service

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"pricetracker/worker/internal/repository"
	"time"
)

type ItemService struct {
	repo *repository.ItemRepository
}

func NewItemService(repo *repository.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) ProcessTask(messageBody []byte) error {

	type Task struct {
		Id  int    `json:"id"`
		Url string `json:"url"`
	}

	var t Task

	if err := json.Unmarshal(messageBody, &t); err != nil {
		log.Printf("Ошибка декодирования Json: %v", err)
		return err
	}
	log.Printf("Начинаю парсинг для URL: %v", t.Url)
	time.Sleep(2 *time.Second)
	price := float64(rand.IntN(1000)) + 100

	err := s.repo.UpdatePrice(price, t.Id); if err != nil {
		log.Printf("Ошибка обновления БД: %v", err)
		return err
	}
	return nil

}
