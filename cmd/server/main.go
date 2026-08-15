package main

import (
	"log"
	"net/http"

	"github.com/IGMA-IGMA/go-distributed-lock/internal/database"
	"github.com/IGMA-IGMA/go-distributed-lock/internal/handler"
	"github.com/IGMA-IGMA/go-distributed-lock/internal/model"
	"github.com/IGMA-IGMA/go-distributed-lock/internal/repository"
	"github.com/IGMA-IGMA/go-distributed-lock/internal/service"
	"github.com/redis/go-redis/v9"
)

func main() {
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Product{})

	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count == 0 {
		db.Create(&model.Product{Name: "Laptop", Quantity: 10})
		db.Create(&model.Product{Name: "Phone", Quantity: 20})
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo, rdb)
	h := handler.NewProductHandler(svc)

	http.HandleFunc("/products/update-quantity", h.UpdateQuantity)

	log.Println("server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
