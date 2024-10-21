package middleware

import (
  "net/http"
  //"encoding/json"
  "github.com/golang-jwt/jwt/v4"
  //"github.com/jakeNorman007/go_chime/auth/users"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
  //var userLoginRequest users.LoginUserRequest
  return func(w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
      http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
      return
    }

    tokenString = tokenString[len("Bearer "):]

    //claims := &jwt.MapClaims{}

    return
  }
}
