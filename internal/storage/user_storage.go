package storage

import (
	"database/sql"
	"errors"

	"github.com/Balyshev/notes-api/internal/models"
	"github.com/lib/pq"
)

func (s *Storage) CreateUser(username string) (*models.User, error) {
	query := `
		INSERT INTO users (username, created_at)
		VALUES ($1, NOW())
		RETURNING id,username,created_at
	`

	user := &models.User{}
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, models.ErrUsernameExists
			}
		}
		return nil, err
	}
	return user, nil
}

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

func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := s.db.QueryRow(query, username).Scan(
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
