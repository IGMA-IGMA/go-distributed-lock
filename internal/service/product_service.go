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
	ErrBusy            = errors.New("resource is busy")
	ErrInsufficient    = errors.New("insufficient quantity")
	ErrVersionConflict = errors.New("version conflict")
	ErrProductNotFound = errors.New("product not found")
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
		if errors.Is(err, repository.ErrNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	if product.Quantity+delta < 0 {
		return ErrInsufficient
	}

	lockKey := fmt.Sprintf("lock:product:%d", id)
	distributedLock := lock.NewDistributedLock(s.rdb, lockKey, 10*time.Second)

	acquired, err := distributedLock.TryLock(ctx)
	if err != nil {
		return err
	}

	if !acquired {
		return ErrBusy
	}

	defer distributedLock.Unlock(context.Background())

	err = s.repo.UpdateQuantity(id, delta, product.Version)
	if err != nil {
		if errors.Is(err, repository.ErrVersionConflict) {
			return ErrVersionConflict
		}
		return err
	}

	return nil
}
