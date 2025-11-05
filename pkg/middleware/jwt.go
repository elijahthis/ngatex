package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secretKey []byte
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: []byte(secretKey),
	}
}

type contextKey string

const claimsKey contextKey = "claims"

func (j *JWTAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Forbidden: Missing Authorization Header", http.StatusForbidden)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.secretKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
