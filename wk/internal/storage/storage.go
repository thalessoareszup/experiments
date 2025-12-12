package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"example.com/wk/internal/workflow"
)

// DB wraps a sql.DB for workflow state.
type DB struct {
	SQL *sql.DB
}

// Open opens the state database at the given path.
func Open(dbPath string) (*DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &DB{SQL: db}, nil
}

// Close closes the underlying DB.
func (db *DB) Close() error { return db.SQL.Close() }

// migrate ensures required tables exist and adds any missing columns.
func migrate(sqldb *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS runs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			started_at TIMESTAMP NOT NULL,
			status TEXT NOT NULL,
			current_step_index INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS run_steps (
			run_id INTEGER NOT NULL,
			step_index INTEGER NOT NULL,
			step_id TEXT NOT NULL,
			step_name TEXT NOT NULL,
			step_description TEXT,
			PRIMARY KEY (run_id, step_index),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
	}

	for _, stmt := range stmts {
		if _, err := sqldb.Exec(stmt); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}

	// Add confirmation columns if they don't exist (for existing databases)
	if err := addColumnIfNotExists(sqldb, "run_steps", "requires_confirmation", "BOOLEAN NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := addColumnIfNotExists(sqldb, "run_steps", "confirmed_at", "TIMESTAMP"); err != nil {
		return err
	}

	// Create reports table for agent status updates
	if _, err := sqldb.Exec(`CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		run_id INTEGER NOT NULL,
		step_index INTEGER NOT NULL,
		report TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
	)`); err != nil {
		return fmt.Errorf("migrate reports: %w", err)
	}

	return nil
}

// addColumnIfNotExists adds a column to a table if it doesn't already exist.
func addColumnIfNotExists(sqldb *sql.DB, table, column, columnDef string) error {
	// Check if column exists
	var count int
	err := sqldb.QueryRow(
		`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?`,
		table, column,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("check column existence: %w", err)
	}

	if count == 0 {
		// Column doesn't exist, add it
		stmt := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columnDef)
		if _, err := sqldb.Exec(stmt); err != nil {
			return fmt.Errorf("add column %s.%s: %w", table, column, err)
		}
	}

	return nil
}

// Run represents a workflow execution run.
type Run struct {
	ID               int64
	StartedAt        time.Time
	Status           string
	CurrentStepIndex int
}

// StepSnapshot represents a step stored in the DB for a particular run.
type StepSnapshot struct {
	Index                int
	ID                   string
	Name                 string
	Description          string
	RequiresConfirmation bool
	ConfirmedAt          *time.Time
}

// StartRun creates a new run for the given workflow, starting at step 0 and
// storing all steps in the database.
func (db *DB) StartRun(ctx context.Context, wf *workflow.Workflow) (*Run, error) {
	if len(wf.Steps) == 0 {
		return nil, fmt.Errorf("cannot start run: workflow has no steps")
	}

	tx, err := db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now().UTC()
	res, err := tx.ExecContext(ctx,
		`INSERT INTO runs (started_at, status, current_step_index)
		 VALUES (?, ?, ?)`,
		now, "in_progress", 0,
	)
	if err != nil {
		return nil, fmt.Errorf("insert run: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	for i, s := range wf.Steps {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO run_steps (run_id, step_index, step_id, step_name, step_description, requires_confirmation)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			id, i, s.ID, s.Name, s.Description, s.RequiresConfirmation,
		); err != nil {
			return nil, fmt.Errorf("insert run_steps: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &Run{
		ID:               id,
		StartedAt:        now,
		Status:           "in_progress",
		CurrentStepIndex: 0,
	}, nil
}

// LatestRun returns the most recent run, or sql.ErrNoRows if none.
func (db *DB) LatestRun(ctx context.Context) (*Run, error) {
	row := db.SQL.QueryRowContext(ctx,
		`SELECT id, started_at, status, current_step_index
		 FROM runs
		 ORDER BY id DESC
		 LIMIT 1`)

	var r Run
	if err := row.Scan(&r.ID, &r.StartedAt, &r.Status, &r.CurrentStepIndex); err != nil {
		return nil, err
	}
	return &r, nil
}

// StepCount returns the number of steps stored for a given run.
func (db *DB) StepCount(ctx context.Context, runID int64) (int, error) {
	row := db.SQL.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM run_steps WHERE run_id = ?`, runID)
	var n int
	if err := row.Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

// CurrentStep returns the snapshot of the current step for the given run.
func (db *DB) CurrentStep(ctx context.Context, runID int64, index int) (*StepSnapshot, error) {
	row := db.SQL.QueryRowContext(ctx,
		`SELECT step_index, step_id, step_name, step_description, requires_confirmation, confirmed_at
		 FROM run_steps
		 WHERE run_id = ? AND step_index = ?`, runID, index)

	var s StepSnapshot
	if err := row.Scan(&s.Index, &s.ID, &s.Name, &s.Description, &s.RequiresConfirmation, &s.ConfirmedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

// ErrNoRuns is returned when there are no runs recorded.
var ErrNoRuns = errors.New("no runs found")

// LatestRunWithCurrentStep returns the latest run, its current step, and
// total number of steps, or ErrNoRuns if none exist.
func (db *DB) LatestRunWithCurrentStep(ctx context.Context) (*Run, *StepSnapshot, int, error) {
	run, err := db.LatestRun(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, 0, ErrNoRuns
		}
		return nil, nil, 0, err
	}

	total, err := db.StepCount(ctx, run.ID)
	if err != nil {
		return nil, nil, 0, err
	}

	step, err := db.CurrentStep(ctx, run.ID, run.CurrentStepIndex)
	if err != nil {
		return nil, nil, 0, err
	}

	return run, step, total, nil
}

// ErrAlreadyAtLastStep is returned when trying to advance beyond the last step.
var ErrAlreadyAtLastStep = errors.New("already at last step")

// AdvanceLatestRun moves the latest run to the next step, updating status as needed,
// and returns the updated run, new current step, and total steps.
func (db *DB) AdvanceLatestRun(ctx context.Context) (*Run, *StepSnapshot, int, error) {
	run, err := db.LatestRun(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, 0, ErrNoRuns
		}
		return nil, nil, 0, err
	}

	total, err := db.StepCount(ctx, run.ID)
	if err != nil {
		return nil, nil, 0, err
	}

	if total == 0 {
		return nil, nil, 0, fmt.Errorf("run has no steps recorded")
	}

	if run.CurrentStepIndex >= total-1 {
		// Mark as completed but do not advance beyond last step.
		if _, err := db.SQL.ExecContext(ctx,
			`UPDATE runs SET status = ? WHERE id = ?`, "completed", run.ID,
		); err != nil {
			return nil, nil, 0, err
		}
		run.Status = "completed"
		return run, nil, total, ErrAlreadyAtLastStep
	}

	newIndex := run.CurrentStepIndex + 1
	if _, err := db.SQL.ExecContext(ctx,
		`UPDATE runs SET current_step_index = ? WHERE id = ?`, newIndex, run.ID,
	); err != nil {
		return nil, nil, 0, err
	}

	run.CurrentStepIndex = newIndex

	step, err := db.CurrentStep(ctx, run.ID, newIndex)
	if err != nil {
		return nil, nil, 0, err
	}

	return run, step, total, nil
}

// ConfirmStep marks the specified step as confirmed with the current timestamp.
// Returns error if step does not exist, does not require confirmation, or is already confirmed.
func (db *DB) ConfirmStep(ctx context.Context, runID int64, stepIndex int) error {
	now := time.Now().UTC()

	result, err := db.SQL.ExecContext(ctx,
		`UPDATE run_steps
		 SET confirmed_at = ?
		 WHERE run_id = ? AND step_index = ? AND requires_confirmation = 1 AND confirmed_at IS NULL`,
		now, runID, stepIndex,
	)
	if err != nil {
		return fmt.Errorf("confirm step: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("step does not exist, does not require confirmation, or is already confirmed")
	}

	return nil
}

// IsAwaitingConfirmation returns true if the run's current step requires confirmation
// and has not yet been confirmed.
func (db *DB) IsAwaitingConfirmation(ctx context.Context, runID int64, stepIndex int) (bool, error) {
	var requiresConf bool
	var confirmedAt *time.Time

	err := db.SQL.QueryRowContext(ctx,
		`SELECT requires_confirmation, confirmed_at
		 FROM run_steps
		 WHERE run_id = ? AND step_index = ?`,
		runID, stepIndex,
	).Scan(&requiresConf, &confirmedAt)

	if err != nil {
		return false, err
	}

	return requiresConf && confirmedAt == nil, nil
}

// Report represents an agent status update for a step.
type Report struct {
	ID        int64
	RunID     int64
	StepIndex int
	Report    string
	CreatedAt time.Time
}

// AddReport inserts a new report for the current step of the latest run.
func (db *DB) AddReport(ctx context.Context, runID int64, stepIndex int, report string) error {
	now := time.Now().UTC()
	_, err := db.SQL.ExecContext(ctx,
		`INSERT INTO reports (run_id, step_index, report, created_at)
		 VALUES (?, ?, ?, ?)`,
		runID, stepIndex, report, now,
	)
	if err != nil {
		return fmt.Errorf("insert report: %w", err)
	}
	return nil
}

// ReportsForStep returns all reports for a given run and step index.
func (db *DB) ReportsForStep(ctx context.Context, runID int64, stepIndex int) ([]Report, error) {
	rows, err := db.SQL.QueryContext(ctx,
		`SELECT id, run_id, step_index, report, created_at
		 FROM reports
		 WHERE run_id = ? AND step_index = ?
		 ORDER BY created_at ASC`,
		runID, stepIndex,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var r Report
		if err := rows.Scan(&r.ID, &r.RunID, &r.StepIndex, &r.Report, &r.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, rows.Err()
}

// ReportsForRun returns all reports for a given run, grouped by step index.
func (db *DB) ReportsForRun(ctx context.Context, runID int64) (map[int][]Report, error) {
	rows, err := db.SQL.QueryContext(ctx,
		`SELECT id, run_id, step_index, report, created_at
		 FROM reports
		 WHERE run_id = ?
		 ORDER BY step_index ASC, created_at ASC`,
		runID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := make(map[int][]Report)
	for rows.Next() {
		var r Report
		if err := rows.Scan(&r.ID, &r.RunID, &r.StepIndex, &r.Report, &r.CreatedAt); err != nil {
			return nil, err
		}
		reports[r.StepIndex] = append(reports[r.StepIndex], r)
	}
	return reports, rows.Err()
}
