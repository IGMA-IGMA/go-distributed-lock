package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	UpdateQuantity(ctx context.Context, id uint, delta int) error
}

type ProductHandler struct {
	svc ProductService
}

func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

type UpdateQuantityRequest struct {
	Delta int `json:"delta" binding:"required"`
}

// GinHandler возвращает обработчик для Gin
func (h *ProductHandler) GinHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	h.updateQuantity(c.Writer, c.Request, uint(id), req.Delta)
}

// HTTPHandler возвращает обработчик для net/http
func (h *ProductHandler) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	h.updateQuantity(w, r, uint(id), req.Delta)
}

func (h *ProductHandler) updateQuantity(w http.ResponseWriter, r *http.Request, id uint, delta int) {
	err := h.svc.UpdateQuantity(r.Context(), id, delta)
	if err != nil {
		switch err {
		case service.ErrBusy:
			http.Error(w, err.Error(), http.StatusConflict)
		case service.ErrInsufficient:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case service.ErrVersionConflict:
			http.Error(w, err.Error(), http.StatusConflict)
		case service.ErrProductNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}