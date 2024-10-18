package users

type UsersStore interface {
  CreateUser(User) error
  GetAllUsers() ([]*User, error)
  UpdateUserById(id int, email string, username string) error
  DeleteUser(id int) error
}

type User struct {
  ID        int     `json:"id"`
  Email     string  `json:"email"`
  Username  string  `json:"username"`
  Password  string  `json:"password"`
}

type CreateUserRequest struct {
  Email     string  `json:"email"`
  Username  string  `json:"username"`
  Password  string  `json:"password"`
}

type CreateUserResponse struct {
  ID        string  `json:"id"`
  Email     string  `json:"email"`
  Username  string  `json:"username"`
}

type LoginUserRequest struct {
  Email     string  `json:"email"`
  Password  string  `json:"password"`
}

type LoginUserResponse struct {
  accessToken string
  ID          string  `json:"id"`
  Username  string  `json:"username"`
}

func NewUser(email string, username string) *User {
  return &User {
    Email: email,
    Username: username,
  }
}
