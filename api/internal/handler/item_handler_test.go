package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pricetracker/api/internal/repository"
	"strings"
	"testing"
)

type mockItemService struct {
	getAllItemsFn  func(ctx context.Context) ([]repository.Item, error)
	processItemFn func(ctx context.Context, url string) (int, error)
	getHistoryFn  func(ctx context.Context, id int) ([]repository.HistoryItem, error)
}


func (m *mockItemService) GetAllItems(ctx context.Context) ([]repository.Item, error) {
	return m.getAllItemsFn(ctx)
}

func (m *mockItemService) ProcessItem(ctx context.Context, url string) (int, error) {
	return m.processItemFn(ctx, url)
}

func (m *mockItemService) GetHistory(ctx context.Context, id int) ([]repository.HistoryItem, error) {
	return m.getHistoryFn(ctx, id)
}

func ptr(s string) *string { return &s }

func TestGetAll(t *testing.T) {
	t.Run("успешный ответ возвращает товары в JSON", func(t *testing.T) {
		svc := &mockItemService{
			getAllItemsFn: func(ctx context.Context) ([]repository.Item, error) {
				return []repository.Item{
					{ID: 1, Title: ptr("Смартфон"), CurrentPrice: 50000, Status: "processed", Source: "citilink"},
					{ID: 2, Title: ptr("Ноутбук"), CurrentPrice: 90000, Status: "processed", Source: "citilink"},
				}, nil
			},
		}

		h := NewItemHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/items", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("ожидали 200, получили %d", resp.StatusCode)
		}

		ct := resp.Header.Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("ожидали Content-Type application/json, получили %q", ct)
		}

		var items []repository.Item
		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			t.Fatalf("не удалось декодировать ответ: %v", err)
		}

		if len(items) != 2 {
			t.Errorf("ожидали 2 товара, получили %d", len(items))
		}
	})

	t.Run("ошибка сервиса возвращает 500", func(t *testing.T) {
		svc := &mockItemService{
			getAllItemsFn: func(ctx context.Context) ([]repository.Item, error) {
				return nil, errors.New("connection refused")
			},
		}

		h := NewItemHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/items", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("ожидали 500, получили %d", w.Code)
		}
	})
}

func TestTrack(t *testing.T) {
	t.Run("валидный URL citilink → 201 Created", func(t *testing.T) {
		svc := &mockItemService{
			processItemFn: func(ctx context.Context, url string) (int, error) {
				return 42, nil 
			},
		}

		h := NewItemHandler(svc)

		body := `{"url": "https://www.citilink.ru/catalog/smartfony/"}`
		req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Track(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("ожидали 201, получили %d", w.Code)
		}

		var resp TrackResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("не удалось декодировать ответ: %v", err)
		}

		if resp.ID != 42 {
			t.Errorf("ожидали ID=42, получили %d", resp.ID)
		}
	})

	t.Run("неподдерживаемый домен → 400 Bad Request", func(t *testing.T) {
		svc := &mockItemService{} 
		h := NewItemHandler(svc)

		body := `{"url": "https://avito.ru/item/123"}`
		req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Track(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("ожидали 400, получили %d", w.Code)
		}
	})

	t.Run("сломанный JSON → 400 Bad Request", func(t *testing.T) {
		svc := &mockItemService{}
		h := NewItemHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewBufferString(`{invalid json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Track(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("ожидали 400, получили %d", w.Code)
		}
	})

	t.Run("пустое тело → 400 Bad Request", func(t *testing.T) {
		svc := &mockItemService{}
		h := NewItemHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewBufferString(`{}`))
		w := httptest.NewRecorder()

		h.Track(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("ожидали 400, получили %d (URL пустой — не поддерживается)", w.Code)
		}
	})
}
