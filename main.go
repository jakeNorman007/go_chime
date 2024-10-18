package main

import (
  "log"
	_ "github.com/joho/godotenv/autoload"
  "github.com/jakeNorman007/go_chime/server"
  "github.com/jakeNorman007/go_chime/db"
)

func main() {
  databaseConnection := db.NewDatabaseConnection()
  db.InitDatabase(databaseConnection)

  server := server.NewService("PORT")
  if err := server.Run(); err != nil {
      log.Fatal(err)
  }
}
