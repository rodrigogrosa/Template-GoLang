package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// AuthConfig holds JWT authentication configuration
type AuthConfig struct {
	Secret    string
	PublicKey *rsa.PublicKey // For RS256
	Issuer    string
	Enabled   bool
}

type contextKey string

const userContextKey contextKey = "user"

// JWTAuth middleware validates JWT tokens
func JWTAuth(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if auth is disabled
			if !config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			// Remove Bearer prefix
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, `{"error":"Invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Secret), nil
			})

			if err != nil {
				log.Warn().Err(err).Msg("Invalid token")
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"Invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			// Validate issuer if configured
			if config.Issuer != "" {
				if iss, ok := claims["iss"].(string); !ok || iss != config.Issuer {
					http.Error(w, `{"error":"Invalid token issuer"}`, http.StatusUnauthorized)
					return
				}
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extracts user claims from context
func GetUserFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(userContextKey).(jwt.MapClaims)
	return claims, ok
}
