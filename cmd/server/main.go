package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go-app/internal/task"
	"go-app/internal/web"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	var err error
	// Open SQLite database file
	db, err = sql.Open("sqlite", "./app.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0,
		task_group TEXT DEFAULT ''
	);`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatalf("Failed to initialize DB schema: %v", err)
	}
	apiToken := os.Getenv("TODO_API_TOKEN")
	if apiToken == "" {
		log.Fatal("TODO_API_TOKEN must be set")
	}

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /tasks", getTasks)
	protectedMux.HandleFunc("POST /tasks", createTask)
	protectedMux.HandleFunc("PATCH /tasks/{id}/toggle", toggleTask)
	protectedMux.HandleFunc("DELETE /tasks/{id}", deleteTask)
	protectedMux.HandleFunc("GET /groups", getGroups)
	protectedMux.HandleFunc("GET /web/todos", getWebTodos)
	protectedMux.HandleFunc("GET /web/groups", getWebGroups)
	protectedMux.HandleFunc("POST /web/todos", createWebTodo)
	protectedMux.HandleFunc("PATCH /web/todos/{id}/toggle", toggleWebTodo)
	protectedMux.HandleFunc("DELETE /web/todos/{id}", deleteTask)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", web.RenderPage)
	mux.Handle("/", requireBearerToken(apiToken, protectedMux))

	fmt.Println("Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

func requireBearerToken(expected string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, prefix) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		provided := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
		if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := loadTasks(r.URL.Query().Get("group"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func loadTasks(groupFilter string) ([]task.Task, error) {
	var rows *sql.Rows
	var err error
	if groupFilter == "" {
		rows, err = db.Query("SELECT id, title, done, COALESCE(task_group, '') FROM tasks ORDER BY id ASC")
	} else {
		rows, err = db.Query("SELECT id, title, done, COALESCE(task_group, '') FROM tasks WHERE task_group = ? ORDER BY id ASC", groupFilter)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []task.Task{}
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.Group); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func getWebTodos(w http.ResponseWriter, r *http.Request) {
	tasks, err := loadTasks(r.URL.Query().Get("group"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]interface{}, len(tasks))
	for i := range tasks {
		items[i] = tasks[i]
	}
	web.RenderList(w, items)
}

func getWebGroups(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT DISTINCT task_group FROM tasks WHERE COALESCE(task_group, '') <> '' ORDER BY task_group")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	groups := []string{}
	for rows.Next() {
		var group string
		if err := rows.Scan(&group); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createWebTodo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	group := strings.TrimSpace(r.FormValue("group"))

	result, err := db.Exec("INSERT INTO tasks (title, task_group) VALUES (?, ?)", title, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.RenderItem(w, task.Task{ID: int(id), Title: title, Group: group})
}

func toggleWebTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET done = NOT done WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var t task.Task
	err = db.QueryRow("SELECT id, title, done, COALESCE(task_group, '') FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Title, &t.Done, &t.Group)
	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.RenderItem(w, t)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO tasks (title, task_group) VALUES (?, ?)", t.Title, t.Group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	t.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func toggleTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET done = NOT done WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var t task.Task
	err = db.QueryRow("SELECT id, title, done, COALESCE(task_group, '') FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Title, &t.Done, &t.Group)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Dynamic ID remapping in SQLite
	_, _ = db.Exec(`
		CREATE TABLE tasks_temp AS SELECT title, done, task_group FROM tasks ORDER BY id ASC;
		DROP TABLE tasks;
		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0,
			task_group TEXT DEFAULT ''
		);
		INSERT INTO tasks (title, done, task_group) SELECT title, done, task_group FROM tasks_temp;
		DROP TABLE tasks_temp;
	`)

	w.WriteHeader(http.StatusNoContent)
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT COALESCE(task_group, 'default'), COUNT(*) FROM tasks GROUP BY task_group")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	groups := make(map[string]int)
	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err == nil {
			if name == "" {
				name = "default"
			}
			groups[name] = count
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
