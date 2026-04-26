package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Wild-sergunys/shrtic/internal/model"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"unauthorized","message":"Требуется авторизация"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"unauthorized","message":"Неверный формат токена"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"unauthorized","message":"Недействительный токен"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"unauthorized","message":"Неверные claims"}`, http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, `{"error":"unauthorized","message":"user_id не найден в токене"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), model.UserIDKey, int64(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuthMiddleware(jwtKey []byte) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
						return jwtKey, nil
					})
					if err == nil && token.Valid {
						if claims, ok := token.Claims.(jwt.MapClaims); ok {
							if userID, ok := claims["user_id"].(float64); ok {
								ctx := context.WithValue(r.Context(), model.UserIDKey, int64(userID))
								r = r.WithContext(ctx)
							}
						}
					}
				}
			}
			next(w, r)
		}
	}
}
