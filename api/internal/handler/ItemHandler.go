package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pricetracker/api/internal/repository"
	"strconv"
)

type ItemService interface {
	GetAllItems(ctx context.Context) ([]repository.Item, error)
	ProcessItem(ctx context.Context, url string) (int, error)
	GetHistory(ctx context.Context, id int) ([]repository.HistoryItem, error)
}

type TrackRequest struct {
	URL string `json:"url"`
}

type TrackResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type ItemHandler struct {
	svc ItemService
}

func NewItemHandler(svc ItemService) *ItemHandler {
	return &ItemHandler{
		svc: svc,
	}
}

func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	allItems, err := h.svc.GetAllItems(r.Context())
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allItems); err != nil {
		log.Printf("Ошибка кодирования Json: %v", err)
	}
}

func (h *ItemHandler) Track(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	var req TrackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id, err := h.svc.ProcessItem(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	res := TrackResponse{
		ID:     id,
		Status: "Сохранено",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Ошибка кодирования ответа от Track: %v", err)
	}
}

func (h *ItemHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID товара", http.StatusBadRequest)
		return
	}
	history, err := h.svc.GetHistory(r.Context(), id)
	if err != nil {
		http.Error(w, "Ошибка получения истории", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		log.Printf("Ошибка кодирования JSON", err)
	}
}
