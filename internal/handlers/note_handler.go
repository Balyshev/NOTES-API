package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Balyshev/notes-api/internal/middleware"
	"github.com/Balyshev/notes-api/internal/models"
	"github.com/Balyshev/notes-api/internal/storage"
	"github.com/go-chi/chi/v5"
)

// NoteHandler обрабатывает запросы к /users/{id}/notes
type NoteHandler struct {
	storage *storage.Storage
}

func NewNoteHandler(storage *storage.Storage) *NoteHandler {
	return &NoteHandler{
		storage: storage,
	}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== CreateNote called ===")

	// Получаем user_id из JWT токена (из контекста)
	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получаем user_id из URL
	userIDFromURL, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Проверяем, что пользователь создаёт заметку для себя
	if authenticatedUserID != userIDFromURL {
		respondError(w, http.StatusForbidden, "You can only create notes for yourself")
		return
	}

	fmt.Printf("Authenticated UserID: %d\n", authenticatedUserID)

	// Парсим JSON
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	fmt.Printf("Parsed request: %+v\n", req)

	// Валидируем
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Создаём заметку (используем authenticatedUserID из токена, а не из URL!)
	note, err := h.storage.CreateNote(authenticatedUserID, req.Title, req.Content)
	if err != nil {
		fmt.Println("ERROR: CreateNote failed:", err)
		respondError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	fmt.Printf("Note created: %+v\n", note)
	respondJSON(w, http.StatusCreated, note)
}

// GetUserNotes обрабатывает GET /users/{id}/notes
func (h *NoteHandler) GetUserNotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== GetUserNotes called ===")

	// Получаем user_id из JWT токена
	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получаем user_id из URL
	userIDFromURL, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Проверяем, что пользователь запрашивает свои заметки
	if authenticatedUserID != userIDFromURL {
		respondError(w, http.StatusForbidden, "You can only view your own notes")
		return
	}

	// Парсим query параметры
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	sortOrder := r.URL.Query().Get("sort")

	limit := 10
	offset := 0
	if sortOrder == "" {
		sortOrder = "desc"
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			respondError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid offset parameter")
			return
		}
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		respondError(w, http.StatusBadRequest, "Invalid sort parameter (must be 'asc' or 'desc')")
		return
	}

	fmt.Printf("Query params: limit=%d, offset=%d, sort=%s\n", limit, offset, sortOrder)

	// Получаем заметки
	notes, err := h.storage.GetUserNotes(authenticatedUserID, limit, offset, sortOrder)
	if err != nil {
		fmt.Println("ERROR: GetUserNotes failed:", err)
		respondError(w, http.StatusInternalServerError, "Failed to get notes")
		return
	}

	respondJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== GetNote called ===")

	// Получаем user_id из JWT токена
	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDFromURL, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if authenticatedUserID != userIDFromURL {
		respondError(w, http.StatusForbidden, "You can only view your own notes")
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "note_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	fmt.Printf("UserID: %d, NoteID: %d\n", authenticatedUserID, noteID)

	note, err := h.storage.GetNoteByID(noteID)
	if err != nil {
		if err == models.ErrNoteNotFound {
			respondError(w, http.StatusNotFound, "Note not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get note")
		return
	}

	// Проверяем ownership
	if note.UserID != authenticatedUserID {
		respondError(w, http.StatusForbidden, "You don't have permission to access this note")
		return
	}

	respondJSON(w, http.StatusOK, note)
}

// UpdateNote обрабатывает PUT /users/{id}/notes/{note_id}
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== UpdateNote called ===")

	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDFromURL, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if authenticatedUserID != userIDFromURL {
		respondError(w, http.StatusForbidden, "You can only update your own notes")
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "note_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	fmt.Printf("USER ID: %d, Note ID: %d\n", authenticatedUserID, noteID)

	existingNote, err := h.storage.GetNoteByID(noteID)
	if err != nil {
		if err == models.ErrNoteNotFound {
			respondError(w, http.StatusNotFound, "Note not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get note")
		return
	}

	if existingNote.UserID != authenticatedUserID {
		respondError(w, http.StatusForbidden, "You don't have permission to update this note")
		return
	}

	var req models.UdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	fmt.Printf("Parsed request: %+v\n", req)

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	note, err := h.storage.UpdateNote(noteID, req.Title, req.Content)
	if err != nil {
		fmt.Println("ERROR: UpdateNote failed:", err)
		respondError(w, http.StatusInternalServerError, "Failed to update note")
		return
	}

	fmt.Printf("Note Updated: %+v\n", note)
	respondJSON(w, http.StatusOK, note)
}

// DeleteNote обрабатывает DELETE /users/{id}/notes/{note_id}
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== DeleteNote called ===")

	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDFromURL, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if authenticatedUserID != userIDFromURL {
		respondError(w, http.StatusForbidden, "You can only delete your own notes")
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "note_id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	fmt.Printf("User ID: %d, Note ID: %d\n", authenticatedUserID, noteID)

	existingNote, err := h.storage.GetNoteByID(noteID)
	if err != nil {
		if err == models.ErrNoteNotFound {
			respondError(w, http.StatusNotFound, "Note not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get note")
		return
	}

	if existingNote.UserID != authenticatedUserID {
		respondError(w, http.StatusForbidden, "You don't have permission to delete this note")
		return
	}

	if err := h.storage.DeleteNote(noteID); err != nil {
		fmt.Println("ERROR: DeleteNote failed:", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Note deleted successfully"})
}
