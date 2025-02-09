package profile

import (
	"github.com/labstack/echo/v4"
	"user-service/internal/http/middleware"
)

func MapProfileRoutes(profileGroup *echo.Group, h *Handlers, authMiddleware *middleware.AuthMiddleware) {
	profileGroup.Use(authMiddleware.ParseAndVerifyAccessToken())

	profileGroup.GET("/:id", h.GetByID())
	profileGroup.GET("/", h.Get())

	meGroup := profileGroup.Group("/me")
	meGroup.GET("", h.GetMe())
	meGroup.GET("/follows/", h.GetFollows())
	meGroup.PUT("", h.Update())
	meGroup.PUT("/follow/:id", h.FollowProfile())
	meGroup.PUT("/unfollow/:id", h.UnfollowProfile())
}
