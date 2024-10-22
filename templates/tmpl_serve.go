package templates

import (
  "log"
  "net/http"
  "html/template"
  "github.com/jakeNorman007/go_chime/internals"
  "github.com/jakeNorman007/go_chime/auth/users"
)

func ServeChatTemplates(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/chat" {
    http.Error(w, "Files served to the root path not found", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Files served to the root path could not be fetched", http.StatusNotFound)
    return
  }

  cookie, err := r.Cookie("jwt")
  if err != nil {
    log.Println("Error finding JWT: ", err)
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  username, err := users.ExtractUsernameFromToken(cookie.Value)
  if err != nil {
    log.Println("Error decoding JWT: ", err)
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  tmpl, err := template.ParseFiles("templates/chat_body.html")
  if err != nil {
    log.Fatalf("Error parsing template: %s", err)
    http.Error(w, "Internal server error: ", http.StatusInternalServerError)
    return
  }

  data := struct {
    Username string
  }{
    Username: username,
  }

  err = tmpl.Execute(w, data)
  if err != nil {
    log.Fatalf("Error executing template: %s", err)
    http.Error(w, "Internal server error: ", http.StatusInternalServerError)
    return
  }
}

func ServeAuthenticationTemplates(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/auth" {
    http.Error(w, "Files served to the root path not found", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Files served to the root path could not be fetched", http.StatusNotFound)
    return
  }

  http.ServeFile(w, r, "templates/log_in.html")
}

func ServeMessageTemplates(w http.ResponseWriter, r *http.Request) {
  core := internals.NewCore()
  go core.Run()

  internals.ServeWebSocket(core, w, r);
}
