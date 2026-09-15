package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"opsdesk/internal/config"
	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type AttachmentHandler struct {
	cfg *config.Config
}

func NewAttachmentHandler(cfg *config.Config) *AttachmentHandler {
	_ = os.MkdirAll(cfg.UploadDir, 0755)
	return &AttachmentHandler{cfg: cfg}
}

const MaxUploadSize = 10 * 1024 * 1024 // 10MB

func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
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

	// Verify ticket exists
	var creatorID int64
	err = database.DB.QueryRow("SELECT creator_id FROM tickets WHERE id = ?", ticketID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	if claims.Role == models.RoleEmployee && creatorID != claims.UserID {
		utils.Error(w, http.StatusForbidden, "Cannot attach files to other users' tickets")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		utils.Error(w, http.StatusBadRequest, "File exceeds maximum size of 10MB")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "File field is required")
		return
	}
	defer file.Close()

	// Safe unique filename generation
	randomBytes := make([]byte, 16)
	_, _ = rand.Read(randomBytes)
	safeExt := strings.ToLower(filepath.Ext(header.Filename))
	uniqueName := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), hex.EncodeToString(randomBytes), safeExt)

	destPath := filepath.Join(h.cfg.UploadDir, uniqueName)
	dst, err := os.Create(destPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to write file")
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	now := time.Now()
	res, err := database.DB.Exec(`
		INSERT INTO attachments (ticket_id, filename, file_path, file_size, mime_type, uploader_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, ticketID, filepath.Base(header.Filename), uniqueName, size, mimeType, claims.UserID, now)

	if err != nil {
		_ = os.Remove(destPath)
		utils.Error(w, http.StatusInternalServerError, "Failed to record attachment")
		return
	}

	attachmentID, _ := res.LastInsertId()
	utils.Success(w, http.StatusCreated, map[string]interface{}{
		"id":        attachmentID,
		"filename":  header.Filename,
		"file_size": size,
	}, "Attachment uploaded successfully")
}

func (h *AttachmentHandler) List(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, ticket_id, filename, file_path, file_size, mime_type, uploader_id, created_at
		FROM attachments
		WHERE ticket_id = ?
		ORDER BY created_at DESC
	`, ticketID)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve attachments")
		return
	}
	defer rows.Close()

	attachments := make([]models.Attachment, 0)
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(
			&a.ID, &a.TicketID, &a.Filename, &a.FilePath, &a.FileSize, &a.MimeType, &a.UploaderID, &a.CreatedAt,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to read attachment row")
			return
		}
		attachments = append(attachments, a)
	}

	utils.Success(w, http.StatusOK, attachments)
}

func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	attachmentIDStr := chi.URLParam(r, "attachmentId")
	attachmentID, err := strconv.ParseInt(attachmentIDStr, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid attachment ID")
		return
	}

	var a models.Attachment
	err = database.DB.QueryRow(`
		SELECT id, ticket_id, filename, file_path, file_size, mime_type, uploader_id, created_at
		FROM attachments
		WHERE id = ?
	`, attachmentID).Scan(
		&a.ID, &a.TicketID, &a.Filename, &a.FilePath, &a.FileSize, &a.MimeType, &a.UploaderID, &a.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(w, http.StatusNotFound, "Attachment not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Clean safe path resolution
	safeBase := filepath.Clean(h.cfg.UploadDir)
	targetPath := filepath.Join(safeBase, filepath.Base(a.FilePath))

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		utils.Error(w, http.StatusNotFound, "File not found on storage")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", a.Filename))
	w.Header().Set("Content-Type", a.MimeType)
	http.ServeFile(w, r, targetPath)
}
