package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type TicketHandler struct{}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{}
}

type CreateTicketRequest struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Priority    models.TicketPriority `json:"priority"`
	Department  string                `json:"department"`
	CreatorID   *int64                `json:"creator_id,omitempty"`
	Status      models.TicketStatus   `json:"status,omitempty"`
	AssigneeID  *int64                `json:"assignee_id,omitempty"`
}

type UpdateTicketRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      models.TicketStatus    `json:"status"`
	Priority    models.TicketPriority  `json:"priority"`
	Department  string                 `json:"department"`
	AssigneeID  *int64                 `json:"assignee_id"`
	Notes       string                 `json:"notes"`
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	department := r.URL.Query().Get("department")
	priority := r.URL.Query().Get("priority")

	query := `
		SELECT t.id, t.title, t.description, t.status, t.priority, t.department,
		       t.creator_id, t.assignee_id, t.notes, t.created_at, t.updated_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role,
		       a.id, a.username, a.full_name, a.email, a.department, a.role
		FROM tickets t
		JOIN users u ON t.creator_id = u.id
		LEFT JOIN users a ON t.assignee_id = a.id
		WHERE 1=1
	`
	var args []interface{}

	if status != "" {
		query += " AND t.status = ?"
		args = append(args, status)
	}
	if department != "" {
		query += " AND t.department = ?"
		args = append(args, department)
	}
	if priority != "" {
		query += " AND t.priority = ?"
		args = append(args, priority)
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve tickets")
		return
	}
	defer rows.Close()

	tickets := make([]models.Ticket, 0)
	for rows.Next() {
		var t models.Ticket
		var creator models.User
		var assigneeID sql.NullInt64
		var assignee models.User
		var aID, aUsername, aFullName, aEmail, aDept, aRole sql.NullString

		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Department,
			&t.CreatorID, &assigneeID, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
			&creator.ID, &creator.Username, &creator.FullName, &creator.Email, &creator.Department, &creator.Role,
			&aID, &aUsername, &aFullName, &aEmail, &aDept, &aRole,
		)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse ticket data")
			return
		}

		t.Creator = &creator
		if assigneeID.Valid {
			id := assigneeID.Int64
			t.AssigneeID = &id
			assignee.ID = id
			assignee.Username = aUsername.String
			assignee.FullName = aFullName.String
			assignee.Email = aEmail.String
			assignee.Department = aDept.String
			assignee.Role = models.Role(aRole.String)
			t.Assignee = &assignee
		}

		tickets = append(tickets, t)
	}

	utils.Success(w, http.StatusOK, tickets)
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	var t models.Ticket
	var creator models.User
	var assigneeID sql.NullInt64
	var assignee models.User
	var aID, aUsername, aFullName, aEmail, aDept, aRole sql.NullString

	err = database.DB.QueryRow(`
		SELECT t.id, t.title, t.description, t.status, t.priority, t.department,
		       t.creator_id, t.assignee_id, t.notes, t.created_at, t.updated_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role,
		       a.id, a.username, a.full_name, a.email, a.department, a.role
		FROM tickets t
		JOIN users u ON t.creator_id = u.id
		LEFT JOIN users a ON t.assignee_id = a.id
		WHERE t.id = ?
	`, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Department,
		&t.CreatorID, &assigneeID, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
		&creator.ID, &creator.Username, &creator.FullName, &creator.Email, &creator.Department, &creator.Role,
		&aID, &aUsername, &aFullName, &aEmail, &aDept, &aRole,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to load ticket")
		return
	}

	t.Creator = &creator
	if assigneeID.Valid {
		aid := assigneeID.Int64
		t.AssigneeID = &aid
		assignee.ID = aid
		assignee.Username = aUsername.String
		assignee.FullName = aFullName.String
		assignee.Email = aEmail.String
		assignee.Department = aDept.String
		assignee.Role = models.Role(aRole.String)
		t.Assignee = &assignee
	}

	utils.Success(w, http.StatusOK, t)
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Title == "" || req.Description == "" {
		utils.Error(w, http.StatusBadRequest, "Title and description are required")
		return
	}

	if req.Priority == "" {
		req.Priority = models.PriorityMedium
	}
	if req.Department == "" {
		if claims.Role == models.RoleAdmin {
			req.Department = "IT Operations"
		} else {
			req.Department = "General"
		}
	}

	creatorID := claims.UserID
	if req.CreatorID != nil && *req.CreatorID > 0 {
		creatorID = *req.CreatorID
	}

	status := models.StatusOpen
	if req.Status != "" {
		status = req.Status
	}

	now := time.Now()
	res, err := database.DB.Exec(`
		INSERT INTO tickets (title, description, status, priority, department, creator_id, assignee_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Title, req.Description, status, req.Priority, req.Department, creatorID, req.AssigneeID, now, now)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create ticket")
		return
	}

	ticketID, _ := res.LastInsertId()
	utils.Success(w, http.StatusCreated, map[string]interface{}{
		"id": ticketID,
	}, "Ticket created successfully")
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	var req UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Verify ticket exists
	var creatorID int64
	err = database.DB.QueryRow("SELECT creator_id FROM tickets WHERE id = ?", id).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to verify ticket")
		return
	}

	// Access control: only creator, technician, or admin can update
	if claims.Role == models.RoleEmployee && creatorID != claims.UserID {
		utils.Error(w, http.StatusForbidden, "You can only update your own tickets")
		return
	}

	now := time.Now()
	_, err = database.DB.Exec(`
		UPDATE tickets
		SET title = ?, description = ?, status = ?, priority = ?, department = ?, assignee_id = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`, req.Title, req.Description, req.Status, req.Priority, req.Department, req.AssigneeID, req.Notes, now, id)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update ticket")
		return
	}

	utils.Success(w, http.StatusOK, nil, "Ticket updated successfully")
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM tickets WHERE id = ?", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to delete ticket")
		return
	}

	utils.Success(w, http.StatusOK, nil, "Ticket deleted successfully")
}

func (h *TicketHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.List(w, r)
		return
	}

	sqlQuery := fmt.Sprintf(`
		SELECT t.id, t.title, t.description, t.status, t.priority, t.department,
		       t.creator_id, t.assignee_id, t.notes, t.created_at, t.updated_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role
		FROM tickets t
		JOIN users u ON t.creator_id = u.id
		WHERE t.title LIKE '%%%s%%' OR t.description LIKE '%%%s%%' OR t.department LIKE '%%%s%%'
		ORDER BY t.created_at DESC
	`, query, query, query)

	rows, err := database.DB.Query(sqlQuery)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Search query failed")
		return
	}
	defer rows.Close()

	tickets := make([]models.Ticket, 0)
	for rows.Next() {
		var t models.Ticket
		var creator models.User
		var assigneeID sql.NullInt64

		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Department,
			&t.CreatorID, &assigneeID, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
			&creator.ID, &creator.Username, &creator.FullName, &creator.Email, &creator.Department, &creator.Role,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse search results")
			return
		}
		t.Creator = &creator
		if assigneeID.Valid {
			aid := assigneeID.Int64
			t.AssigneeID = &aid
		}
		tickets = append(tickets, t)
	}

	utils.Success(w, http.StatusOK, tickets)
}
