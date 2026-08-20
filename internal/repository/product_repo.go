package repository

import (
	"errors"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/model"
	"gorm.io/gorm"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrVersionConflict = errors.New("version conflict")
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id uint) (*model.Product, error) {
	var p model.Product
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) UpdateQuantity(id uint, delta int, version int) error {
	result := r.db.Model(&model.Product{}).
		Where("id = ? AND version = ?", id, version).
		Updates(map[string]interface{}{
			"quantity": gorm.Expr("quantity + ?", delta),
			"version":  gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrVersionConflict
	}

	return nil
}
