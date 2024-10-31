package users

import (
  "context"
  "database/sql"
)

type DBTX interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
  	PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type repo struct {
  db DBTX
}

func NewRepo(db DBTX) *repo {
  return &repo {
    db: db,
  }
}

func (r *repo) CreateUser(ctx context.Context, user *User) (*User, error) {
  var lastInsertId int

  query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
  err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
  if err != nil {
    return &User{}, err
  }

  user.ID = int64(lastInsertId)

  return user, nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
  user := User{}

  query := "SELECT id, email, username, password FROM users WHERE email = $1"
  err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
  if err != nil {
    return &User{}, nil
  }

  return &user, nil
}

func (r *repo) GetAllUsers(ctx context.Context) ([]*User, error) {
  query := "SELECT username FROM users"

  rows, err := r.db.QueryContext(ctx, query)
  if err != nil {
    return nil, err
  }

  defer rows.Close()

  var users []*User
  for rows.Next() {
    user := new(User)
    if err := rows.Scan(&user.Username); err != nil {
      return nil, err
    }

    users = append(users, user)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return users, nil
}
