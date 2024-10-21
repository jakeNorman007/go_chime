package utils

import (
  "fmt"
  "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
