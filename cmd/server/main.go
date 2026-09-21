package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"go-app/internal/task"
)

var (
	tasks  = []task.Task{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("POST /tasks", createTask)
	mux.HandleFunc("PATCH /tasks/{id}/toggle", toggleTask)
	mux.HandleFunc("DELETE /tasks/{id}", deleteTask)
	mux.HandleFunc("GET /groups", getGroups)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groupFilter := r.URL.Query().Get("group")

	mu.Lock()
	defer mu.Unlock()

	if groupFilter == "" {
		json.NewEncoder(w).Encode(tasks)
		return
	}

	filtered := []task.Task{}
	for _, t := range tasks {
		if t.Group == groupFilter {
			filtered = append(filtered, t)
		}
	}
	json.NewEncoder(w).Encode(filtered)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	t.ID = nextID
	nextID++
	tasks = append(tasks, t)
	mu.Unlock()

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

	mu.Lock()
	defer mu.Unlock()

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = !tasks[i].Done
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// 1. Remap IDs upon task deletion
	for i := range tasks {
		tasks[i].ID = i + 1
	}
	nextID = len(tasks) + 1

	w.WriteHeader(http.StatusNoContent)
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	groupMap := make(map[string]int)
	for _, t := range tasks {
		g := t.Group
		if g == "" {
			g = "default"
		}
		groupMap[g]++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupMap)
}
