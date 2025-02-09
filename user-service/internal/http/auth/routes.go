package auth

import (
	"github.com/labstack/echo/v4"
	"user-service/internal/http/middleware"
)

func MapAuthRoutes(authGroup *echo.Group, h *Handlers, authMiddleware *middleware.AuthMiddleware) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.POST("/refresh", h.RefreshToken())

	authGroup.GET("/introspect", h.IntrospectToken(), authMiddleware.ParseAndVerifyAccessToken())

	credGroup := authGroup.Group("/credential")
	credGroup.Use(authMiddleware.ParseAndVerifyAccessToken())
	credGroup.PUT("/password", h.UpdatePassword())
	credGroup.PUT("/username", h.UpdateUsername())
}
