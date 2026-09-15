package middleware

import (
	"context"
	"net/http"
	"strings"

	"opsdesk/internal/config"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type contextKey string

const (
	UserContextKey contextKey = "currentUser"
)

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string

			// Try Authorization header first
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenStr = parts[1]
				}
			}

			// Try cookie as fallback
			if tokenStr == "" {
				if cookie, err := r.Cookie("token"); err == nil {
					tokenStr = cookie.Value
				}
			}

			if tokenStr == "" {
				utils.Error(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			claims, err := utils.ValidateToken(tokenStr, cfg.JWTSecret)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCurrentUser(r *http.Request) *utils.Claims {
	if val := r.Context().Value(UserContextKey); val != nil {
		if claims, ok := val.(*utils.Claims); ok {
			return claims
		}
	}
	return nil
}

func RequireRole(roles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetCurrentUser(r)
			if user == nil {
				utils.Error(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				utils.Error(w, http.StatusForbidden, "Access denied: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
