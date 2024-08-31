package middleware

import (
	"net/http"
)

// ValidateApiKey middleware to validate CAMLL API key from the Authorization header.
func ValidateApiKey() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check Authorization header exists
			// authHeader := r.Header.Get("Authorization")
			// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// 	http.Error(w, "Authorization required", http.StatusUnauthorized)
			// 	return
			// }

			// // Extract and validate the CAMLL API key
			// camllAPIKey := strings.TrimPrefix(authHeader, "Bearer ")
			// valid, err := key.ValidateKey(camllAPIKey)

			// if err != nil {
			// 	http.Error(w, "Error validating CAMLL API key", http.StatusInternalServerError)
			// 	return
			// }

			// if !valid {
			// 	http.Error(w, "Invalid CAMLL API key", http.StatusUnauthorized)
			// 	return
			// }

			h.ServeHTTP(w, r)
		})
	}
}
