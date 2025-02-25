package handlers

import (
	"net/http"
	"rest-api/internal/domain"
	"rest-api/internal/lib/storage/postgre"
	"strconv"

	"github.com/gin-gonic/gin"
)

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