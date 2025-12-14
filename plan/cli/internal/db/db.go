package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	dbPath := getDBPath()

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create db directory: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func getDBPath() string {
	if path := os.Getenv("PLAN_DB_PATH"); path != "" {
		return path
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "plan", "plan.db")
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS plans (
		id TEXT PRIMARY KEY,
		parent_id TEXT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'failed')),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (parent_id) REFERENCES plans(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS steps (
		id TEXT PRIMARY KEY,
		plan_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'failed')),
		step_order INTEGER NOT NULL,
		progress INTEGER DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_plans_parent_id ON plans(parent_id);
	CREATE INDEX IF NOT EXISTS idx_plans_status ON plans(status);
	CREATE INDEX IF NOT EXISTS idx_steps_plan_id ON steps(plan_id);
	CREATE INDEX IF NOT EXISTS idx_steps_status ON steps(status);
	`

	_, err := DB.Exec(schema)
	return err
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
