package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Balyshev/notes-api/internal/models"
	"github.com/Balyshev/notes-api/internal/storage"
	"github.com/Balyshev/notes-api/pkg/auth"
)

// AuthHandler обрабатывает авторизацию
type AuthHandler struct {
	storage *storage.Storage
}

// NewAuthHandler создаёт новый AuthHandler
func NewAuthHandler(storage *storage.Storage) *AuthHandler {
	return &AuthHandler{
		storage: storage,
	}
}

// Register обрабатывает POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== Register called ===")

	// 1. Парсим JSON
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// 2. Валидируем
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Хешируем пароль
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		fmt.Println("ERROR: Failed to hash password:", err)
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// 4. Создаём пользователя
	user, err := h.storage.CreateUser(req.Username, passwordHash)
	if err != nil {
		if err == models.ErrUsernameExists {
			respondError(w, http.StatusBadRequest, "Username already exists")
			return
		}
		fmt.Println("ERROR: Failed to create user:", err)
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// 5. Генерируем JWT токен
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		fmt.Println("ERROR: Failed to generate token:", err)
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 6. Возвращаем токен и пользователя
	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	fmt.Printf("User registered: %+v\n", user)
	respondJSON(w, http.StatusCreated, response)
}

// Login обрабатывает POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== Login called ===")

	// 1. Парсим JSON
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// 2. Валидируем
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Получаем пользователя по username
	user, err := h.storage.GetUserByUsername(req.Username)
	if err != nil {
		if err == models.ErrUserNotFound {
			respondError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		fmt.Println("ERROR: Failed to get user:", err)
		respondError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	// 4. Проверяем пароль
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		respondError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// 5. Генерируем JWT токен
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		fmt.Println("ERROR: Failed to generate token:", err)
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 6. Возвращаем токен и пользователя
	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	fmt.Printf("User logged in: %s\n", user.Username)
	respondJSON(w, http.StatusOK, response)
}
