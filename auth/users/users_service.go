package users

import (
  "time"
  "context"
  "strconv"
  "github.com/golang-jwt/jwt/v5"
  "github.com/jakeNorman007/go_chime/auth/utils"
)

const (
  secretKey = "secret"
)

type service struct {
  Repo
  timeout time.Duration
}

func NewService(repo Repo) Service {
  return &service {
    repo,
    time.Duration(2) * time.Second,
  }
}

func (s *service) Createuser(c context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
  ctx, cancel := context.WithTimeout(c, s.timeout)
  defer cancel()

  hashedPassword, err := utils.HashPassword(request.Password)
  if err != nil {
    return nil, err
  }

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

type JWTClaims struct {
  ID        string `json:"id"`
  Username  string `json:"username"`
  jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, request *LoginUserRequest) (*LoginUserResponse, error) {
  ctx, cancel := context.WithTimeout(c, s.timeout)
  defer cancel()

  user, err := s.Repo.GetUserByEmail(ctx, request.Email)
  if err != nil {
    return &LoginUserResponse{}, err
  }

  err = utils.CheckPassword(request.Password, user.Password)
  if err != nil {
    return &LoginUserResponse{}, err
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims {
    ID: strconv.Itoa(int(user.ID)),
    Username: user.Username,
    RegisteredClaims: jwt.RegisteredClaims {
      Issuer: strconv.Itoa(int(user.ID)),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    },
  })

  ss, err := token.SignedString([]byte(secretKey))
  if err != nil {
    return &LoginUserResponse {}, err
  }

  return &LoginUserResponse { accessToken: ss, Username: user.Username, ID: strconv.Itoa(int(user.ID)) }, nil
}
