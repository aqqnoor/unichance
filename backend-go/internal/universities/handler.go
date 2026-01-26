package universities

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Repo Repo
}

func (h Handler) GetByID(c echo.Context) error {
	id := c.Param("id")
	u, err := h.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "university not found"})
	}
	return c.JSON(http.StatusOK, u)
}
