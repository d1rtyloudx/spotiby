package app

import (
	"context"
	"fmt"
	"github.com/d1rtyloudx/spotiby/search-service/internal/http/search"
	"go.uber.org/zap"
	"net/http"
)

func (a *App) runHTTPServer(searchHandlers *search.Handlers) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.HTTP.Port),
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
	}

	a.mapRoutes(searchHandlers)

	err := a.echo.StartServer(server)
	if err != nil {
		return err
	}

	a.log.Info("http server started", zap.Int("port", a.cfg.HTTP.Port))
	return nil
}

func (a *App) mapRoutes(searchHandlers *search.Handlers) {
	apiGroup := a.echo.Group("/api")

	v1 := apiGroup.Group("/v1")

	searchGroup := v1.Group("/search")

	search.MapSearchRoutes(searchGroup, searchHandlers)
}

func (a *App) stopHTTPServer(ctx context.Context) error {
	a.log.Info("shutting down http server")

	if err := a.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
