// set CONFIG_PATH=config\local.yaml
// go run cmd/main.go
package main

import (
	"rest-api/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()
	router.GET("/items", handlers.GetItemsHandler)
	router.GET("/items/{id}", handlers.GetItemHandler)
	router.POST("/items", handlers.CreateItemHandler)
	router.PUT("/items/{id}", handlers.UpdateItemHandler)
	router.DELETE("/items/{id}", handlers.DeleteItemHandler)

	router.Run(":8080")
}

