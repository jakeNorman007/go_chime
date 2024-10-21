package users

import (
  "time"
  "context"
  "strconv"
  "net/http"
  "github.com/golang-jwt/jwt/v4"
  "github.com/jakeNorman007/go_chime/auth/utils"
)

const (
  SecretKey = "secret"
)

type service struct {
  Repo
  timeout time.Duration
}

type JWTClaims struct {
  ID        string      `json:"id"`
  Username  string      `json:"username"`

  jwt.RegisteredClaims
}

func NewService(repo Repo) *service {
  return &service {
    repo,
    time.Duration(2) * time.Second,
  }
}

func (s *service) CreateUser(c context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
  ctx, cancel := context.WithTimeout(c, s.timeout)
  defer cancel()

  //log.Printf("Creating user with password: %s", request.Password)

  hashedPassword, err := utils.HashPassword(request.Password)
  if err != nil {
    return nil, err
  }

  //log.Printf("Hashed password: %s", hashedPassword)
  //log.Printf("Hashed password length: %d", len(hashedPassword))

  user := &User {
    Username: request.Username,
    Email: request.Email,
    Password: hashedPassword,
  }

  r, err := s.Repo.CreateUser(ctx, user)
  if err != nil {
    return nil, err
  }

  response := &CreateUserResponse {
    ID: strconv.Itoa(int(r.ID)),
    Username: r.Username,
    Email: r.Email,
  }

  return response, nil
}

func (s *service) Login(c context.Context, request *LoginUserRequest) (*LoginUserResponse, error) {
  ctx, cancel := context.WithTimeout(c, s.timeout)
  defer cancel()

  user, err := s.Repo.GetUserByEmail(ctx, request.Email)
  if err != nil {
    return &LoginUserResponse{}, err
  }

  //log.Printf("Retrieved user password: %s", user.Password)

  err = utils.CheckPassword(request.Password, user.Password)
  if err != nil {
    //log.Printf("Request Password: %v", request.Password)
    //log.Printf("User Password: %v", user.Password)
    //log.Printf("Password check failed: %v", err)
    return &LoginUserResponse{}, err
  }

  //log.Printf("Request password length: %d", len(request.Password))
  //log.Printf("Request password: %v", request.Password)
  //log.Printf("User password length encrypted: %d", len(user.Password))

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims {
    ID: strconv.Itoa(int(user.ID)),
    Username: user.Username,
    RegisteredClaims: jwt.RegisteredClaims {
      Issuer: strconv.Itoa(int(user.ID)),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    },
  })

  ss, err := token.SignedString([]byte(SecretKey))
  if err != nil {
    return &LoginUserResponse {}, err
  }

  return &LoginUserResponse { accessToken: ss, Username: user.Username, ID: strconv.Itoa(int(user.ID)) }, nil
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
  userID := r.Context().Value("ID").(string)
  username := r.Context().Value("Username").(string)

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Hi " + username + ", your user ID is " + userID))
}
