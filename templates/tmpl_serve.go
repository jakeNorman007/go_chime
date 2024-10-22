package templates

import (
  "net/http"
  "github.com/jakeNorman007/go_chime/internals"
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

  http.ServeFile(w, r, "templates/chat_body.html")
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
