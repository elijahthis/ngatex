package middleware

import "net/http"

type APIKeyAuth struct {
	validKeys map[string]bool
}

func NewAPIKeyAuth(keys []string) *APIKeyAuth {
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}
	return &APIKeyAuth{
		validKeys: keyMap,
	}
}

func (a *APIKeyAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, "Forbidden: Missing API Key", http.StatusForbidden)
			return
		}

		if !a.validKeys[key] {
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
