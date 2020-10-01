package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetDestinationByCityCountry(c echo.Context) error {
	city := c.Param("city")
	country := c.Param("country")
	destination := h.db.ByCityCountry(city, country)
	return c.JSON(http.StatusOK, destination)
}

func (h *Handler) GetDestinationByCountry(c echo.Context) error {
	country := c.Param("country")
	destinations := h.db.ByCountry(country)
	return c.JSON(http.StatusOK, destinations)
}

func (h *Handler) GetDestinations(c echo.Context) error {
	destinations := h.db.All()
	return c.JSON(http.StatusOK, destinations)
}
