package app

import (
	"context"
	"github.com/d1rtyloudx/spotiby-pkg/elastic"
	kafkapkg "github.com/d1rtyloudx/spotiby-pkg/kafka"
	"github.com/d1rtyloudx/spotiby/search-service/internal/config"
	searchhand "github.com/d1rtyloudx/spotiby/search-service/internal/http/search"
	playlistmanager "github.com/d1rtyloudx/spotiby/search-service/internal/kafka/manager/playlist"
	profilemanager "github.com/d1rtyloudx/spotiby/search-service/internal/kafka/manager/profile"
	trackmanager "github.com/d1rtyloudx/spotiby/search-service/internal/kafka/manager/track"
	searchsvc "github.com/d1rtyloudx/spotiby/search-service/internal/service/search"
	elasticstorage "github.com/d1rtyloudx/spotiby/search-service/internal/storage/elastic"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

const poolSize = 10

type App struct {
	echo      *echo.Echo
	kafkaConn *kafka.Conn
	log       *zap.Logger
	cfg       *config.Config
}

func NewApp(cfg *config.Config, log *zap.Logger) *App {
	return &App{
		echo: echo.New(),
		log:  log,
		cfg:  cfg,
	}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := elastic.MustConnect(&a.cfg.Elastic.Connection)

	searchStorage := elasticstorage.NewSearchStorage(client, &a.cfg.Elastic.Indexes)
	searchService := searchsvc.New(searchStorage, a.log)

	searchHandlers := searchhand.New(searchService)

	err := a.connectKafkaTopics(ctx)
	if err != nil {
		return err
	}
	defer a.kafkaConn.Close() //wrap

	//TODO add GROUP ids to cfg
	profileCg := kafkapkg.NewConsumerGroup(a.cfg.Kafka.Connection.Brokers, a.cfg.Kafka.Connection.ConsumerGroups.ProfileName, a.log)
	trackCg := kafkapkg.NewConsumerGroup(a.cfg.Kafka.Connection.Brokers, a.cfg.Kafka.Connection.ConsumerGroups.TrackName, a.log)
	playlistCg := kafkapkg.NewConsumerGroup(a.cfg.Kafka.Connection.Brokers, a.cfg.Kafka.Connection.ConsumerGroups.PlaylistName, a.log)

	profileProcessManager := profilemanager.NewProcessManager(searchService, &a.cfg.Kafka.Topics, a.log)
	trackProcessManager := trackmanager.NewProcessManager(searchService, &a.cfg.Kafka.Topics, a.log)
	playlistProcessManager := playlistmanager.NewProcessManager(searchService, &a.cfg.Kafka.Topics, a.log)

	go func() {
		if err := a.runHTTPServer(searchHandlers); err != nil {
			a.log.Error("failed to run http server", zap.Error(err))
			stop()
		}
	}()

	go profileCg.ConsumeTopic(ctx, a.getTopicGroup(), poolSize, profileProcessManager.ProcessMessages)

	go trackCg.ConsumeTopic(ctx, a.getTopicGroup(), poolSize, trackProcessManager.ProcessMessages)

	go playlistCg.ConsumeTopic(ctx, a.getTopicGroup(), poolSize, playlistProcessManager.ProcessMessages)

	<-ctx.Done()

	err = a.stopHTTPServer(ctx)
	if err != nil {
		a.log.Warn("failed to stop http server", zap.Error(err))
	}

	return nil
}
