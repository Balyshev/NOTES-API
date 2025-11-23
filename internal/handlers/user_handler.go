package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Balyshev/notes-api/internal/models"
	"github.com/Balyshev/notes-api/internal/storage"
)

// обрабатывает запросы к Users
type UserHandler struct {
	storage *storage.Storage
}

func NewUserHandler(storage *storage.Storage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

// обрабатывает POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== CreateUser called ===")

	// 1. Парсим JSON
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("ERROR: Failed to decode JSON:", err)
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	fmt.Printf("Parsed request: %+v\n", req)

	// 2. Валидируем
	if err := req.Validate(); err != nil {
		fmt.Println("ERROR: Validation failed:", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Validation passed")

	// 3. Создаём пользователя
	fmt.Printf("Calling storage.CreateUser with username: %s\n", req.Username)
	user, err := h.storage.CreateUser(req.Username)
	if err != nil {
		fmt.Println("ERROR: storage.CreateUser failed:", err) // ← ЭТО ГЛАВНОЕ!

		if err == models.ErrUsernameExists {
			respondError(w, http.StatusBadRequest, "Username already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	fmt.Printf("User created successfully: %+v\n", user)

	// 4. Возвращаем ответ
	respondJSON(w, http.StatusCreated, user)
}
