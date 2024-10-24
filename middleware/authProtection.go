package middleware

import (
  "net/http"
  "github.com/jakeNorman007/go_chime/auth/users"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/login" || r.URL.Path == "/signup" {
      next.ServeHTTP(w, r)
      return
    }

    cookie, err := r.Cookie("jwt")
    if err != nil {
      http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
      return
    }

    _, err = users.ExtractUsernameFromToken(cookie.Value)
    if err != nil {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    next.ServeHTTP(w, r)
  })
}
