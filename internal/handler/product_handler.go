package handler

import (
	"net/http"
	"strconv"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

type UpdateQuantityRequest struct {
	Delta int `json:"delta" binding:"required"`
}

func (h *ProductHandler) UpdateQuantity(c *gin.Context) {
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

	err = h.svc.UpdateQuantity(c.Request.Context(), uint(id), req.Delta)
	if err != nil {
		switch err {
		case service.ErrBusy:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case service.ErrInsufficient:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case service.ErrVersionConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
