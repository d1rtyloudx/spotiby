package profile

import (
	"context"
	"encoding/json"
	"github.com/d1rtyloudx/spotiby/search-service/internal/config"
	"github.com/d1rtyloudx/spotiby/search-service/internal/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sync"
)

type indexer interface {
	IndexProfile(ctx context.Context, profile model.Profile) error
	DeleteProfile(ctx context.Context, deleteProfile model.DeleteProfile) error
}

type ProcessManager struct {
	indexer indexer
	topics  *config.KafkaTopics
	log     *zap.Logger
}

func NewProcessManager(indexer indexer, topics *config.KafkaTopics, log *zap.Logger) *ProcessManager {
	return &ProcessManager{
		indexer: indexer,
		topics:  topics,
		log:     log,
	}
}

func (m *ProcessManager) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := r.FetchMessage(ctx)
		if err != nil {
			m.log.Warn("error fetching message",
				zap.String("op", "manager.Profile.ProcessMessages - r.FetchMessage"),
				zap.Int("worker_id", workerID),
				zap.Error(err),
			)
			continue
		}

		switch msg.Topic {
		case m.topics.CreateProfileTopic.TopicName:
			m.processIndexProfile(ctx, r, msg)
		case m.topics.UpdateProfileTopic.TopicName:
			m.processIndexProfile(ctx, r, msg)
		case m.topics.DeleteProfileTopic.TopicName:
			m.processDeleteProfile(ctx, r, msg)
		}
	}
}

func (m *ProcessManager) processIndexProfile(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	var profile model.Profile
	if err := json.Unmarshal(msg.Value, &profile); err != nil {
		m.log.Warn("failed to unmarshal profile message", zap.Error(err))
		return
	}

	err := m.indexer.IndexProfile(ctx, profile)
	if err != nil {
		return
	}

	if err := r.CommitMessages(ctx, msg); err != nil {
		m.log.Warn("failed to commit messages", zap.Error(err))
	}

	m.log.Info("successfully process msg", zap.Any("data", profile))
}

func (m *ProcessManager) processDeleteProfile(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	var req model.DeleteProfile
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		m.log.Warn("failed to unmarshal profile msg", zap.Error(err))
		return
	}

	err := m.indexer.DeleteProfile(ctx, req)
	if err != nil {
		return
	}

	if err := r.CommitMessages(ctx, msg); err != nil {
		m.log.Warn("failed to commit msgs", zap.Error(err))
	}

	m.log.Info("successfully process msg", zap.Any("data", req))
}
