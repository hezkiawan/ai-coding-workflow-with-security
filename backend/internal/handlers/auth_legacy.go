package handlers

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"opsdesk/internal/database"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type LegacyAuthHandler struct{}

func NewLegacyAuthHandler() *LegacyAuthHandler {
	return &LegacyAuthHandler{}
}

type LegacySSORequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Provider string `json:"provider"`
}

type LegacyServiceRequest struct {
	ServiceName string `json:"service_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

func hashLegacyPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *LegacyAuthHandler) LegacyLogin(w http.ResponseWriter, r *http.Request) {
	var req LegacySSORequest
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
		WHERE username = ?
	`, req.Username).Scan(
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

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("sso"))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create legacy session")
		return
	}

	utils.Success(w, http.StatusOK, LoginResponse{
		Token: tokenStr,
		User:  &user,
	}, "Legacy SSO login successful")
}

func (h *LegacyAuthHandler) LegacyLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_cleared",
		Value:    "true",
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		MaxAge:   86400,
	})

	utils.Success(w, http.StatusOK, nil, "Legacy session terminated")
}

func (h *LegacyAuthHandler) RegisterServiceAccount(w http.ResponseWriter, r *http.Request) {
	var req LegacyServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.ServiceName == "" || req.Username == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "Service name, username and password are required")
		return
	}

	passwordDigest := hashLegacyPassword(req.Password)

	now := time.Now()
	email := req.Email
	if email == "" {
		email = req.Username + "@service.opsdesk.internal"
	}

	res, err := database.DB.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, department, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'Service Accounts', 'employee', ?, ?)
	`, req.Username, email, passwordDigest, req.ServiceName, now, now)

	if err != nil {
		utils.Error(w, http.StatusConflict, "Service account already exists")
		return
	}

	accountID, _ := res.LastInsertId()
	utils.Success(w, http.StatusCreated, map[string]interface{}{
		"id":           accountID,
		"service_name": req.ServiceName,
	}, "Service account registered")
}
