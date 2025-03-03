package playlist

import (
	"context"
	"github.com/d1rtyloudx/spotiby/search-service/internal/model"
)

type indexer interface {
	IndexPlaylist(ctx context.Context, playlist model.Playlist) error
}
