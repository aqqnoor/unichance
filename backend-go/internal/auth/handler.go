package auth

import (
  "net/http"

  "github.com/labstack/echo/v4"
)

type Handler struct { Svc Service }

type authReq struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

func (h Handler) Register(c echo.Context) error {
  var req authReq
  if err := c.Bind(&req); err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"error":"bad body"}) }
  token, user, err := h.Svc.Register(c.Request().Context(), req.Email, req.Password)
  if err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()}) }
  return c.JSON(http.StatusCreated, map[string]any{"token": token, "user": user})
}

func (h Handler) Login(c echo.Context) error {
  var req authReq
  if err := c.Bind(&req); err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"error":"bad body"}) }
  token, user, err := h.Svc.Login(c.Request().Context(), req.Email, req.Password)
  if err != nil { return c.JSON(http.StatusUnauthorized, map[string]string{"error":"invalid credentials"}) }
  return c.JSON(http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h Handler) Me(c echo.Context) error {
  u := c.Get("user")
  return c.JSON(http.StatusOK, map[string]any{"user": u})
}
