package server

import (
  "os"
  "fmt"
  "log"
  "time"
  "strconv"
  "net/http"
  "github.com/rs/cors"
  "github.com/jakeNorman007/go_chime/templates"
  "github.com/jakeNorman007/go_chime/middleware"
  "github.com/jakeNorman007/go_chime/internals"
)

type Service struct {
    addr    string
}

func NewService(addr string) *Service {
    port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

    return &Service {
        addr:   fmt.Sprintf(":%d", port),
    }
}

func (s *Service) Run() error {
  router := http.NewServeMux()
  core := internals.NewCore()

  go core.Run()

  c := cors.Default().Handler(middleware.Logging(router))

  server := http.Server {
    Addr:           s.addr,
    Handler:        c,
    IdleTimeout:    time.Minute,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   30 * time.Second,
  }

  router.HandleFunc("/", templates.ServeIndexTemplates)
  router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    internals.ServeWebSocket(core, w, r)
  })

  log.Println("Server listeneing on port", s.addr)
  return server.ListenAndServe()
}
