package middleware

import (
  "net/http"
  "github.com/jakeNorman007/go_chime/auth/users"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Allow access to the /login and /signup routes
    if r.URL.Path == "/login" || r.URL.Path == "/signup" {
      next.ServeHTTP(w, r)
      return
    }

    // Get the JWT from the cookie
    cookie, err := r.Cookie("jwt")
    if err != nil {
      // No JWT cookie found, return unauthorized
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    // Validate the JWT here (using your existing ExtractUsernameFromToken or similar logic)
    _, err = users.ExtractUsernameFromToken(cookie.Value)
    if err != nil {
      // JWT is invalid or expired
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    // If we reached this point, the user is authenticated
    next.ServeHTTP(w, r)
  })
}
