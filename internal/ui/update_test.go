package ui

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-app/internal/client"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRefreshKeyFetchesTasks(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/tasks" {
			http.Error(w, "unexpected request", http.StatusNotFound)
			return
		}
		requests++
		_, _ = io.WriteString(w, `[{"id":7,"title":"Fresh task","done":false}]`)
	}))
	defer server.Close()

	model := NewModel(client.NewClient(server.URL, ""))
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	if cmd == nil {
		t.Fatal("refresh key did not start a task fetch")
	}

	msg := cmd()
	if _, ok := msg.(tasksMsg); !ok {
		t.Fatalf("expected task response, got %T", msg)
	}
	updated, _ := model.Update(msg)
	refreshed, ok := updated.(Model)
	if !ok {
		t.Fatalf("expected updated Model, got %T", updated)
	}
	if requests != 1 {
		t.Fatalf("expected one task refresh request, got %d", requests)
	}
	if len(refreshed.tasks) != 1 || refreshed.tasks[0].Title != "Fresh task" {
		t.Fatalf("expected refreshed task list, got %+v", refreshed.tasks)
	}
	if !strings.Contains(refreshed.View(), "r: refresh") {
		t.Fatal("refresh shortcut is missing from the TUI help")
	}
}
