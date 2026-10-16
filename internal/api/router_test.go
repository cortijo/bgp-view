package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/bgpadvert/bgpadvert/internal/store"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewRouter(st)
}

func TestHealth(t *testing.T) {
	r := newTestRouter(t)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("/health = %d, esperado 200", rec.Code)
	}
}

func TestReady(t *testing.T) {
	r := newTestRouter(t)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("/ready = %d, esperado 200", rec.Code)
	}
}

func TestAPIKey(t *testing.T) {
	t.Setenv("BGP_API_KEY", "segredo")
	h := apiKey(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	// sem chave
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/prefixes", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("sem chave = %d, esperado 401", rec.Code)
	}

	// com chave
	req := httptest.NewRequest(http.MethodGet, "/v1/prefixes", nil)
	req.Header.Set("X-API-Key", "segredo")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("com chave = %d, esperado 204", rec.Code)
	}
}
