package users

import (
  "net/http"
)

type Handler struct {
  Service
}

func NewHandler(s Service) *Handler {
  return &Handler {
    Service: s,
  }
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Invalid http method for creating user.", http.StatusMethodNotAllowed)
    return
  }

  if err := r.ParseForm(); err != nil {
    http.Error(w, "Failed to parse sign up form.", http.StatusBadRequest)
  }

  createUserRequest := CreateUserRequest {
    Username: r.FormValue("username"),
    Email: r.FormValue("email"),
    Password: r.FormValue("password"),
  }

  _, err := h.Service.CreateUser(r.Context(), &createUserRequest)
  if err != nil {
    http.Redirect(w, r, "/unauthorized.html", http.StatusUnauthorized)
  }

  w.Header().Set("HX-Redirect", "/log_in")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Invalid http method for logging in user.", http.StatusMethodNotAllowed)
    return
  }

  if err := r.ParseForm(); err != nil {
    http.Error(w, "Failed to parse log in form", http.StatusBadRequest)
  }

  userLoginRequest := LoginUserRequest {
    Email: r.FormValue("email"),
    Password: r.FormValue("password"),
  }

  user, err := h.Service.Login(r.Context(), &userLoginRequest)
  if err != nil {
    http.Redirect(w, r, "/unauthorized.html", http.StatusUnauthorized)
    return
  }

  http.SetCookie(w, &http.Cookie {
    Name: "jwt",
    Value: user.accessToken,
    Path: "/",
    MaxAge: 60 * 60 * 24,
    HttpOnly: true,
    Secure: false,
    SameSite: http.SameSiteLaxMode,
  })

  w.Header().Set("HX-Redirect", "/chat")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie {
    Name: "jwt",
    Value: "", 
    Path: "/",
    MaxAge: -1,
    HttpOnly: true,
  })

  http.Redirect(w, r, "/log_in", http.StatusSeeOther)
}
