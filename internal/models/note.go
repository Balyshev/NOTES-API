package models

import "time"

//Данные заметки
type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//createNoteRequest - данные для создания заметки
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `jsson:"content"`
}

//updateNoteRequest - данные для обновления заметки
type UdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

//Validate проверяет createNoteRequest
func (c *CreateNoteRequest) Validate() error {
	if c.Title == "" {
		return ErrTitleRequired
	}
	if len(c.Title) > 255 {
		return ErrTitleTooLong
	}
	if c.Content == "" {
		return ErrContentRequired
	}
	return nil
}

func (r *UdateNoteRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleRequired
	}
	if len(r.Title) > 255 {
		return ErrTitleTooLong
	}
	if r.Content == "" {
		return ErrContentRequired
	}
	return nil
}
