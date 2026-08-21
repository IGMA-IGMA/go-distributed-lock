package lock

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestDistributedLock(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	ctx := context.Background()

	lock1 := NewDistributedLock(client, "test:lock", 10*time.Second)
	lock2 := NewDistributedLock(client, "test:lock", 10*time.Second)

	// первый должен захватить
	ok, err := lock1.TryLock(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected lock to be acquired")
	}

	// второй не должен
	ok, err = lock2.TryLock(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected lock to be busy")
	}

	// освобождаем
	err = lock1.Unlock(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// теперь второй может захватить
	ok, err = lock2.TryLock(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected lock to be acquired after unlock")
	}
}
