package http

import (
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"

	"unichance-backend-go/internal/auth"
	appMw "unichance-backend-go/internal/middleware"
	"unichance-backend-go/internal/programs"
	"unichance-backend-go/internal/profile"
	"unichance-backend-go/internal/shortlist"
)

type Deps struct {
	AuthHandler      auth.Handler
	ProgramsHandler  programs.Handler
	ProfileHandler   profile.Handler
	ShortlistHandler shortlist.Handler
	JwtSecret        string
}


func NewRouter(d Deps) *echo.Echo {
	e := echo.New()

	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	e.GET("/health", func(c echo.Context) error { return c.String(200, "ok") })

	// auth (public)
	e.POST("/auth/register", d.AuthHandler.Register)
	e.POST("/auth/login", d.AuthHandler.Login)

	// auth/me (protected)
	e.GET("/auth/me", d.AuthHandler.Me, appMw.RequireAuth(d.JwtSecret))

	// programs (public)
	e.GET("/programs", d.ProgramsHandler.List)

	// profile (protected)
	e.GET("/profile/me", d.ProfileHandler.GetMe, appMw.RequireAuth(d.JwtSecret))
	e.POST("/profile/me", d.ProfileHandler.UpsertMe, appMw.RequireAuth(d.JwtSecret))
	e.POST("/score", d.ProfileHandler.ScoreProgram, appMw.RequireAuth(d.JwtSecret))

	// shortlist (protected)  <-- middleware дегенді appMw деп жаз
	e.GET("/shortlist", d.ShortlistHandler.List, appMw.RequireAuth(d.JwtSecret))
	e.POST("/shortlist", d.ShortlistHandler.Add, appMw.RequireAuth(d.JwtSecret))
	e.DELETE("/shortlist/:program_id", d.ShortlistHandler.Remove, appMw.RequireAuth(d.JwtSecret))

	return e
}

