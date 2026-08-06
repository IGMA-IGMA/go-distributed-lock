package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type DistributedLock struct {
	client *redis.Client
	key    string
	value  string
	ttl    time.Duration
}

func NewDistributedLock(client *redis.Client, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		client: client,
		key:    key,
		value:  "1",
		ttl:    ttl,
	}
}

func (l *DistributedLock) TryLock(ctx context.Context) (bool, error) {
	return l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
}

func (l *DistributedLock) Unlock(ctx context.Context) error {
	return l.client.Del(ctx, l.key).Err()
}