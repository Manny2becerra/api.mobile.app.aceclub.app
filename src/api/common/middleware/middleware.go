package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"context"

	"api-mobile-app/src/api/common/constants"
	errorHandling "api-mobile-app/src/api/common/errorHandling"
	Logging "api-mobile-app/src/api/common/logging"

	auth "api-mobile-app/src/api/common/authentication"
)

type AuthRequest struct {
	AccessToken string          `json:"accessToken"`
	Data        json.RawMessage `json:"data"`
}

func Middleware(router http.Handler, JWTAccessSecret string, logger *Logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timing the request
		start := time.Now()
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.StartTimeKey, start)
		r = r.WithContext(ctx)

		// Check if this is a multipart form request
		contentType := r.Header.Get("Content-Type")
		isMultipartForm := strings.Contains(contentType, "multipart/form-data")

		var claims *auth.AccessClaims
		var err error

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			var authReq AuthRequest
			if isMultipartForm {
				// For multipart/form-data requests, get token from Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					errorDetails := errorHandling.Unauthorized("Missing authorization header", nil)
					logger.LogRequest(r, &errorDetails)
					http.Error(w, "not authorized", http.StatusUnauthorized)
					return
				}

				claims, err = auth.Authenticate(authHeader, JWTAccessSecret)
			} else {
				// For JSON requests, get token from body
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					errorDetails := errorHandling.BadRequest("Failed to read request body", nil)
					logger.LogRequest(r, &errorDetails)
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				r.Body.Close()

				// Create a copy of the original body for later
				bodyReader := bytes.NewBuffer(bodyBytes)

				// Parse the auth request
				if err := json.NewDecoder(bodyReader).Decode(&authReq); err != nil {
					errorDetails := errorHandling.BadRequest("Invalid request body format", nil)
					logger.LogRequest(r, &errorDetails)
					http.Error(w, "not authorized", http.StatusBadRequest)
					return
				}

				// Check if token exists
				if authReq.AccessToken == "" {
					errorDetails := errorHandling.Unauthorized("Missing authentication token", nil)
					logger.LogRequest(r, &errorDetails)
					http.Error(w, "not authorized", http.StatusUnauthorized)
					return
				}

				// Authenticate the token
				claims, err = auth.Authenticate(authReq.AccessToken, JWTAccessSecret)
				if err != nil {
					errorDetails := errorHandling.Unauthorized("Invalid authentication token", nil)
					logger.LogRequest(r, &errorDetails)
					http.Error(w, "not authorized", http.StatusUnauthorized)
					return
				}

				// Create a new reader with the original body
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// Handle authentication errors
			if err != nil {
				var detail string
				if err == auth.ErrExpiredToken {
					detail = "Token has expired"
				} else if err == auth.ErrInvalidToken {
					detail = "Invalid token format"
				} else {
					detail = "Authentication failed"
				}

				errorDetails := errorHandling.Unauthorized(detail, nil)
				logger.LogRequest(r, &errorDetails)
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}
			// Add authenticated identity to context
			if claims == nil {
				errorDetails := errorHandling.Unauthorized("no claims found", nil)
				logger.LogRequest(r, &errorDetails)
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, constants.UserEmailKey, claims.Subject)
			ctx = context.WithValue(ctx, constants.UserRoleKey, claims.Role)
			ctx = context.WithValue(ctx, constants.AccessToken, authReq.AccessToken)
			if claims.EstablishmentID != "" {
				ctx = context.WithValue(ctx, constants.EstablishmentIDKey, claims.EstablishmentID)
			}
			r = r.WithContext(ctx)
		} else {
			// For other methods (GET, DELETE), use header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errorDetails := errorHandling.Unauthorized("Missing authorization header", nil)
				logger.LogRequest(r, &errorDetails)
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			// Authenticate the token from header
			claims, err = auth.Authenticate(authHeader, JWTAccessSecret)
			if err != nil {
				var detail string
				if err == auth.ErrExpiredToken {
					detail = "Token has expired"
				} else if err == auth.ErrInvalidToken {
					detail = "Invalid token format"
				} else {
					detail = "Authentication failed"
				}

				errorDetails := errorHandling.Unauthorized(detail, nil)
				logger.LogRequest(r, &errorDetails)
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			// Add authenticated identity to context
			ctx = context.WithValue(ctx, constants.UserEmailKey, claims.Subject)
			ctx = context.WithValue(ctx, constants.UserRoleKey, claims.Role)
			if claims.EstablishmentID != "" {
				ctx = context.WithValue(ctx, constants.EstablishmentIDKey, claims.EstablishmentID)
			}
			r = r.WithContext(ctx)
		}

		// Continue to the next handler
		router.ServeHTTP(w, r)
	})
}
