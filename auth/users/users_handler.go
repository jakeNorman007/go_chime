package users

import (
  "net/http"
  "encoding/json"
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
  var userRequest CreateUserRequest

  if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  response, err := h.Service.CreateUser(r.Context(), &userRequest)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(response); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
  var userLoginRequest LoginUserRequest

  if err := json.NewDecoder(r.Body).Decode(&userLoginRequest); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  user, err := h.Service.Login(r.Context(), &userLoginRequest)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
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

  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(user); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie {
    Name: "jwt",
    Value: "", 
    Path: "/",
    MaxAge: -1,
    HttpOnly: true,
  })

  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(map[string]string{"Message": "logout successful"}); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
		return
  }
}
