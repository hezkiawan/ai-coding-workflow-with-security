package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

func (h *CommentHandler) ListByTicket(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.ticket_id, c.user_id, c.content, c.created_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.ticket_id = ?
		ORDER BY c.created_at ASC
	`, ticketID)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load comments")
		return
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var c models.Comment
		var u models.User

		if err := rows.Scan(
			&c.ID, &c.TicketID, &c.UserID, &c.Content, &c.CreatedAt,
			&u.ID, &u.Username, &u.FullName, &u.Email, &u.Department, &u.Role,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse comments")
			return
		}
		c.Author = &u
		comments = append(comments, c)
	}

	utils.Success(w, http.StatusOK, comments)
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ticketIDStr := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Content == "" {
		utils.Error(w, http.StatusBadRequest, "Comment content cannot be empty")
		return
	}

	// Verify ticket exists
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM tickets WHERE id = ?", ticketID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	now := time.Now()
	res, err := database.DB.Exec(`
		INSERT INTO comments (ticket_id, user_id, content, created_at)
		VALUES (?, ?, ?, ?)
	`, ticketID, claims.UserID, req.Content, now)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	commentID, _ := res.LastInsertId()
	utils.Success(w, http.StatusCreated, map[string]interface{}{
		"id": commentID,
	}, "Comment added successfully")
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "commentId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	var authorID int64
	err = database.DB.QueryRow("SELECT user_id FROM comments WHERE id = ?", id).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Comment not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if claims.Role != models.RoleAdmin && authorID != claims.UserID {
		utils.Error(w, http.StatusForbidden, "You can only delete your own comments")
		return
	}

	_, err = database.DB.Exec("DELETE FROM comments WHERE id = ?", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	utils.Success(w, http.StatusOK, nil, "Comment deleted successfully")
}
