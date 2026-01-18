package auth

import (
  "context"
  "time"

  "github.com/golang-jwt/jwt/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  "golang.org/x/crypto/bcrypt"
)

type Service struct {
  DB *pgxpool.Pool
  JwtSecret string
}

type User struct {
  ID string `json:"id"`
  Email string `json:"email"`
}

func (s Service) Register(ctx context.Context, email, password string) (string, User, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil { return "", User{}, err }

  var id string
  err = s.DB.QueryRow(ctx,
    `INSERT INTO users(email,password_hash) VALUES ($1,$2) RETURNING id`,
    email, string(hash),
  ).Scan(&id)
  if err != nil { return "", User{}, err }

  token, err := s.issueToken(id, email)
  if err != nil { return "", User{}, err }

  return token, User{ID:id, Email:email}, nil
}

func (s Service) Login(ctx context.Context, email, password string) (string, User, error) {
  var id, hash string
  err := s.DB.QueryRow(ctx,
    `SELECT id, password_hash FROM users WHERE email=$1`,
    email,
  ).Scan(&id, &hash)
  if err != nil { return "", User{}, err }

  if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
    return "", User{}, err
  }

  token, err := s.issueToken(id, email)
  if err != nil { return "", User{}, err }

  return token, User{ID:id, Email:email}, nil
}

func (s Service) issueToken(userID, email string) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
    "email": email,
    "iat": time.Now().Unix(),
    "exp": time.Now().Add(7*24*time.Hour).Unix(),
  }
  t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return t.SignedString([]byte(s.JwtSecret))
}
