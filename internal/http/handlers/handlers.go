package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Item – структура, описывающая сущность с полями ID, Name и Description.
// Теги json используются для правильной сериализации/десериализации данных.
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// items – in-memory хранилище для наших item-ов.
// nextID – переменная для генерации уникальных идентификаторов.
var items = map[int]Item{}
var nextID = 1

// getItems – обработчик для GET /items.
// Он собирает все item из хранилища и возвращает их в формате JSON.
func GetItemsHandler(c *gin.Context) {
	var itemList []Item
	for _, item := range items {
		itemList = append(itemList, item)
	}
	c.JSON(http.StatusOK, itemList)
}

// getItem – обработчик для GET /items/:id.
// Извлекает параметр id из URL, конвертирует его в число и возвращает item, если он найден.
func GetItemHandler(c *gin.Context) {
	idStr := c.Param("id")              // Получаем строковое значение параметра id.
	id, err := strconv.Atoi(idStr)       // Преобразуем строку в число.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	item, exists := items[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// createItem – обработчик для POST /items.
// Читает JSON-данные из тела запроса, создает новый item, присваивает ему уникальный ID и сохраняет его.
func CreateItemHandler(c *gin.Context) {
	var newItem Item
	// ShouldBindJSON пытается распарсить JSON из тела запроса в структуру newItem.
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newItem.ID = nextID  // Присваиваем уникальный ID.
	nextID++
	items[newItem.ID] = newItem
	c.JSON(http.StatusCreated, newItem)
}

// updateItem – обработчик для PUT /items/:id.
// Извлекает id из URL, парсит JSON из запроса, обновляет существующий item и возвращает обновленный объект.
func UpdateItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	var updatedItem Item
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := items[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	updatedItem.ID = id         // Гарантируем, что ID остается неизменным.
	items[id] = updatedItem       // Сохраняем обновленный item.
	c.JSON(http.StatusOK, updatedItem)
}

// deleteItem – обработчик для DELETE /items/:id.
// Извлекает id из URL, проверяет наличие item, удаляет его из хранилища и возвращает статус 204 (No Content).
func DeleteItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	if _, exists := items[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	delete(items, id)
	c.Status(http.StatusNoContent)
}