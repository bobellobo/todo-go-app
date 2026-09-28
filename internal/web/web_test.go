package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
}
