package postgre

import (
    "fmt"
    "log"
    "os"

    "github.com/jmoiron/sqlx" // Импортируем sqlx
    _ "github.com/lib/pq"      // Драйвер PostgreSQL для регистрации в database/sql
)

func main() {
    // Чтение переменных окружения для параметров подключения.
    // Это помогает не хардкодить данные, а задавать их вне кода.
    host, hostExists := os.LookupEnv("DB_HOST")
    port, portExists := os.LookupEnv("DB_PORT")
    user, userExists := os.LookupEnv("DB_USER")
    password, passExists := os.LookupEnv("DB_PASSWORD")
    dbname, dbnameExists := os.LookupEnv("DB_NAME")

    // Проверяем наличие всех необходимых переменных окружения.
    if !hostExists || !portExists || !userExists || !passExists || !dbnameExists {
        log.Fatal("Одна или несколько переменных окружения (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME) не установлены")
    }

    // Формирование DSN (Data Source Name) для подключения.
    // sslmode=disable используется для локальной разработки.
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Используем sqlx.Connect для установки подключения.
    // Функция Connect сразу выполняет проверку (ping) подключения.
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatalf("Ошибка подключения к базе данных: %v", err)
    }
    defer db.Close()

    fmt.Println("Успешное подключение к PostgreSQL через sqlx!")

    // Пример выполнения SQL-запроса: создание таблицы "items" (если она ещё не существует)
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR NOT NULL,
		description TEXT NOT NULL
    );`
    if _, err := db.Exec(createTableQuery); err != nil {
        log.Fatalf("Ошибка при создании таблицы: %v", err)
    }
    fmt.Println("Таблица 'items' успешно создана!")
}
