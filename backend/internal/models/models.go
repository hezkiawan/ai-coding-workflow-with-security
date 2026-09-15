package models

import (
	"time"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleTechnician Role = "technician"
	RoleEmployee   Role = "employee"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Department   string    `json:"department"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserProfileUpdate struct {
	FullName   string `json:"full_name"`
	Department string `json:"department"`
}

type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusResolved   TicketStatus = "resolved"
	StatusClosed     TicketStatus = "closed"
)

type TicketPriority string

const (
	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
	PriorityUrgent TicketPriority = "urgent"
)

type Ticket struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      TicketStatus   `json:"status"`
	Priority    TicketPriority `json:"priority"`
	Department  string         `json:"department"`
	CreatorID   int64          `json:"creator_id"`
	AssigneeID  *int64         `json:"assignee_id,omitempty"`
	Notes       string         `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Populated relations
	Creator  *User `json:"creator,omitempty"`
	Assignee *User `json:"assignee,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Populated relations
	Author *User `json:"author,omitempty"`
}

type AssetStatus string

const (
	AssetStatusActive      AssetStatus = "active"
	AssetStatusMaintenance AssetStatus = "maintenance"
	AssetStatusDecommissioned AssetStatus = "decommissioned"
)

type Asset struct {
	ID           int64       `json:"id"`
	AssetTag     string      `json:"asset_tag"`
	Name         string      `json:"name"`
	Category     string      `json:"category"` // Laptop, Server, Switch, Monitor, etc.
	Model        string      `json:"model"`
	SerialNumber string      `json:"serial_number"`
	Status       AssetStatus `json:"status"`
	Location     string      `json:"location"`
	AssignedToID *int64      `json:"assigned_to_id,omitempty"`
	IPAddress    string      `json:"ip_address,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// Populated relation
	AssignedUser *User `json:"assigned_user,omitempty"`
}

type Attachment struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	Filename  string    `json:"filename"`
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	UploaderID int64    `json:"uploader_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id,omitempty"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  int64     `json:"entity_id"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`

	User *User `json:"user,omitempty"`
}

type IntegrationWebhook struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	TargetURL string    `json:"target_url"`
	SecretKey string    `json:"secret_key"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
}
