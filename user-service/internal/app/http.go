package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"user-service/internal/http/auth"
	"user-service/internal/http/middleware"
	"user-service/internal/http/profile"
)

func (a *App) runHTTPServer(
	profileHandlers *profile.Handlers,
	authHandlers *auth.Handlers,
	authMiddleware *middleware.AuthMiddleware,
) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.HTTP.Port),
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
	}

	a.mapRoutes(profileHandlers, authHandlers, authMiddleware)

	err := a.echo.StartServer(server)
	if err != nil {
		return err
	}

	a.log.Info("http server started", zap.Int("port", a.cfg.HTTP.Port))
	return nil
}

func (a *App) mapRoutes(
	profileHandlers *profile.Handlers,
	authHandlers *auth.Handlers,
	authMiddleware *middleware.AuthMiddleware,
) {
	apiGroup := a.echo.Group("/api")

	v1 := apiGroup.Group("/v1")

	authGroup := v1.Group("/auth")

	profileGroup := v1.Group("/profile")

	profile.MapProfileRoutes(profileGroup, profileHandlers, authMiddleware)
	auth.MapAuthRoutes(authGroup, authHandlers, authMiddleware)
}

func (a *App) stopHTTPServer(ctx context.Context) error {
	a.log.Info("shutting down http server")

	if err := a.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
