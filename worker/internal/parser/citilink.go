package parser

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings" 

	"github.com/PuerkitoBio/goquery"
)

type CitilinkParser struct{}

func NewCitilinkParser() *CitilinkParser {
	return &CitilinkParser{}
}

func (p *CitilinkParser) Parse(ctx context.Context, url string) ([]ParsedItem, error) {
	var items []ParsedItem

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Сервис вернул статус %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Ошибка парсинга Html: %v", err)
	}

	titles := doc.Find(`a[data-meta-name="Snippet__title"]`)
	prices := doc.Find(`span[data-meta-name="Snippet__price"]`)
	images := doc.Find(`div[data-meta-name="Snippet__images"] [img]`)

	titles.Each(func(i int, s *goquery.Selection) {
		title := s.AttrOr("title", "Без названия")
		
		// Достаем ссылку на сам товар
		href := s.AttrOr("href", "")
		if href != "" && strings.HasPrefix(href, "/") {
			href = "https://www.citilink.ru" + href
		}

		priceStr := prices.Eq(i).AttrOr("data-meta-price", "0")
		imgUrl := images.Eq(i).AttrOr("src", "")

		priceFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			priceFloat = 0
		}

		items = append(items, ParsedItem{
			Title:      title,
			Price:      priceFloat,
			ImageURL:   imgUrl,
			ProductURL: href, 
		})
	})

	if len(items) == 0 {
		return nil, fmt.Errorf("на странице не найдены товары")
	}

	return items, nil
}