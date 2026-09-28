package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-app/internal/task"
)

func TestRegisterStaticRoutesServesPWAAssetsPublicly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", RenderPage)
	RegisterStaticRoutes(mux)

	for _, path := range []string{"/manifest.json", "/sw.js", "/static/sw.js"} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Errorf("GET %s returned status %d", path, response.Code)
		}
	}

	manifestResponse := httptest.NewRecorder()
	mux.ServeHTTP(manifestResponse, httptest.NewRequest(http.MethodGet, "/manifest.json", nil))
	var manifest map[string]any
	if err := json.Unmarshal(manifestResponse.Body.Bytes(), &manifest); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}
	if got := manifestResponse.Header().Get("Content-Type"); !strings.Contains(got, "application/manifest+json") && !strings.Contains(got, "application/json") {
		t.Errorf("manifest content type = %q", got)
	}

	workerResponse := httptest.NewRecorder()
	mux.ServeHTTP(workerResponse, httptest.NewRequest(http.MethodGet, "/sw.js", nil))
	if got := workerResponse.Header().Get("Service-Worker-Allowed"); got != "/" {
		t.Errorf("service worker scope header = %q, want /", got)
	}

	pageResponse := httptest.NewRecorder()
	mux.ServeHTTP(pageResponse, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(pageResponse.Body.String(), `rel="manifest" href="/manifest.json"`) ||
		!strings.Contains(pageResponse.Body.String(), `navigator.serviceWorker.register('/sw.js')`) {
		t.Error("page does not link and register the PWA manifest and service worker")
	}
	if !strings.Contains(pageResponse.Body.String(), `hx-on::after-request="if (event.detail.successful) refreshAfterCreate(this)"`) ||
		!strings.Contains(pageResponse.Body.String(), `htmx.trigger('#todo-list', 'refreshTasks')`) {
		t.Error("successful task creation does not refresh the task list")
	}
}

func TestRenderItemHasClickableCompletionCheckbox(t *testing.T) {
	for _, test := range []struct {
		name      string
		done      bool
		isChecked bool
		label     string
	}{
		{name: "not done", label: `aria-label="Complete Write tests"`},
		{name: "done", done: true, isChecked: true, label: `aria-label="Mark Write tests as not done"`},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			RenderItem(response, task.Task{ID: 3, Title: "Write tests", Done: test.done})
			body := response.Body.String()
			if !strings.Contains(body, `type="checkbox"`) ||
				!strings.Contains(body, `hx-patch="/web/todos/3/toggle"`) ||
				!strings.Contains(body, `hx-trigger="change"`) ||
				!strings.Contains(body, test.label) {
				t.Fatalf("rendered item is missing the toggle checkbox behavior: %s", body)
			}
			if strings.Contains(body, "checked") != test.isChecked {
				t.Fatalf("checkbox checked state does not match done=%t: %s", test.done, body)
			}
		})
	}
}
