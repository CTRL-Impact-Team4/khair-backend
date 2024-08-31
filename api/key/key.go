package key

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Service       string
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 10) // Generates a 20-character hex string
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return "cml-" + hex.EncodeToString(bytes), nil
}

func GenKey(info UserInfo) (string, error) {
	key, _ := generateAPIKey()
	// TODO: add to db with userinfo
	return key, nil
}

func ValidateKey(api_key string) (bool, error) {
	// TODO: check db that key exists and is valid
	// (not deleted by user and below usage cap)
	// -> probably needs usage cap func
	return true, nil
}

func InvalidateKey(api_key string) error {
	// TODO: in db invalidate (delete from user perspective )
	return nil
}

func HandleDeleteKey(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Authorization header must be in 'Bearer {token}' format", http.StatusBadRequest)
		return
	}

	apiKey := parts[1]

	v, err := ValidateKey(apiKey)

	if err != nil {
		http.Error(w, "Error validating CAMLL API key", http.StatusUnauthorized)
		return
	}

	if !v {
		http.Error(w, "Invalid CAMLL API key", http.StatusUnauthorized)
		return
	}

	err = InvalidateKey(apiKey)
	if err != nil {
		log.Printf("Failed to invalidate API key: %s, error: %v", apiKey, err)
		http.Error(w, "Failed to invalidate API key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API key invalidated successfully"))
}

// TODO: billing info and usage endpoints
