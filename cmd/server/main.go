package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	ok, err := rdb.SetNX(ctx, "test:lock", "1", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("lock acquired:", ok)
}