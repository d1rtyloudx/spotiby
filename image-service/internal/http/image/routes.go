package image

import "github.com/labstack/echo/v4"

func MapImageRoutes(imageGroup *echo.Group, h *Handlers) {
	imageGroup.POST("/profile/:id", h.UploadProfile())
	imageGroup.POST("/track/:id", h.UploadTrack())
	imageGroup.POST("/playlist/:id", h.UploadPlaylist())

	imageGroup.GET("/profile/:fileName", h.GetProfile())
	imageGroup.GET("/track/:fileName", h.GetTrack())
	imageGroup.GET("/playlist/:fileName", h.GetPlaylist())
}
