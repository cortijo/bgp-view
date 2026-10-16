package api

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// apiKey autentica por header X-API-Key (chave única via env BGP_API_KEY).
func apiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := os.Getenv("BGP_API_KEY")
		got := r.Header.Get("X-API-Key")
		if want == "" || subtle.ConstantTimeCompare([]byte(want), []byte(got)) != 1 {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
