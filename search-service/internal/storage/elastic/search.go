package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/d1rtyloudx/spotiby-pkg/lib"
	"github.com/d1rtyloudx/spotiby/search-service/internal/config"
	"github.com/d1rtyloudx/spotiby/search-service/internal/model"
	"github.com/elastic/go-elasticsearch/v8"
)

var searchFields = []string{"description", "display_name", "name"}

type SearchStorage struct {
	client  *elasticsearch.Client
	indexes *config.ElasticIndexes
}

func NewSearchStorage(client *elasticsearch.Client, indexes *config.ElasticIndexes) *SearchStorage {
	return &SearchStorage{
		client:  client,
		indexes: indexes,
	}
}

func (s *SearchStorage) delete(ctx context.Context, index string, documentID string) error {
	const op = "storage.elastic.SearchStorage.delete"

	resp, err := s.client.Delete(
		index,
		documentID,
		s.client.Delete.WithPretty(),
		s.client.Delete.WithHuman(),
		s.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("%s - s.client.Delete: %w", op, err)
	}

	defer resp.Body.Close() //wrap

	if resp.IsError() {
		return fmt.Errorf("%s - resp.IsError: %s", op, resp.Status())
	}

	return nil
}

func (s *SearchStorage) index(ctx context.Context, index string, documentID string, data interface{}) error {
	const op = "storage.elastic.SearchStorage.index"
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	resp, err := s.client.Index(
		index,
		bytes.NewReader(dataBytes),
		s.client.Index.WithPretty(),
		s.client.Index.WithHuman(),
		s.client.Index.WithDocumentID(documentID),
		s.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("%s - s.client.Index: %w", op, err)
	}

	defer resp.Body.Close() //wrap

	if resp.IsError() {
		return fmt.Errorf("%s - resp.IsError: %s", op, resp.Status())
	}

	return nil
}

func (s *SearchStorage) IndexProfile(ctx context.Context, profile model.Profile) error {
	return s.index(ctx, s.indexes.ProfileIndex.Name, profile.ID, profile)
}

func (s *SearchStorage) IndexTrack(ctx context.Context, track model.Track) error {
	return s.index(ctx, s.indexes.TrackIndex.Name, track.ID, track)
}

func (s *SearchStorage) IndexPlaylist(ctx context.Context, playlist model.Playlist) error {
	return s.index(ctx, s.indexes.PlaylistIndex.Name, playlist.ID, playlist)
}

func (s *SearchStorage) DeleteProfile(ctx context.Context, id string) error {
	return s.delete(ctx, s.indexes.ProfileIndex.Name, id)
}

func (s *SearchStorage) DeletePlaylist(ctx context.Context, id string) error {
	return s.delete(ctx, s.indexes.PlaylistIndex.Name, id)
}

func (s *SearchStorage) DeleteTrack(ctx context.Context, id string) error {
	return s.delete(ctx, s.indexes.TrackIndex.Name, id)
}

func (s *SearchStorage) Search(ctx context.Context, term string, paginationQuery lib.PaginationQuery) (model.SearchResponse, error) {
	const op = "elastic.SearchStorage.Search"

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  term,
				"fields": searchFields,
			},
		},
	}

	searchQueryBytes, err := json.Marshal(&searchQuery)
	if err != nil {
		return model.SearchResponse{}, fmt.Errorf("%s - json.Marshal: %w", op, err)
	}

	resp, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.indexes.ProfileIndex.Name),
		s.client.Search.WithBody(bytes.NewReader(searchQueryBytes)),
		s.client.Search.WithPretty(),
		s.client.Search.WithHuman(),
		s.client.Search.WithSize(int(paginationQuery.Limit)),
		s.client.Search.WithFrom(int(paginationQuery.GetOffset())),
	)
	if err != nil {
		return model.SearchResponse{}, fmt.Errorf("%s - s.client.Search: %w", op, err)
	}
	defer resp.Body.Close() //wrap

	if resp.IsError() {
		return model.SearchResponse{}, fmt.Errorf("%s - resp.IsError: %s", op, resp.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.SearchResponse{}, fmt.Errorf("%s - json.NewDecoder: %w", op, err)
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	paginationResponse := lib.NewPaginationResponse(
		uint64(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		paginationQuery.Limit,
		paginationQuery.Page,
	)

	searchResponse := model.SearchResponse{
		Playlists:  make([]model.Playlist, 0),
		Tracks:     make([]model.Track, 0),
		Profiles:   make([]model.Profile, 0),
		Pagination: paginationResponse,
	}

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		index := hitMap["_index"].(string)

		sourceJSON, err := json.Marshal(source)
		if err != nil {
			continue
		}

		switch index {
		case s.indexes.ProfileIndex.Name:
			var profile model.Profile
			if err := json.Unmarshal(sourceJSON, &profile); err != nil {
				continue
			}
			searchResponse.Profiles = append(searchResponse.Profiles, profile)

		case s.indexes.TrackIndex.Name:
			var track model.Track
			if err := json.Unmarshal(sourceJSON, &track); err != nil {
				continue
			}
			searchResponse.Tracks = append(searchResponse.Tracks, track)

		case s.indexes.PlaylistIndex.Name:
			var playlist model.Playlist
			if err := json.Unmarshal(sourceJSON, &playlist); err != nil {
				continue
			}
			searchResponse.Playlists = append(searchResponse.Playlists, playlist)
		}
	}

	return searchResponse, nil
}
