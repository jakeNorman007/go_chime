package utils

import (
  "fmt"
  "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return "", fmt.Errorf("Failed to hash password: %w", err)
  }

  return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, password string) error {
  err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
  if err != nil {
    return err
  }

  return err
}

func PasswordCheck() {
    password := "yourPassword"
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Println("Error hashing password:", err)
        return
    }

    fmt.Println("Hashed password:", string(hashedPassword))

    // Simulate storing and retrieving from the database
    storedHash := string(hashedPassword) // Replace with actual DB retrieval
    err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
    if err != nil {
        fmt.Println("Password verification failed:", err)
    } else {
        fmt.Println("Password verification succeeded.")
    }
}
