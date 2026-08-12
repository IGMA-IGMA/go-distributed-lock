package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/lock"
	"github.com/IGMA-IGMA/go-distributed-lock/internal/repository"
	"github.com/redis/go-redis/v9"
)

var (
	ErrBusy         = errors.New("resource is busy")
	ErrInsufficient = errors.New("insufficient quantity")
)

type ProductService struct {
	repo *repository.ProductRepository
	rdb  *redis.Client
}

func NewProductService(repo *repository.ProductRepository, rdb *redis.Client) *ProductService {
	return &ProductService{repo: repo, rdb: rdb}
}

func (s *ProductService) UpdateQuantity(ctx context.Context, id uint, delta int) error {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if product.Quantity+delta < 0 {
		return ErrInsufficient
	}

	lock := lock.NewDistributedLock(s.rdb, fmt.Sprintf("lock:product:%d", id), 10*time.Second)

	acquired, err := lock.TryLock(ctx)
	if err != nil {
		return err
	}

	if !acquired {
		return ErrBusy
	}

	defer lock.Unlock(context.Background())

	return s.repo.UpdateQuantity(id, delta)
}
