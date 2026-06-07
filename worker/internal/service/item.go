package service

import (
	"encoding/json"
	"log"
	"net/http"
	"pricetracker/worker/internal/parser"
	"pricetracker/worker/internal/repository"
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

		resp, err := http.Get(t.Url)
		if err != nil{
			log.Printf("Ошибка скачивания страницы: %v", err)
			return err
		}
		log.Printf("Страница скачана! Статус: %d", resp.StatusCode)
		defer resp.Body.Close()

	p := parser.NewCitilinkParser()
	
	parsedItem, err := p.Parse(t.Url)
	if err != nil {
		log.Printf("Ошибка парсинга: %v", err)
	}
	
	err = s.repo.UpdateItemData(t.Id, parsedItem.Title, parsedItem.ImageUrl, parsedItem.Price)
	if err != nil {
		log.Printf("Ошибка обновления БД: %v", err)
		return err
	}
	return nil

}
