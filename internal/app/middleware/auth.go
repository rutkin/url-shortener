package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/service"
	"go.uber.org/zap"
)

const password = "password"

// error user not found
var ErrNotFound = errors.New("userID not found")

// get user id from cookie and decrypt
func GetUserIDFromCookie(r *http.Request) (string, error) {
	userIDCookie, err := r.Cookie("userID")
	if err != nil {
		return "", ErrNotFound
	}

	key := sha256.Sum256([]byte(password))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		logger.Log.Error("failed to create new cipher", zap.String("error", err.Error()))
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		logger.Log.Error("failed to create new gcm", zap.String("error", err.Error()))
		return "", err
	}

	nonce := key[(len(key) - aesgcm.NonceSize()):]

	data, err := hex.DecodeString(userIDCookie.Value)
	if err != nil {
		logger.Log.Error("failed to decode userID cookie", zap.String("error", err.Error()))
		return "", err
	}

	userID, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		logger.Log.Error("failed to decrypt userID cookie", zap.String("error", err.Error()))
		return "", err
	}
	return string(userID), err
}

// set user id to cookie and encrypt
func SetUserIDToCookies(w http.ResponseWriter, userID string) error {
	key := sha256.Sum256([]byte(password))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		logger.Log.Error("failed to create new cipher", zap.String("error", err.Error()))
		return err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		logger.Log.Error("failed to create new gcm", zap.String("error", err.Error()))
		return err
	}

	nonce := key[(len(key) - aesgcm.NonceSize()):]

	encryptedUserID := aesgcm.Seal(nil, nonce, []byte(userID), nil)
	userIDcookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString(encryptedUserID)}
	http.SetCookie(w, userIDcookie)
	return nil
}

// middleware that set user id
func WithUserID(h http.Handler) http.Handler {
	authFn := func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if errors.Is(err, ErrNotFound) {
			userID = uuid.NewString()
			err = SetUserIDToCookies(w, userID)
		}

		if err != nil {
			logger.Log.Error("failed to set user id", zap.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}

		logger.Log.Info("WithUserID", zap.String("userID", userID))
		ctx := context.WithValue(r.Context(), service.UserIDKey, userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(authFn)
}

// middleware that check userid exists
func WithAuth(h http.Handler) http.Handler {
	authFn := func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), service.UserIDKey, userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(authFn)
}
