package users

import (
  "database/sql"
)

type Store struct {
  db *sql.DB
}

func NewStore(db *sql.DB) *Store {
  return &Store {
    db: db,
  }
}

func (s *Store) CreateUser(user *User) error {
  _, err := s.db.Query("INSERT INTO users (email, username, password) VALUES($1, $2, $3)", user.Email, user.Username, user.Password)
  if err != nil {
    return err
  }

  return nil
}

func (s *Store) GetAllUsers() ([]*User, error) {
  rows, err := s.db.Query("SELECT * FROM users")
  if err != nil {
    return nil, err
  }
  
  users := make([]*User, 0)
  for rows.Next() {
    p, err := ScanRowsIntoUser(rows)
    if err != nil {
      return nil, err
    }

    users = append(users, p)
  }

  return users, nil
}

func (s *Store) GetUserById(id int) (*User, error) {
  rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
  if err != nil {
    return nil, err
  }

  p := new(User)

  for rows.Next() {
    p, err = ScanRowsIntoUser(rows)
    if err != nil {
      return nil, err
    }
  }

  return p, nil
}

func (s *Store) UpdateUserById(id int, email string, username string, password string) error {
  _, err := s.db.Query("UPDATE users SET email = $1, username = $2, password = $3", email, username, password)
  if err != nil {
    return err
  }

  return nil
}

func(s *Store) DeleteUser(id int) error {
  _, err := s.db.Query("DELETE FROM users WHERE id = $1", id)
  return err
}

func ScanRowsIntoUser(rows *sql.Rows) (*User, error) {
  user := new(User)

  err := rows.Scan(&user.ID, &user.Username, &user.Email)
  if err != nil {
    return nil, err
  }

  return user, nil
}
