package handler

import (
	"log"
	"net/http"
	"pricetracker/api/internal/service"
	"encoding/json"
)

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
