package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const dbFile = "tasks.db"

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// Create tasks table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTask(db *sql.DB, description string) error {
	_, err := db.Exec("INSERT INTO tasks (description) VALUES (?)", description)
	if err != nil {
		return err
	}
	fmt.Printf("Task created: %s\n", description)
	return nil
}

func listTasks(db *sql.DB) error {
	rows, err := db.Query("SELECT id, description, created_at FROM tasks ORDER BY id")
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("\nTasks:")
	fmt.Println("------")
	
	hasRows := false
	for rows.Next() {
		hasRows = true
		var task Task
		var createdAt string
		
		err := rows.Scan(&task.ID, &task.Description, &createdAt)
		if err != nil {
			return err
		}
		
		// Parse the timestamp
		task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		
		fmt.Printf("[%d] %s (created: %s)\n", task.ID, task.Description, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	
	if !hasRows {
		fmt.Println("No tasks found.")
	}
	
	return rows.Err()
}

func printUsage() {
	fmt.Println("Task Manager CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  taskmanager create <description>  - Create a new task")
	fmt.Println("  taskmanager list                  - List all tasks")
	fmt.Println("\nExamples:")
	fmt.Println("  taskmanager create \"Buy groceries\"")
	fmt.Println("  taskmanager list")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Error: task description is required")
			fmt.Println("\nUsage: taskmanager create <description>")
			os.Exit(1)
		}
		description := os.Args[2]
		if err := createTask(db, description); err != nil {
			log.Fatalf("Failed to create task: %v", err)
		}

	case "list":
		if err := listTasks(db); err != nil {
			log.Fatalf("Failed to list tasks: %v", err)
		}

	default:
		fmt.Printf("Error: unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
}
