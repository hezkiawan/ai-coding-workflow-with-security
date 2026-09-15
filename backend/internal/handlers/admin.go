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

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

type UpdateUserRoleRequest struct {
	Role models.Role `json:"role"`
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, username, email, full_name, department, role, created_at, updated_at
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load users")
		return
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FullName, &u.Department, &u.Role, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse user")
			return
		}
		users = append(users, u)
	}

	utils.Success(w, http.StatusOK, users)
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Role != models.RoleAdmin && req.Role != models.RoleTechnician && req.Role != models.RoleEmployee {
		utils.Error(w, http.StatusBadRequest, "Invalid role specified")
		return
	}

	now := time.Now()
	_, err = database.DB.Exec("UPDATE users SET role = ?, updated_at = ? WHERE id = ?", req.Role, now, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	claims := middleware.GetCurrentUser(r)
	if claims != nil {
		_, _ = database.DB.Exec(`
			INSERT INTO audit_logs (user_id, action, entity, entity_id, details, ip_address, created_at)
			VALUES (?, 'USER_ROLE_UPDATE', 'User', ?, ?, ?, ?)
		`, claims.UserID, id, "Updated role to "+string(req.Role), r.RemoteAddr, now)
	}

	utils.Success(w, http.StatusOK, nil, "User role updated successfully")
}

func (h *AdminHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT a.id, a.user_id, a.action, a.entity, a.entity_id, a.details, a.ip_address, a.created_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to fetch audit logs")
		return
	}
	defer rows.Close()

	logs := make([]models.AuditLog, 0)
	for rows.Next() {
		var l models.AuditLog
		var u models.User
		var uID sql.NullInt64
		var uName, uFull, uEmail, uDept, uRole sql.NullString

		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.Entity, &l.EntityID, &l.Details, &l.IPAddress, &l.CreatedAt,
			&uID, &uName, &uFull, &uEmail, &uDept, &uRole,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse log row")
			return
		}

		if uID.Valid {
			u.ID = uID.Int64
			u.Username = uName.String
			u.FullName = uFull.String
			u.Email = uEmail.String
			u.Department = uDept.String
			u.Role = models.Role(uRole.String)
			l.User = &u
		}
		logs = append(logs, l)
	}

	utils.Success(w, http.StatusOK, logs)
}
