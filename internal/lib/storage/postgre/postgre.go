package postgre

import (
	"fmt"
	"os"
	"rest-api/internal/domain"

	"github.com/jmoiron/sqlx" // Импортируем sqlx
	_ "github.com/lib/pq"     // Драйвер PostgreSQL для регистрации в database/sql
)

type Storage struct {
	db *sqlx.DB
} 

func New() (*Storage, error){
    // Чтение переменных окружения для параметров подключения.
    // Это помогает не хардкодить данные, а задавать их вне кода.
    host, hostExists := os.LookupEnv("DB_HOST")
    port, portExists := os.LookupEnv("DB_PORT")
    user, userExists := os.LookupEnv("DB_USER")
    password, passExists := os.LookupEnv("DB_PASSWORD")
    dbname, dbnameExists := os.LookupEnv("DB_NAME")

    // Проверяем наличие всех необходимых переменных окружения.
    if !hostExists || !portExists || !userExists || !passExists || !dbnameExists {
        return &Storage{}, fmt.Errorf("failed to find env variables")
    }

    // Формирование DSN (Data Source Name) для подключения.
    // sslmode=disable используется для локальной разработки.
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    fmt.Println(dsn)
    // Используем sqlx.Connect для установки подключения.
    // Функция Connect сразу выполняет проверку (ping) подключения.
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return &Storage{}, fmt.Errorf("failed storage connection: %w", err)
    }
    fmt.Println("Успешное подключение к PostgreSQL через sqlx!")

    // Пример выполнения SQL-запроса: создание таблицы "items" (если она ещё не существует)
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS items(
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT NOT NULL);
    `

    if _, err := db.Exec(createTableQuery); err != nil {
        return &Storage{}, fmt.Errorf("failed to execute: %w", err)
    }
    return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
    return s.db.Close()
}

func (s *Storage) GetItems() ([]domain.Item, error){
    var items []domain.Item
    query := `SELECT id, name, description FROM items`

    if err := s.db.Select(&items, query); err != nil{
        return nil, fmt.Errorf("failed to get items: %w", err)
    }

    return items, nil
}

func (s *Storage) GetItem(id int) (domain.Item, error){
    var item domain.Item
    query := `SELECT id, name, description FROM items WHERE id = $1`
    if err := s.db.Get(&item, query, id); err != nil{
        return domain.Item{}, fmt.Errorf("failed to get item: %w", err)     
    }
    return item, nil
}

func (s *Storage) CreateItem(item domain.Item) (domain.Item, error){
    query := `INSERT INTO items(name, description) VALUES($1, $2) RETURNING id`
    err := s.db.QueryRow(query, item.Name, item.Description).Scan(&item.ID)
    if err != nil {
        return domain.Item{}, fmt.Errorf("failed to create item: %w", err)
    }
    return item, nil
}

func (s *Storage) UpdateItem(id int, item domain.Item) (domain.Item, error){
    query := `UPDATE items SET name = $1, description = $2 WHERE id = $3 RETURNING id`
    err := s.db.QueryRow(query, item.Name, item.Description, id).Scan(&item.ID)
    if err != nil {
        return domain.Item{}, fmt.Errorf("failed to update item: %w", err) 
        }            
    return item, nil
}

func (s *Storage) DeleteItem(id int) (domain.Item, error){
    var deletedItem domain.Item
    query := `SELECT id, name, description FROM items WHERE id = $1`
    if err := s.db.Get(&deletedItem, query, id); err != nil{
        return domain.Item{}, fmt.Errorf("failed to get item: %w", err)
    }

    query = "DELETE FROM items WHERE id = $1"
    if _, err := s.db.Exec(query, id); err != nil{
        return domain.Item{}, fmt.Errorf("failed to delete item: %w", err)     
    }
    return deletedItem, nil
}