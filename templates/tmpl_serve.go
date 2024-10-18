package templates

import (
  "net/http"
  "github.com/jakeNorman007/go_chime/internals"
)

func ServeIndexTemplates(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.Error(w, "Files served to the root path not found", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Files served to the root path could not be fetched", http.StatusNotFound)
    return
  }

  http.ServeFile(w, r, "templates/index.html")
}

func ServeMessageTemplates(w http.ResponseWriter, r *http.Request) {
  core := internals.NewCore()
  go core.Run()

  internals.ServeWebSocket(core, w, r);
}
