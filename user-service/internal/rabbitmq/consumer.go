package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	"github.com/d1rtyloudx/spotiby/user-service/internal/dto"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type profileUpdater interface {
	Update(ctx context.Context, id string, req dto.UpdateProfileRequest) error
}

type Consumer struct {
	updater  profileUpdater
	AmqpConn *amqp091.Connection
	AmqpChan *amqp091.Channel
	log      *zap.Logger
	cfg      *config.RabbitMQConfig
}

func NewConsumer(updater profileUpdater, cfg *config.RabbitMQConfig, log *zap.Logger) (*Consumer, error) {
	amqpConn, err := rabbitmq.NewRabbitMQConn(&cfg.Connection)
	if err != nil {
		return nil, err
	}

	amqpChan, err := amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = rabbitmq.DeclareQueue(amqpChan, cfg.ImageBinding.QueueName)
	if err != nil {
		return nil, err
	}

	err = rabbitmq.DeclareExchange(amqpChan, cfg.ImageBinding.ExchangeName, cfg.ImageBinding.ExchangeKind)
	if err != nil {
		return nil, err
	}

	err = rabbitmq.BindExchangeAndQueue(
		amqpChan,
		cfg.ImageBinding.ExchangeName,
		cfg.ImageBinding.QueueName,
		cfg.ImageBinding.RoutingKey,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		updater:  updater,
		AmqpConn: amqpConn,
		AmqpChan: amqpChan,
		log:      log,
		cfg:      cfg,
	}, nil
}

func (c *Consumer) UploadProfileAvatar(ctx context.Context, deliveries <-chan amqp091.Delivery, workerID int) func() error {
	return func() error {
		c.log.Info(
			"starting consumer",
			zap.Int("workerID", workerID),
			zap.String("queue", c.cfg.ImageBinding.QueueName),
		)

		for {
			select {
			case <-ctx.Done():
				c.log.Error("image consumer ctx done", zap.Error(ctx.Err()))
				return ctx.Err()

			case msg, ok := <-deliveries:
				if !ok {
					return fmt.Errorf("deliveries channel closed")
				}

				var req dto.UpdateAvatarProfileRequest
				err := json.Unmarshal(msg.Body, &req)
				if err != nil {
					return fmt.Errorf("failed to unmarshal update avatar profile: %w", err)
				}

				c.log.Info(
					"consume delivery",
					zap.Int("workerID", workerID),
					zap.Any("message", req),
				)

				err = c.updater.Update(ctx, req.ID, dto.UpdateProfileRequest{
					AvatarURL: req.AvatarURL,
				})
				if err != nil {
					continue
				}

				c.log.Info("successfully updated avatar profile")
			}
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.AmqpConn.Close(); err != nil {
		c.log.Warn("failed to close connection", zap.Error(err))
		return err
	}

	if err := c.AmqpChan.Close(); err != nil {
		c.log.Warn("failed to close channel", zap.Error(err))
		return err
	}

	return nil
}
