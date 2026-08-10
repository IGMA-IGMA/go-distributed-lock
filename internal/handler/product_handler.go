package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var req struct {
		Delta int `json:"delta"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := h.svc.UpdateQuantity(r.Context(), uint(id), req.Delta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
