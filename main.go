package main

import (
  "log"
	_ "github.com/joho/godotenv/autoload"
  "github.com/jakeNorman007/go_chime/db"
  "github.com/jakeNorman007/go_chime/server"
  "github.com/jakeNorman007/go_chime/auth/users"
)

func main() {
  databaseConnection := db.NewDatabaseConnection()
  db.InitDatabase(databaseConnection)

  userRepo := users.NewRepo(databaseConnection)
  userService := users.NewService(userRepo)

  server := server.NewService("PORT", userService)
  if err := server.Run(); err != nil {
      log.Fatal(err)
  }
}
