package main

import (
  "context"
  "log"

  "github.com/joho/godotenv"

  "unichance-backend-go/internal/auth"
  "unichance-backend-go/internal/config"
  "unichance-backend-go/internal/db"
  httpRouter "unichance-backend-go/internal/http"
  "unichance-backend-go/internal/programs"
  "unichance-backend-go/internal/profile"
)

func main() {
  _ = godotenv.Load(".env")

  cfg := config.Load()
  if cfg.DatabaseURL == "" || cfg.JwtSecret == "" {
    log.Fatal("DATABASE_URL and JWT_SECRET are required")
  }

  pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
  if err != nil { log.Fatal(err) }
  defer pool.Close()

  // auth
  authSvc := auth.Service{DB: pool, JwtSecret: cfg.JwtSecret}
  authH := auth.Handler{Svc: authSvc}

  // programs
  progRepo := programs.Repo{DB: pool}
  progH := programs.Handler{Repo: progRepo}

  // profile + scoring endpoints
  profRepo := profile.Repo{DB: pool}
  profH := profile.Handler{Repo: profRepo, DB: pool}

  e := httpRouter.NewRouter(httpRouter.Deps{
    AuthHandler:     authH,
    ProgramsHandler: progH,
    ProfileHandler:  profH,
    JwtSecret:       cfg.JwtSecret,
  })

  log.Println("api listening on :" + cfg.Port)
  log.Fatal(e.Start(":" + cfg.Port))
}
