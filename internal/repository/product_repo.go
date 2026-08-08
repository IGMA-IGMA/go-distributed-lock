package repository

import (
	"github.com/IGMA-IGMA/go-distributed-lock/internal/model"
	"gorm.io/gorm"
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
	return &p, err
}

func (r *ProductRepository) UpdateQuantity(id uint, delta int) error {
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		Update("quantity", gorm.Expr("quantity + ?", delta)).Error
}
