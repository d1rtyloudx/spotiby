package track

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
	IndexTrack(ctx context.Context, track model.Track) error
	DeleteTrack(ctx context.Context, req model.DeleteTrack) error
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
			m.log.Warn("failed to fetch message",
				zap.String("op", "manager.track.ProcessMessages - r.FetchMessage"),
				zap.Int("worker_id", workerID),
				zap.Error(err),
			)
			continue
		}

		switch msg.Topic {
		case m.topics.CreateTrackTopic.TopicName:
			m.processIndexTrack(ctx, r, msg)
		case m.topics.DeleteTrackTopic.TopicName:
			m.processDeleteTrack(ctx, r, msg)
		case m.topics.UpdateTrackTopic.TopicName:
			m.processIndexTrack(ctx, r, msg)
		}
	}
}

func (m *ProcessManager) processIndexTrack(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	var track model.Track
	if err := json.Unmarshal(msg.Value, &track); err != nil {
		m.log.Warn("failed to unmarshal track msg", zap.Error(err))
		return
	}

	err := m.indexer.IndexTrack(ctx, track)
	if err != nil {
		return
	}

	if err := r.CommitMessages(ctx, msg); err != nil {
		m.log.Warn("failed to commit msgs", zap.Error(err))
	}

	m.log.Info("successfully process msg", zap.Any("data", track))
}

func (m *ProcessManager) processDeleteTrack(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	var req model.DeleteTrack
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		m.log.Warn("failed to unmarshal track msg", zap.Error(err))
		return
	}

	err := m.indexer.DeleteTrack(ctx, req)
	if err != nil {
		return
	}

	if err := r.CommitMessages(ctx, msg); err != nil {
		m.log.Warn("failed to commit msgs", zap.Error(err))
	}

	m.log.Info("successfully process msg", zap.Any("data", req))
}
