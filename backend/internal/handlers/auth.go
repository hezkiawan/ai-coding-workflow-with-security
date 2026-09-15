package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"opsdesk/internal/config"
	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type RegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Department string `json:"department"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		utils.Error(w, http.StatusBadRequest, "All required fields must be provided")
		return
	}

	if len(req.Password) < 8 {
		utils.Error(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to process credentials")
		return
	}

	if req.Department == "" {
		req.Department = "General"
	}

	now := time.Now()
	res, err := database.DB.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, department, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'employee', ?, ?)
	`, req.Username, req.Email, hashedPassword, req.FullName, req.Department, now, now)

	if err != nil {
		utils.Error(w, http.StatusConflict, "Username or email already exists")
		return
	}

	userID, _ := res.LastInsertId()
	user := &models.User{
		ID:         userID,
		Username:   req.Username,
		Email:      req.Email,
		FullName:   req.FullName,
		Department: req.Department,
		Role:       models.RoleEmployee,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	token, err := utils.GenerateToken(user, h.cfg.JWTSecret, 24*time.Hour)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	utils.Success(w, http.StatusCreated, LoginResponse{
		Token: token,
		User:  user,
	}, "User registered successfully")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	var user models.User
	var passwordHash string
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, full_name, department, role, created_at, updated_at
		FROM users
		WHERE username = ? OR email = ?
	`, req.Username, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.FullName,
		&user.Department,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Authentication failed")
		return
	}

	if !utils.CheckPasswordHash(req.Password, passwordHash) {
		utils.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(&user, h.cfg.JWTSecret, 24*time.Hour)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	utils.Success(w, http.StatusOK, LoginResponse{
		Token: token,
		User:  &user,
	}, "Login successful")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, http.StatusOK, nil, "Logged out successfully")
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, username, email, full_name, department, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`, claims.UserID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Department,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		utils.Error(w, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.UserProfileUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	now := time.Now()
	_, err := database.DB.Exec(`
		UPDATE users
		SET full_name = ?, department = ?, updated_at = ?
		WHERE id = ?
	`, req.FullName, req.Department, now, claims.UserID)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.Success(w, http.StatusOK, nil, "Profile updated successfully")
}
