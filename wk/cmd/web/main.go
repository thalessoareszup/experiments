package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example.com/wk/internal/storage"
)

// pageData holds data passed to the HTML template.
type pageData struct {
	Runs          []storage.Run
	StepsByRun    map[int64][]storage.StepSnapshot
	MessagesByRun map[int64]map[int][]storage.Message
	Error         string
	WorkflowFile  string
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	workflowFile := homeDir + "/.local/wk/workflow.yaml"
	dbFile := homeDir + "/.local/wk/wk.db"

	db, err := storage.Open(dbFile)
	if err != nil {
		log.Fatalf("open state db: %v", err)
	}
	defer db.Close()

	tmpl := template.Must(template.New("monitor.html").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format(time.RFC3339)
		},
	}).ParseFiles(filepath.Join("cmd", "web", "templates", "monitor.html")))

	mux := http.NewServeMux()

	mux.HandleFunc("/dbinfo", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		runs, stepsByRun, messagesByRun, err := loadDBInfo(ctx, db)
		data := pageData{Runs: runs, StepsByRun: stepsByRun, MessagesByRun: messagesByRun, WorkflowFile: workflowFile}
		if err != nil {
			log.Printf("loadDBInfo error: %v", err)
			data.Error = err.Error()
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template execute: %v", err)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dbinfo", http.StatusSeeOther)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Starting web server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("http server: %v", err)
	}
}

// loadDBInfo retrieves all runs, steps, and messages from the state DB.
func loadDBInfo(ctx context.Context, db *storage.DB) ([]storage.Run, map[int64][]storage.StepSnapshot, map[int64]map[int][]storage.Message, error) {
	rows, err := db.SQL.QueryContext(ctx,
		`SELECT id, started_at, status, current_step_index
		 FROM runs
		 ORDER BY id DESC`)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	var runs []storage.Run
	for rows.Next() {
		var r storage.Run
		if err := rows.Scan(&r.ID, &r.StartedAt, &r.Status, &r.CurrentStepIndex); err != nil {
			return nil, nil, nil, err
		}
		runs = append(runs, r)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	stepsByRun := make(map[int64][]storage.StepSnapshot, len(runs))
	messagesByRun := make(map[int64]map[int][]storage.Message, len(runs))

	for _, r := range runs {
		srows, err := db.SQL.QueryContext(ctx,
			`SELECT step_index, step_id, step_name, step_description, requires_confirmation, confirmed_at
			 FROM run_steps
			 WHERE run_id = ?
			 ORDER BY step_index ASC`, r.ID)
		if err != nil {
			return nil, nil, nil, err
		}

		var steps []storage.StepSnapshot
		for srows.Next() {
			var s storage.StepSnapshot
			if err := srows.Scan(&s.Index, &s.ID, &s.Name, &s.Description, &s.RequiresConfirmation, &s.ConfirmedAt); err != nil {
				srows.Close()
				return nil, nil, nil, err
			}
			steps = append(steps, s)
		}
		srows.Close()
		if err := srows.Err(); err != nil {
			return nil, nil, nil, err
		}

		stepsByRun[r.ID] = steps

		messages, err := db.MessagesForRun(ctx, r.ID)
		if err != nil {
			return nil, nil, nil, err
		}
		messagesByRun[r.ID] = messages
	}

	return runs, stepsByRun, messagesByRun, nil
}
