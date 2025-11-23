package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Balyshev/notes-api/internal/models"
)

// CreateNote создаёт новую заметку
func (s *Storage) CreateNote(userID int, title, content string) (*models.Note, error) {
	query := `
		INSERT INTO notes (user_id, title, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, user_id, title, content, created_at, updated_at
	`

	note := &models.Note{}
	err := s.db.QueryRow(query, userID, title, content).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return note, nil
}

// GetNoteByID получает заметку по ID
func (s *Storage) GetNoteByID(noteID int) (*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE id = $1
	`

	note := &models.Note{}
	err := s.db.QueryRow(query, noteID).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoteNotFound
		}
		return nil, err
	}

	return note, nil
}

// GetUserNotes получает все заметки пользователя с пагинацией и сортировкой
func (s *Storage) GetUserNotes(userID, limit, offset int, sortOrder string) ([]*models.Note, error) {
	// Проверяем sortOrder (защита от SQL injection)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // по умолчанию
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY created_at %s
		LIMIT $2 OFFSET $3
	`, sortOrder)

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// UpdateNote обновляет заметку
func (s *Storage) UpdateNote(noteID int, title, content string) (*models.Note, error) {
	query := `
		UPDATE notes
		SET title = $1, content = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, user_id, title, content, created_at, updated_at
	`

	note := &models.Note{}
	err := s.db.QueryRow(query, title, content, noteID).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoteNotFound
		}
		return nil, err
	}

	return note, nil
}

// DeleteNote удаляет заметку
func (s *Storage) DeleteNote(noteID int) error {
	query := `DELETE FROM notes WHERE id = $1`

	result, err := s.db.Exec(query, noteID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNoteNotFound
	}

	return nil
}
