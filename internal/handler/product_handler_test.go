package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	err error
}

func (m *mockService) UpdateQuantity(ctx context.Context, id uint, delta int) error {
	return m.err
}

func TestUpdateQuantity_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{err: nil}
	h := NewProductHandler(mock)

	router := gin.New()
	router.POST("/products/:id/update-quantity", h.UpdateQuantity)

	body := `{"delta": -1}`
	req, _ := http.NewRequest("POST", "/products/1/update-quantity", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ok", resp["status"])
}

func TestUpdateQuantity_Busy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{err: service.ErrBusy}
	h := NewProductHandler(mock)

	router := gin.New()
	router.POST("/products/:id/update-quantity", h.UpdateQuantity)

	body := `{"delta": -1}`
	req, _ := http.NewRequest("POST", "/products/1/update-quantity", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateQuantity_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{err: nil}
	h := NewProductHandler(mock)

	router := gin.New()
	router.POST("/products/:id/update-quantity", h.UpdateQuantity)

	req, _ := http.NewRequest("POST", "/products/1/update-quantity", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
