package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Balyshev/notes-api/internal/handlers"
	"github.com/Balyshev/notes-api/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//загружаем переменные .env в окружение ОС
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
	//подключаемся к БД
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	defer db.Close()
	fmt.Println("Connected to db")

	//создаём storage подключение к БД
	store := storage.New(db)

	//создаём Handlers
	userHandler := handlers.NewUserHandler(store)

	//Настраиваем роутер
	r := chi.NewRouter()

	//middleware для логирования
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//роуты для пользователей
	r.Post("/users", userHandler.CreateUser)

	//запускаем сервер
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server start on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

// инициализируем подключение к БД
func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port =%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
