package shortlist

import (
  "net/http"
  "github.com/labstack/echo/v4"
  "unichance-backend-go/internal/middleware"
)

type Handler struct { Repo Repo }

type addReq struct { ProgramID string `json:"program_id"` }

func (h Handler) List(c echo.Context) error {
  u := c.Get("user").(middleware.CtxUser)
  items, err := h.Repo.List(c.Request().Context(), u.ID)
  if err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()}) }
  return c.JSON(http.StatusOK, map[string]any{"items": items})
}

func (h Handler) Add(c echo.Context) error {
  u := c.Get("user").(middleware.CtxUser)
  var req addReq
  if err := c.Bind(&req); err != nil || req.ProgramID == "" {
    return c.JSON(http.StatusBadRequest, map[string]string{"error":"program_id required"})
  }
  if err := h.Repo.Add(c.Request().Context(), u.ID, req.ProgramID); err != nil {
    return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
  }
  return c.NoContent(http.StatusNoContent)
}

func (h Handler) Remove(c echo.Context) error {
  u := c.Get("user").(middleware.CtxUser)
  id := c.Param("program_id")
  if id == "" {
    return c.JSON(http.StatusBadRequest, map[string]string{"error":"program_id required"})
  }
  if err := h.Repo.Remove(c.Request().Context(), u.ID, id); err != nil {
    return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
  }
  return c.NoContent(http.StatusNoContent)
}
