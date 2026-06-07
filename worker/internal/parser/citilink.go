package parser

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type CitilinkParser struct{}

func NewCitilinkParser() *CitilinkParser{
	return &CitilinkParser{}
}

func (p *CitilinkParser) Parse(url string) (ParsedItem, error) {
	var item ParsedItem

	resp, err := http.Get(url)
	if err != nil {
		return item, fmt.Errorf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return item, fmt.Errorf("Сервис вернул статус %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return item, fmt.Errorf("Ошибка парсинга Html: %v", err)
	}
	item.Title = doc.Find(`a[data-meta-name="Snippet__title"]`).AttrOr("title", "Без названия")

	item.ImageUrl = doc.Find(`div[data-meta-name="Snippet__images"] [img]`).AttrOr("src", "")

	priceStr := doc.Find(`span[data-meta-name="Snippet__price"]`).AttrOr("data-meta-price", "")

	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("Не удалось перевести строку в число: %v", err)
	} else {
		item.Price = priceFloat
	}

	return item, nil
}