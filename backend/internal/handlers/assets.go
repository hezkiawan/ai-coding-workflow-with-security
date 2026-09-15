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

type AssetHandler struct{}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{}
}

type CreateAssetRequest struct {
	AssetTag     string             `json:"asset_tag"`
	Name         string             `json:"name"`
	Category     string             `json:"category"`
	Model        string             `json:"model"`
	SerialNumber string             `json:"serial_number"`
	Status       models.AssetStatus `json:"status"`
	Location     string             `json:"location"`
	AssignedToID *int64             `json:"assigned_to_id"`
	IPAddress    string             `json:"ip_address"`
}

type UpdateAssetRequest struct {
	Name         string             `json:"name"`
	Category     string             `json:"category"`
	Model        string             `json:"model"`
	SerialNumber string             `json:"serial_number"`
	Status       models.AssetStatus `json:"status"`
	Location     string             `json:"location"`
	AssignedToID *int64             `json:"assigned_to_id"`
	IPAddress    string             `json:"ip_address"`
}

func (h *AssetHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	category := r.URL.Query().Get("category")

	query := `
		SELECT a.id, a.asset_tag, a.name, a.category, a.model, a.serial_number,
		       a.status, a.location, a.assigned_to_id, a.ip_address, a.created_at, a.updated_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role
		FROM assets a
		LEFT JOIN users u ON a.assigned_to_id = u.id
		WHERE 1=1
	`
	var args []interface{}

	if status != "" {
		query += " AND a.status = ?"
		args = append(args, status)
	}
	if category != "" {
		query += " AND a.category = ?"
		args = append(args, category)
	}

	query += " ORDER BY a.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load assets")
		return
	}
	defer rows.Close()

	assets := make([]models.Asset, 0)
	for rows.Next() {
		var a models.Asset
		var assignedID sql.NullInt64
		var u models.User
		var uID, uUsername, uFullName, uEmail, uDept, uRole sql.NullString

		err := rows.Scan(
			&a.ID, &a.AssetTag, &a.Name, &a.Category, &a.Model, &a.SerialNumber,
			&a.Status, &a.Location, &assignedID, &a.IPAddress, &a.CreatedAt, &a.UpdatedAt,
			&uID, &uUsername, &uFullName, &uEmail, &uDept, &uRole,
		)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse asset data")
			return
		}

		if assignedID.Valid {
			id := assignedID.Int64
			a.AssignedToID = &id
			u.ID = id
			u.Username = uUsername.String
			u.FullName = uFullName.String
			u.Email = uEmail.String
			u.Department = uDept.String
			u.Role = models.Role(uRole.String)
			a.AssignedUser = &u
		}

		assets = append(assets, a)
	}

	utils.Success(w, http.StatusOK, assets)
}

func (h *AssetHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var a models.Asset
	var assignedID sql.NullInt64
	var u models.User
	var uID, uUsername, uFullName, uEmail, uDept, uRole sql.NullString

	err = database.DB.QueryRow(`
		SELECT a.id, a.asset_tag, a.name, a.category, a.model, a.serial_number,
		       a.status, a.location, a.assigned_to_id, a.ip_address, a.created_at, a.updated_at,
		       u.id, u.username, u.full_name, u.email, u.department, u.role
		FROM assets a
		LEFT JOIN users u ON a.assigned_to_id = u.id
		WHERE a.id = ?
	`, id).Scan(
		&a.ID, &a.AssetTag, &a.Name, &a.Category, &a.Model, &a.SerialNumber,
		&a.Status, &a.Location, &assignedID, &a.IPAddress, &a.CreatedAt, &a.UpdatedAt,
		&uID, &uUsername, &uFullName, &uEmail, &uDept, &uRole,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Asset not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to load asset")
		return
	}

	if assignedID.Valid {
		uid := assignedID.Int64
		a.AssignedToID = &uid
		u.ID = uid
		u.Username = uUsername.String
		u.FullName = uFullName.String
		u.Email = uEmail.String
		u.Department = uDept.String
		u.Role = models.Role(uRole.String)
		a.AssignedUser = &u
	}

	utils.Success(w, http.StatusOK, a)
}

func (h *AssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil || (claims.Role != models.RoleAdmin && claims.Role != models.RoleTechnician) {
		utils.Error(w, http.StatusForbidden, "Admin or Technician role required to create assets")
		return
	}

	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.AssetTag == "" || req.Name == "" || req.Category == "" {
		utils.Error(w, http.StatusBadRequest, "Asset tag, name and category are required")
		return
	}

	if req.Status == "" {
		req.Status = models.AssetStatusActive
	}

	now := time.Now()
	res, err := database.DB.Exec(`
		INSERT INTO assets (asset_tag, name, category, model, serial_number, status, location, assigned_to_id, ip_address, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.AssetTag, req.Name, req.Category, req.Model, req.SerialNumber, req.Status, req.Location, req.AssignedToID, req.IPAddress, now, now)

	if err != nil {
		utils.Error(w, http.StatusConflict, "Asset tag already exists or database error")
		return
	}

	assetID, _ := res.LastInsertId()
	utils.Success(w, http.StatusCreated, map[string]interface{}{
		"id": assetID,
	}, "Asset registered successfully")
}

func (h *AssetHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil || (claims.Role != models.RoleAdmin && claims.Role != models.RoleTechnician) {
		utils.Error(w, http.StatusForbidden, "Admin or Technician role required to modify assets")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	now := time.Now()
	_, err = database.DB.Exec(`
		UPDATE assets
		SET name = ?, category = ?, model = ?, serial_number = ?, status = ?, location = ?, assigned_to_id = ?, ip_address = ?, updated_at = ?
		WHERE id = ?
	`, req.Name, req.Category, req.Model, req.SerialNumber, req.Status, req.Location, req.AssignedToID, req.IPAddress, now, id)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	utils.Success(w, http.StatusOK, nil, "Asset updated successfully")
}

func (h *AssetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM assets WHERE id = ?", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to delete asset")
		return
	}

	utils.Success(w, http.StatusOK, nil, "Asset deleted successfully")
}

func (h *AssetHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status")

	sqlQuery := `
		SELECT a.id, a.asset_tag, a.name, a.category, a.model, a.serial_number,
		       a.status, a.location, a.assigned_to_id, a.ip_address, a.created_at, a.updated_at
		FROM assets a
		WHERE 1=1
	`
	var args []interface{}

	if query != "" {
		sqlQuery += " AND (a.name LIKE ? OR a.asset_tag LIKE ? OR a.model LIKE ? OR a.serial_number LIKE ?)"
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	if status != "" {
		sqlQuery += " AND a.status = ?"
		args = append(args, status)
	}

	sqlQuery += " ORDER BY a.name ASC"

	rows, err := database.DB.Query(sqlQuery, args...)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Asset search failed")
		return
	}
	defer rows.Close()

	assets := make([]models.Asset, 0)
	for rows.Next() {
		var a models.Asset
		var assignedID sql.NullInt64

		if err := rows.Scan(
			&a.ID, &a.AssetTag, &a.Name, &a.Category, &a.Model, &a.SerialNumber,
			&a.Status, &a.Location, &assignedID, &a.IPAddress, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse search results")
			return
		}
		if assignedID.Valid {
			aid := assignedID.Int64
			a.AssignedToID = &aid
		}
		assets = append(assets, a)
	}

	utils.Success(w, http.StatusOK, assets)
}
