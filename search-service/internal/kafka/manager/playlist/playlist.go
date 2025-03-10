package playlist

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
	IndexPlaylist(ctx context.Context, playlist model.Playlist) error
	DeletePlaylist(ctx context.Context, deletePlaylist model.DeletePlaylist) error
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
			m.log.Warn(
				"error fetching message",
				zap.String("op", "manager.Playlist.ProcessMessages - r.FetchMessage"),
				zap.Int("worker_id", workerID),
				zap.Error(err),
			)
			continue
		}

		switch msg.Topic {
		case m.topics.CreatePlaylistTopic.TopicName:
			m.processIndexPlaylist(ctx, r, msg)
		case m.topics.UpdateProfileTopic.TopicName:
			m.processIndexPlaylist(ctx, r, msg)
		case m.topics.DeletePlaylistTopic.TopicName:
			m.processDeletePlaylist(ctx, r, msg)
		}
	}
}

func (m *ProcessManager) processIndexPlaylist(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	var playlist model.Playlist
	if err := json.Unmarshal(msg.Value, &playlist); err != nil {
		m.log.Warn("failed to unmarshal playlist msg")
		return
	}

	err := m.indexer.IndexPlaylist(ctx, playlist)
	if err != nil {
		return
	}

	if err := r.CommitMessages(ctx, msg); err != nil {
		m.log.Warn("failed to commit msgs", zap.Error(err))
	}

	m.log.Info("successfully process msg", zap.Any("data", playlist))
}

func (m *ProcessManager) processDeletePlaylist(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	var req model.DeletePlaylist
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		m.log.Warn("failed to unmarshal playlist msg")
		return
	}

	err := m.indexer.DeletePlaylist(ctx, req)
	if err != nil {
		return
	}

	if err := r.CommitMessages(ctx, msg); err != nil {
		m.log.Warn("failed to commit msgs", zap.Error(err))
	}

	m.log.Info("successfully process msg", zap.Any("data", req))
}
