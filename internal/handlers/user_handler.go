package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Balyshev/notes-api/internal/models"
	"github.com/Balyshev/notes-api/internal/storage"
	"github.com/Balyshev/notes-api/pkg/auth"
)

// UserHandler обрабатывает запросы к /users
type UserHandler struct {
	storage *storage.Storage
}

// NewUserHandler создаёт новый UserHandler
func NewUserHandler(storage *storage.Storage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

// CreateUser теперь требует пароль (используем AuthHandler.Register вместо этого)
// Оставляем для обратной совместимости, но лучше использовать /auth/register
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== CreateUser called (deprecated, use /auth/register) ===")

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Хешируем пароль
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user, err := h.storage.CreateUser(req.Username, passwordHash)
	if err != nil {
		fmt.Println("ERROR: storage.CreateUser failed:", err)

		if err == models.ErrUsernameExists {
			respondError(w, http.StatusBadRequest, "Username already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	fmt.Printf("User created successfully: %+v\n", user)
	respondJSON(w, http.StatusCreated, user)
}
