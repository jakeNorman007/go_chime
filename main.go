package main

import (
  "log"
	_ "github.com/joho/godotenv/autoload"
  "github.com/jakeNorman007/go_chime/server"
  "github.com/jakeNorman007/go_chime/db"
  "github.com/jakeNorman007/go_chime/auth/users"
  "github.com/jakeNorman007/go_chime/auth/utils"
)

func main() {
  databaseConnection := db.NewDatabaseConnection()
  db.InitDatabase(databaseConnection)

  userRepo := users.NewRepo(databaseConnection)
  userService := users.NewService(userRepo)

  utils.PasswordCheck()

  server := server.NewService("PORT", userService)
  if err := server.Run(); err != nil {
      log.Fatal(err)
  }
}
