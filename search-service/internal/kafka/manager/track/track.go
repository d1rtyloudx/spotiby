package track

import (
	"context"
	"github.com/d1rtyloudx/spotiby/search-service/internal/model"
)

type indexer interface {
	IndexTrack(ctx context.Context, track model.Track) error
}
