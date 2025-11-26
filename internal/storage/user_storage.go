package storage

import (
	"database/sql"
	"errors"

	"github.com/Balyshev/notes-api/internal/models"
	"github.com/lib/pq"
)

// CreateUser создаёт нового пользователя
func (s *Storage) CreateUser(username, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, username, created_at
	`

	user := &models.User{}
	err := s.db.QueryRow(query, username, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		// Проверяем, не дубликат ли username
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return nil, models.ErrUsernameExists
			}
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUsername получает пользователя по username (для логина)
func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash, // Теперь получаем хеш пароля
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID получает пользователя по ID
func (s *Storage) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
