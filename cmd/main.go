// set CONFIG_PATH=config\local.yaml
// export DB_HOST=localhost
// export DB_PORT=5433
// export DB_USER=myuser
// export DB_PASSWORD=root
// export DB_NAME=restdb
// go run cmd/main.go
package main

import (
	"log"
	"rest-api/internal/http/handlers"
	"rest-api/internal/lib/storage/postgre"

	"github.com/gin-gonic/gin"
)

func main() {
    storage, err := postgre.New()
    if err != nil {
        log.Fatalf("failed to create storage: %v", err)
    }
    defer storage.Close() 

    router := gin.Default()
    // Передавайте storage в обработчики через контекст или синглтон
    router.GET("/items", handlers.GetItemsHandler(storage))
    router.GET("/items/:id", handlers.GetItemHandler(storage))
    router.POST("/items", handlers.CreateItemHandler(storage))
    router.PUT("/items/:id", handlers.UpdateItemHandler(storage))
    router.DELETE("/items/:id", handlers.DeleteItemHandler(storage))

    if err := router.Run(":8080"); err != nil {
        log.Fatalf("failed to run server: %v", err)
    }
}

