package app

import (
	"context"
	"github.com/d1rtyloudx/spotiby-pkg/postgres"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/d1rtyloudx/spotiby-pkg/redis"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	authhand "github.com/d1rtyloudx/spotiby/user-service/internal/http/auth"
	"github.com/d1rtyloudx/spotiby/user-service/internal/http/middleware"
	profilehand "github.com/d1rtyloudx/spotiby/user-service/internal/http/profile"
	rabbitmqconsumer "github.com/d1rtyloudx/spotiby/user-service/internal/rabbitmq"
	authsvc "github.com/d1rtyloudx/spotiby/user-service/internal/service/auth"
	profilesvc "github.com/d1rtyloudx/spotiby/user-service/internal/service/profile"
	postgresstorage "github.com/d1rtyloudx/spotiby/user-service/internal/storage/postgres"
	rediscache "github.com/d1rtyloudx/spotiby/user-service/internal/storage/redis"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	echo *echo.Echo
	log  *zap.Logger
	cfg  *config.Config
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

	db := postgres.MustConnect(&a.cfg.Postgres)
	defer db.Close() //wrap

	client := redis.MustConnect(&a.cfg.Redis)
	defer client.Close() // wrap

	credentialStorage := postgresstorage.NewCredentialStorage(db)
	profileStorage := postgresstorage.NewProfileStorage(db)
	tokenBlacklist := rediscache.NewTokenBlacklist(client)

	authService := authsvc.New(credentialStorage, profileStorage, tokenBlacklist, a.log, &a.cfg.Token)
	profileService := profilesvc.New(profileStorage, a.log)

	profileHandlers := profilehand.New(profileService)
	authHandlers := authhand.New(authService)

	authMiddleware := middleware.NewAuthMiddleware(authService, a.log)

	imageConsumer, err := rabbitmqconsumer.NewConsumer(profileService, &a.cfg.RabbitMQ, a.log)
	if err != nil {
		return err
	}
	defer imageConsumer.Close() //wrap

	go func() {
		err = rabbitmq.ConsumeQueue(
			ctx,
			imageConsumer.AmqpChan,
			a.cfg.RabbitMQ.ImageBinding.Concurrency,
			a.cfg.RabbitMQ.ImageBinding.QueueName,
			a.cfg.RabbitMQ.ImageBinding.ConsumerTag,
			imageConsumer.UploadProfileAvatar,
		)
		if err != nil {
			a.log.Error("failed to consume query", zap.Error(err))
			stop()
		}
	}()

	go func() {
		err = a.runHTTPServer(profileHandlers, authHandlers, authMiddleware)
		if err != nil {
			a.log.Error("failed to run http server", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()

	err = a.stopHTTPServer(ctx)
	if err != nil {
		a.log.Warn("failed to stop http server", zap.Error(err))
	}

	return nil
}
