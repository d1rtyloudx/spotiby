package image

import "github.com/labstack/echo/v4"

func MapImageRoutes(imageGroup *echo.Group, h *Handlers) {
	uploadGroup := imageGroup.Group("/upload")

	uploadGroup.POST("/profile/:id", h.UploadProfile())
	uploadGroup.POST("/track/:id", h.UploadTrack())
	uploadGroup.POST("/playlist/:id", h.UploadPlaylist())
}
