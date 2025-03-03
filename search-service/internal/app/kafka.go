package app

import (
	"context"
	kafkapkg "github.com/d1rtyloudx/spotiby-pkg/kafka"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func (a *App) connectKafkaTopics(ctx context.Context) error {
	conn, err := kafkapkg.NewKafkaConn(a.cfg.Kafka.Connection.Brokers[0])
	if err != nil {
		return err
	}

	err = a.initKafkaTopics(ctx, conn)
	if err != nil {
		return err
	}

	brokers, err := conn.Brokers()
	if err != nil {
		return err
	}

	a.log.Info(
		"kafka successfully connected to brokers",
		zap.Any("brokers", brokers),
	)

	a.kafkaConn = conn

	return nil
}

func (a *App) getTopicGroup() []string {
	return []string{
		a.cfg.Kafka.Topics.CreateProfileTopic.TopicName,
		a.cfg.Kafka.Topics.UpdateProfileTopic.TopicName,
		a.cfg.Kafka.Topics.DeleteProfileTopic.TopicName,
	}
}

func (a *App) initKafkaTopics(ctx context.Context, conn *kafka.Conn) error {
	profileCreateTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.Topics.CreateProfileTopic.TopicName,
		NumPartitions:     a.cfg.Kafka.Topics.CreateProfileTopic.Partitions,
		ReplicationFactor: a.cfg.Kafka.Topics.CreateProfileTopic.ReplicationFactor,
	}

	profileUpdateTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.Topics.UpdateProfileTopic.TopicName,
		NumPartitions:     a.cfg.Kafka.Topics.UpdateProfileTopic.Partitions,
		ReplicationFactor: a.cfg.Kafka.Topics.UpdateProfileTopic.ReplicationFactor,
	}

	profileDeleteTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.Topics.DeleteProfileTopic.TopicName,
		NumPartitions:     a.cfg.Kafka.Topics.DeleteProfileTopic.Partitions,
		ReplicationFactor: a.cfg.Kafka.Topics.DeleteProfileTopic.ReplicationFactor,
	}

	err := conn.CreateTopics(
		profileCreateTopic,
		profileUpdateTopic,
		profileDeleteTopic,
	)
	if err != nil {
		a.log.Error("app.kafka.initKafkaTopics - conn.CreateTopics", zap.Error(err))
	}

	a.log.Info("successfully connected kafka topics")

	return nil
}
