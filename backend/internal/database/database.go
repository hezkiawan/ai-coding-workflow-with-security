package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Optimize SQLite settings
	if _, err := db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA foreign_keys=ON;
		PRAGMA busy_timeout=5000;
	`); err != nil {
		log.Printf("Warning: failed to set PRAGMA: %v", err)
	}

	if err := migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	if err := seedInitialData(db); err != nil {
		log.Printf("Warning: failed to seed initial data: %v", err)
	}

	DB = db
	return db, nil
}

func migrateSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			full_name TEXT NOT NULL,
			department TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'employee',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'open',
			priority TEXT NOT NULL DEFAULT 'medium',
			department TEXT NOT NULL,
			creator_id INTEGER NOT NULL,
			assignee_id INTEGER,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ticket_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_tag TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			model TEXT NOT NULL,
			serial_number TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			location TEXT NOT NULL,
			assigned_to_id INTEGER,
			ip_address TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (assigned_to_id) REFERENCES users(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS attachments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ticket_id INTEGER NOT NULL,
			filename TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			mime_type TEXT NOT NULL,
			uploader_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
			FOREIGN KEY (uploader_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			action TEXT NOT NULL,
			entity TEXT NOT NULL,
			entity_id INTEGER NOT NULL,
			details TEXT NOT NULL,
			ip_address TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS integrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			target_url TEXT NOT NULL,
			secret_key TEXT NOT NULL,
			is_enabled INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS revoked_tokens (
			token_id TEXT PRIMARY KEY,
			revoked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("schema execution failed on %s: %w", query, err)
		}
	}
	return nil
}

func seedInitialData(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("AdminPass123!"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	techPassword, _ := bcrypt.GenerateFromPassword([]byte("TechPass123!"), bcrypt.DefaultCost)
	empPassword, _ := bcrypt.GenerateFromPassword([]byte("UserPass123!"), bcrypt.DefaultCost)

	now := time.Now()

	// Seed users
	_, err = db.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, department, role, created_at, updated_at)
		VALUES
		('admin', 'admin@opsdesk.internal', ?, 'System Administrator', 'IT Operations', 'admin', ?, ?),
		('sarah.tech', 'sarah.tech@opsdesk.internal', ?, 'Sarah Connor', 'Infrastructure', 'technician', ?, ?),
		('john.doe', 'john.doe@opsdesk.internal', ?, 'John Doe', 'Finance', 'employee', ?, ?);
	`, string(hashedPassword), now, now, string(techPassword), now, now, string(empPassword), now, now)
	if err != nil {
		return fmt.Errorf("failed to insert initial users: %w", err)
	}

	// Seed tickets
	_, err = db.Exec(`
		INSERT INTO tickets (title, description, status, priority, department, creator_id, assignee_id, notes, created_at, updated_at)
		VALUES
		('VPN Access Request for Remote Work', 'Need corporate VPN provisioning for the Q3 audit period.', 'in_progress', 'high', 'IT Operations', 3, 2, 'Verified employee manager approval.', ?, ?),
		('Office 365 License Activation Error', 'Encountering error code 0x80070005 when opening Excel.', 'open', 'medium', 'Finance', 3, NULL, '', ?, ?),
		('Server Room Core Switch Firmware Upgrade', 'Scheduled downtime window for firmware update on sw-core-01.', 'resolved', 'urgent', 'Infrastructure', 2, 2, 'Upgrade completed successfully at 02:00 AM.', ?, ?);
	`, now, now, now, now, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert initial tickets: %w", err)
	}

	// Seed comments
	_, err = db.Exec(`
		INSERT INTO comments (ticket_id, user_id, content, created_at)
		VALUES
		(1, 2, 'Reviewing your access request with IT security manager.', ?),
		(1, 3, 'Thank you Sarah, let me know if additional forms are required.', ?);
	`, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert initial comments: %w", err)
	}

	// Seed assets
	_, err = db.Exec(`
		INSERT INTO assets (asset_tag, name, category, model, serial_number, status, location, assigned_to_id, ip_address, created_at, updated_at)
		VALUES
		('AST-1001', 'MacBook Pro 16 - Finance', 'Laptop', 'Apple M3 Pro', 'C02G1234MD6R', 'active', 'Building A - Fl 3', 3, '10.10.4.55', ?, ?),
		('AST-1002', 'Dell PowerEdge R750', 'Server', 'Dell R750xs', 'SRV-8942-TX', 'active', 'Datacenter Rack 4', 2, '10.10.1.10', ?, ?),
		('AST-1003', 'Cisco Catalyst 9300', 'Switch', 'C9300-48P', 'FOC2349001', 'active', 'MDF Server Room', 2, '10.10.1.2', ?, ?),
		('AST-1004', 'Lenovo ThinkPad X1', 'Laptop', 'ThinkPad X1 Carbon Gen 11', 'PF-98213-US', 'maintenance', 'IT Repair Depot', NULL, '10.10.4.89', ?, ?);
	`, now, now, now, now, now, now, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert initial assets: %w", err)
	}

	// Seed audit log
	_, err = db.Exec(`
		INSERT INTO audit_logs (user_id, action, entity, entity_id, details, ip_address, created_at)
		VALUES
		(1, 'SYSTEM_INIT', 'System', 1, 'Initial database schema and default policies established', '127.0.0.1', ?);
	`, now)

	return err
}
