package search

import "github.com/labstack/echo/v4"

func MapSearchRoutes(searchGroup *echo.Group, h *Handlers) {
	searchGroup.GET("/search/:term", h.Search())
}
