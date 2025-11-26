package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Balyshev/notes-api/internal/handlers"
	"github.com/Balyshev/notes-api/internal/middleware"
	"github.com/Balyshev/notes-api/internal/storage"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Подключаемся к БД
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("✅ Connected to database")

	// 3. Создаём storage
	store := storage.New(db)

	// 4. Создаём handlers
	authHandler := handlers.NewAuthHandler(store)
	userHandler := handlers.NewUserHandler(store)
	noteHandler := handlers.NewNoteHandler(store)

	// 5. Настраиваем роутер
	r := chi.NewRouter()

	// Middleware (применяются ко всем роутам)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/*", http.StripPrefix("/", fs))

	// Публичные роуты (без авторизации)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/users", userHandler.CreateUser) // Deprecated, использовать /auth/register

	// Защищённые роуты (требуют JWT токен)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Применяем JWT middleware

		// Роуты для заметок
		r.Post("/users/{id}/notes", noteHandler.CreateNote)
		r.Get("/users/{id}/notes", noteHandler.GetUserNotes)
		r.Get("/users/{id}/notes/{note_id}", noteHandler.GetNote)
		r.Put("/users/{id}/notes/{note_id}", noteHandler.UpdateNote)
		r.Delete("/users/{id}/notes/{note_id}", noteHandler.DeleteNote)
	})

	// 6. Запускаем сервер
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s...\n", port)
	fmt.Println("📝 Public endpoints:")
	fmt.Println("   POST /auth/register - Register new user")
	fmt.Println("   POST /auth/login - Login")
	fmt.Println("🔒 Protected endpoints (require JWT token):")
	fmt.Println("   POST   /users/{id}/notes")
	fmt.Println("   GET    /users/{id}/notes")
	fmt.Println("   GET    /users/{id}/notes/{note_id}")
	fmt.Println("   PUT    /users/{id}/notes/{note_id}")
	fmt.Println("   DELETE /users/{id}/notes/{note_id}")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDB инициализирует подключение к БД
func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
