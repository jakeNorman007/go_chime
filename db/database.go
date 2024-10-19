package db

import (
  "os"
  "fmt"
  "log"
  "database/sql"
  _ "github.com/joho/godotenv/autoload"
  _ "github.com/jackc/pgx/v5/stdlib"
)

type dbService struct {
  db *sql.DB
}

func NewDatabaseConnection() *sql.DB {
  connectionString := fmt.Sprintf(os.Getenv("DB_URL"))

  db, err := sql.Open("pgx", connectionString)
  if err != nil {
    log.Fatal(err)
  }

  return db
}

func InitDatabase(db *sql.DB) {
  err := db.Ping()

  if err != nil {
    log.Fatal(err)
  }

  log.Println("Connection to database established...")
}
