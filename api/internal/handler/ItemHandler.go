package handler

import (
	"net/http"
	"pricetracker/api/internal/service"
	"encoding/json"
)

type TrackRequest struct {
	Url string `json:"url"`
}

type TrackResponse struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type ItemHandler struct {
	service *service.ItemService
}

func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	allItems, err := h.service.GetAllItems()
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allItems); err != nil {
		http.Error(w, "Ошибка кодирования Json", http.StatusInternalServerError)
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

	id, err := h.service.ProcessItem(r.Context(), req.Url)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	res := TrackResponse  {
		Id: id,
		Status: "Сохранено",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Ошибка декодирования Json", http.StatusInternalServerError)
	}
}