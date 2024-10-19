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
  "github.com/jakeNorman007/go_chime/internals"
  "github.com/jakeNorman007/go_chime/middleware"
  "github.com/jakeNorman007/go_chime/auth/users"
)

type Service struct {
    addr          string
    usersService  users.Service
}

func NewService(addr string, usersService users.Service) *Service {
    port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

    return &Service {
        addr:   fmt.Sprintf(":%d", port),
        usersService: usersService,
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

  usersHandler := users.NewHandler(s.usersService)

  router.HandleFunc("POST /signup", usersHandler.CreateUser)
  router.HandleFunc("POST /login", usersHandler.Login)
  router.HandleFunc("POST /logout", usersHandler.Logout)

  router.HandleFunc("/", templates.ServeIndexTemplates)

  router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    internals.ServeWebSocket(core, w, r)
  })

  log.Println("Server listeneing on port", s.addr)
  return server.ListenAndServe()
}
