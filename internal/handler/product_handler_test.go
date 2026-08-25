package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/service"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	err error
}

func (m *mockService) UpdateQuantity(ctx context.Context, id uint, delta int) error {
	return m.err
}

func TestUpdateQuantity_Success(t *testing.T) {
	h := NewProductHandler(&mockService{err: nil})

	body := `{"delta": -1}`
	req := httptest.NewRequest("POST", "/products/1/update-quantity", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// Вызываем напрямую без Gin router
	h.UpdateQuantity(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ok", resp["status"])
}

func TestUpdateQuantity_Busy(t *testing.T) {
	h := NewProductHandler(&mockService{err: service.ErrBusy})

	body := `{"delta": -1}`
	req := httptest.NewRequest("POST", "/products/1/update-quantity", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.UpdateQuantity(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateQuantity_InvalidBody(t *testing.T) {
	h := NewProductHandler(&mockService{err: nil})

	req := httptest.NewRequest("POST", "/products/1/update-quantity", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.UpdateQuantity(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
