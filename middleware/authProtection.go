package middleware

import (
  "strings"
  "context"
  "net/http"
  "github.com/golang-jwt/jwt/v4"
  "github.com/jakeNorman007/go_chime/auth/users"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    authenticationHeader := r.Header.Get("Authorization")

    if !strings.HasPrefix(authenticationHeader, "Bearer ") {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    tokenString := strings.TrimPrefix(authenticationHeader, "Bearer ")

    claims := &users.JWTClaims {}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface {}, error) {
      return []byte(users.SecretKey), nil
    })

    if err != nil || !token.Valid {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    ctx := context.WithValue(r.Context(), "ID", claims.ID)
    ctx = context.WithValue(ctx, "Username", claims.Username)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}
