package handler

import (
	"encoding/json"
	"errors"
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Delta int `json:"delta"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.svc.UpdateQuantity(r.Context(), uint(id), req.Delta)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBusy):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, service.ErrInsufficient):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
