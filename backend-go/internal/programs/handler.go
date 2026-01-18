package programs

import (
  "net/http"
  "strconv"
  "strings"

  "github.com/labstack/echo/v4"
)

type Handler struct { Repo Repo }

func splitCSV(s string) []string {
  if strings.TrimSpace(s) == "" { return nil }
  parts := strings.Split(s, ",")
  out := make([]string, 0, len(parts))
  for _, p := range parts {
    p = strings.TrimSpace(p)
    if p != "" { out = append(out, p) }
  }
  return out
}

func (h Handler) List(c echo.Context) error {
  page, _ := strconv.Atoi(c.QueryParam("page"))
  limit, _ := strconv.Atoi(c.QueryParam("limit"))

  var minT *float64
  if v := c.QueryParam("min_tuition"); v != "" {
    f, _ := strconv.ParseFloat(v, 64); minT = &f
  }
  var maxT *float64
  if v := c.QueryParam("max_tuition"); v != "" {
    f, _ := strconv.ParseFloat(v, 64); maxT = &f
  }
  var sch *bool
  if v := c.QueryParam("scholarship"); v != "" {
    b := (v == "true"); sch = &b
  }

  params := ListParams{
    Q: c.QueryParam("q"),
    Countries: splitCSV(c.QueryParam("countries")),
    Levels: splitCSV(c.QueryParam("levels")),
    Fields: splitCSV(c.QueryParam("fields")),
    Currency: strings.TrimSpace(c.QueryParam("currency")),
    MinTuition: minT,
    MaxTuition: maxT,
    Scholarship: sch,
    Sort: c.QueryParam("sort"),
    Page: page,
    Limit: limit,
  }

  items, total, err := h.Repo.List(c.Request().Context(), params)
  if err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()}) }

  return c.JSON(http.StatusOK, map[string]any{
    "page": params.Page,
    "limit": params.Limit,
    "total": total,
    "items": items,
  })
}
