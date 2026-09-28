package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateWebTodoStoresAndRendersGroup(t *testing.T) {
	testDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	if _, err := testDB.Exec(`
		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0,
			task_group TEXT DEFAULT ''
		)
	`); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = testDB
	defer func() { db = previousDB }()

	form := url.Values{
		"title": {"Plan the garden"},
		"group": {"  personal  "},
	}
	request := httptest.NewRequest(http.MethodPost, "/web/todos", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	createWebTodo(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var title, group string
	if err := testDB.QueryRow("SELECT title, task_group FROM tasks").Scan(&title, &group); err != nil {
		t.Fatal(err)
	}
	if title != "Plan the garden" {
		t.Errorf("expected title to be stored, got %q", title)
	}
	if group != "personal" {
		t.Errorf("expected trimmed group to be stored, got %q", group)
	}
	if !strings.Contains(response.Body.String(), `class="group-badge">personal</span>`) {
		t.Errorf("expected response to render the task group, got %s", response.Body.String())
	}

	groupsResponse := httptest.NewRecorder()
	getWebGroups(groupsResponse, httptest.NewRequest(http.MethodGet, "/web/groups", nil))
	if groupsResponse.Code != http.StatusOK {
		t.Fatalf("expected groups status %d, got %d: %s", http.StatusOK, groupsResponse.Code, groupsResponse.Body.String())
	}
	var groups []string
	if err := json.NewDecoder(groupsResponse.Body).Decode(&groups); err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || groups[0] != "personal" {
		t.Errorf("expected the group filter to offer personal, got %v", groups)
	}
}
