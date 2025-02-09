package app

import (
	"context"
	"fmt"
	miniolib "github.com/d1rtyloudx/spotiby-pkg/minio"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	imagehand "github.com/d1rtyloudx/spotiby/user-service/internal/http/image"
	imagesvc "github.com/d1rtyloudx/spotiby/user-service/internal/service/image"
	miniostorage "github.com/d1rtyloudx/spotiby/user-service/internal/storage/minio"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	echo   *echo.Echo
	client *minio.Client
	log    *zap.Logger
	cfg    *config.Config
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	return &App{
		echo: echo.New(),
		log:  log,
		cfg:  cfg,
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.client = miniolib.MustConnect(&a.cfg.Minio.Connection)

	a.mustInitBuckets(ctx)

	imageStorage := miniostorage.NewImageStorage(
		a.client,
		fmt.Sprintf("%s:%d", a.cfg.Minio.Connection.Host, a.cfg.Minio.Connection.Port),
	)

	publisher := rabbitmq.MustCreatePublisher(&a.cfg.RabbitMQ.Connection)
	defer publisher.Close()

	imageService := imagesvc.New(imageStorage, publisher, a.log, &a.cfg.RabbitMQ)
	imageHandlers := imagehand.New(imageService, &a.cfg.Minio.Buckets, a.log)

	go func() {
		err := a.runHTTPServer(imageHandlers)
		if err != nil {
			a.log.Error("failed to start http server", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()

	err := a.stopHTTPServer(ctx)
	if err != nil {
		a.log.Warn("failed to stop http server", zap.Error(err))
	}

	return nil
}
