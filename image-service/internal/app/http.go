package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"image-service/internal/http/image"
	"net/http"
)

func (a *App) runHTTPServer(imageHandlers *image.Handlers) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.HTTP.Port),
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
	}

	a.mapRoutes(imageHandlers)

	err := a.echo.StartServer(server)
	if err != nil {
		return err
	}

	a.log.Info("http server started", zap.Int("port", a.cfg.HTTP.Port))
	return nil
}

func (a *App) mapRoutes(imageHandlers *image.Handlers) {
	apiGroup := a.echo.Group("/api")

	v1 := apiGroup.Group("/v1")

	imageGroup := v1.Group("/image")

	image.MapImageRoutes(imageGroup, imageHandlers)
}

func (a *App) stopHTTPServer(ctx context.Context) error {
	a.log.Info("shutting down http server")

	if err := a.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
