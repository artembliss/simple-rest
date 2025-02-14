// set CONFIG_PATH=config\local.yaml
// go run cmd/main.go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Item struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
} 

var items = map[int]Item{}
var nextId = 1

func main(){
	router := gin.Default()
	router.GET("/items", getItemsHandler)
	router.GET("/items/{id}", getItemHandler)
	router.POST("/items", createItemHandler)
	router.PUT("/items/{id}", updateItemHandler)
	router.DELETE("/items/{id}", deleteItemHandler)


	router.Run(":8080")
}

func getItemsHandler(c *gin.Context){
	var itemsList []Item
	for _, item := range items{
		itemsList = append(itemsList, item)
	}
	c.JSON(http.StatusOK, itemsList)
}

func getItemHandler(c *gin.Context){
	itemId := c.Param("id")
	id, err := strconv.Atoi(itemId)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	item, exists := items[id] 
	if !exists{
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func createItemHandler(c *gin.Context){
	var newItem Item
	if err := c.ShouldBindJSON(&newItem); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return	
	}
	newItem.Id = nextId
	items[nextId] = newItem
	nextId++
	c.JSON(http.StatusCreated, newItem)
}

func updateItemHandler(c *gin.Context){
	var updatedItem Item
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := items[id]; !exists{
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	if err := c.ShouldBindJSON(&updatedItem); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	updatedItem.Id = id
	items[id] = updatedItem
	c.JSON(http.StatusOK, updatedItem)
}

func deleteItemHandler(c *gin.Context){
	deleteId, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := items[deleteId]; !exists{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	delete(items, deleteId)
	c.Status(http.StatusNoContent)
}
