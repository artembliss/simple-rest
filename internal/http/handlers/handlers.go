package handlers

import (
	"net/http"
	"rest-api/internal/domain"
	"rest-api/internal/lib/storage/postgre"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Он собирает все item из хранилища и возвращает их в формате JSON.
func GetItemsHandler(s *postgre.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := s.GetItems()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// Извлекает параметр id из URL, конвертирует его в число и возвращает item, если он найден.
func GetItemHandler(s *postgre.Storage) gin.HandlerFunc{
	return func(c *gin.Context) {
		var item domain.Item

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		item, err = s.GetItem(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	}
	}



// Читает JSON-данные из тела запроса, создает новый item, присваивает ему уникальный ID и сохраняет его.
func CreateItemHandler(s *postgre.Storage) gin.HandlerFunc{
	return func(c *gin.Context){
		var item domain.Item
        if err := c.ShouldBindJSON(&item); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        createdItem, err := s.CreateItem(item)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, createdItem)
	} 
}

// Извлекает id из URL, парсит JSON из запроса, обновляет существующий item и возвращает обновленный объект.
func UpdateItemHandler(s *postgre.Storage) gin.HandlerFunc{
	return func(c *gin.Context){
		var item domain.Item

		idStr :=  c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		err = c.ShouldBindJSON(&item)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		updatedItem, err := s.UpdateItem(id, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedItem)
		}
	}



// Извлекает id из URL, проверяет наличие item, удаляет его из хранилища и возвращает статус 204 (No Content).
func DeleteItemHandler(s *postgre.Storage) gin.HandlerFunc{
	return func(c *gin.Context){
		var deletedItem domain.Item
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}
		deletedItem, err = s.DeleteItem(id)	
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, deletedItem)
	}
}