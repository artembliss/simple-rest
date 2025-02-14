package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine{
		// Переводим Gin в режим тестирования для подавления логов.
		gin.SetMode(gin.TestMode)
		router := gin.Default()
	
		// Регистрируем те же маршруты, что и в основном приложении.
		router.GET("/items", getItemsHandler)
		router.GET("/items/:id", getItemHandler)
		router.POST("/items", createItemHandler)
		router.PUT("/items/:id", updateItemHandler)
		router.DELETE("/items/:id", deleteItemHandler)
		return router
}


// TestCreateAndGetItem тестирует создание item через POST и получение через GET.
func TestCreateAndGetItem(t *testing.T) {
	// Инициализируем роутер для тестов.
	router := setUpRouter()

	// Создаем новый item для теста.
	newItem := Item{Name: "Test Item", Description: "Test Description"}
	jsonValue, err := json.Marshal(newItem)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	// Создаем POST-запрос к /items с телом, содержащим JSON.
	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("Error creating POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Используем httptest.NewRecorder для имитации http.ResponseWriter.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что статус ответа равен http.StatusCreated (201).
	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status %d but got %d", http.StatusCreated, resp.Code)
	}

	// Декодируем ответ сервера в структуру Item.
	var createdItem Item
	err = json.Unmarshal(resp.Body.Bytes(), &createdItem)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	// Проверяем, что созданный item получил непустой ID.
	if createdItem.Id == 0 {
		t.Fatal("Expected non-zero item ID")
	}

	// Тестируем получение созданного item через GET-запрос.
	getURL := "/items/" + strconv.Itoa(createdItem.Id)
	req, err = http.NewRequest("GET", getURL, nil)
	if err != nil {
		t.Fatalf("Error creating GET request: %v", err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Проверяем, что статус ответа равен http.StatusOK (200).
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, resp.Code)
	}
}