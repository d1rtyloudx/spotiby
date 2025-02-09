package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/d1rtyloudx/spotiby-pkg/lib"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/d1rtyloudx/spotiby/user-service/internal/storage"
	"github.com/jmoiron/sqlx"
)

type ProfileStorage struct {
	db *sqlx.DB
}

func NewProfileStorage(db *sqlx.DB) *ProfileStorage {
	return &ProfileStorage{
		db: db,
	}
}

func (s *ProfileStorage) get(ctx context.Context, builder squirrel.SelectBuilder, pageQuery lib.PaginationQuery) ([]model.Profile, lib.PaginationResponse, error) {
	const op = "postgres.ProfileStorage.get"

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, lib.PaginationResponse{}, fmt.Errorf("%s - builder.ToSql: %w", op, err)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS subquery", query)

	row := s.db.QueryRowContext(ctx, countQuery, args...)

	var totalCount uint64
	if err := row.Scan(&totalCount); err != nil {
		return nil, lib.PaginationResponse{}, fmt.Errorf("%s - rows.Scan: %w", op, err)
	}

	if totalCount == 0 {
		return []model.Profile{}, lib.NewPaginationResponse(0, pageQuery.Limit, pageQuery.Page), nil
	}

	offset := pageQuery.GetOffset()

	queryWithPaging, args, err := builder.
		Offset(offset).
		Limit(pageQuery.Limit).
		ToSql()
	if err != nil {
		return nil, lib.PaginationResponse{}, fmt.Errorf("%s - builder.ToSql: %w", op, err)
	}

	rows, err := s.db.QueryContext(ctx, queryWithPaging, args...)
	if err != nil {
		return nil, lib.PaginationResponse{}, fmt.Errorf("%s - s.db.QueryContext: %w", op, err)
	}
	defer rows.Close() //wrap

	var profiles []model.Profile
	for rows.Next() {
		var profile model.Profile
		if err := rows.Scan(
			&profile.ID,
			&profile.DisplayName,
			&profile.FirstName,
			&profile.LastName,
			&profile.Description,
			&profile.CredentialID,
			&profile.AvatarURL,
		); err != nil {
			return nil, lib.PaginationResponse{}, fmt.Errorf("%s - rows.Scan: %w", op, err)
		}

		profiles = append(profiles, profile)
	}

	pageResp := lib.NewPaginationResponse(totalCount, pageQuery.Limit, pageQuery.Page)

	return profiles, pageResp, nil
}

func (s *ProfileStorage) Get(ctx context.Context, pageQuery lib.PaginationQuery) ([]model.Profile, lib.PaginationResponse, error) {
	builder := squirrel.
		Select("*").
		From("profiles")

	profiles, pageResp, err := s.get(ctx, builder, pageQuery)
	if err != nil {
		return nil, lib.PaginationResponse{}, err
	}

	return profiles, pageResp, nil
}

func (s *ProfileStorage) GetFollows(ctx context.Context, profileID string, pageQuery lib.PaginationQuery) ([]model.Profile, lib.PaginationResponse, error) {
	builder := squirrel.
		Select("p.*").
		From("followers AS f").
		InnerJoin("profiles AS p ON f.following_id = p.id").
		Where(squirrel.Eq{"f.follower_id": profileID}).
		PlaceholderFormat(squirrel.Dollar)

	profiles, pageResp, err := s.get(ctx, builder, pageQuery)
	if err != nil {
		return nil, lib.PaginationResponse{}, err
	}

	return profiles, pageResp, nil
}

func (s *ProfileStorage) GetByID(ctx context.Context, id string) (model.Profile, error) {
	return s.getByField(ctx, "id", id)
}

func (s *ProfileStorage) GetByCredentialID(ctx context.Context, id string) (model.Profile, error) {
	return s.getByField(ctx, "credential_id", id)
}

func (s *ProfileStorage) getByField(ctx context.Context, field string, value interface{}) (model.Profile, error) {
	const op = "postgres.ProfileStorage.getByField"

	query, args, err := squirrel.
		Select("id, display_name, first_name, last_name, description, avatar_url, credential_id").
		From("profiles").
		Where(squirrel.Eq{field: value}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Profile{}, fmt.Errorf("%s - squirrel.Select: %w", op, err)
	}

	var profile model.Profile
	row := s.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(
		&profile.ID,
		&profile.DisplayName,
		&profile.FirstName,
		&profile.LastName,
		&profile.Description,
		&profile.AvatarURL,
		&profile.CredentialID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Profile{}, storage.ErrNotFound
		}
		return model.Profile{}, fmt.Errorf("%s - s.db.QueryRowContext: %w", op, err)
	}

	return profile, nil
}

func (s *ProfileStorage) Update(ctx context.Context, profile model.Profile) error {
	const op = "postgres.ProfileStorage.Update"

	builder := squirrel.
		Update("profiles").
		Where(squirrel.Eq{"id": profile.ID}).
		PlaceholderFormat(squirrel.Dollar)

	if profile.FirstName != "" {
		builder = builder.Set("first_name", profile.FirstName)
	}

	if profile.LastName != "" {
		builder = builder.Set("last_name", profile.LastName)
	}

	if profile.Description != "" {
		builder = builder.Set("description", profile.Description)
	}

	if profile.DisplayName != "" {
		builder = builder.Set("display_name", profile.DisplayName)
	}

	if profile.AvatarURL != "" {
		builder = builder.Set("avatar_url", profile.AvatarURL)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s - builder.ToSql: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s - s.db.ExecContext: %w", op, err)
	}

	return nil
}

func (s *ProfileStorage) FollowProfile(ctx context.Context, followerID string, followeeID string) error {
	const op = "postgres.ProfileStorage.FollowProfile"

	query, args, err := squirrel.
		Insert("followers").
		Columns("follower_id", "following_id").
		Values(followerID, followeeID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s - squirrel.Insert: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s - s.db.ExecContext: %w", op, err)
	}

	return nil
}

func (s *ProfileStorage) UnfollowProfile(ctx context.Context, followerID string, followeeID string) error {
	const op = "postgres.ProfileStorage.UnfollowProfile"

	query, args, err := squirrel.
		Delete("followers").
		Where(squirrel.Eq{"follower_id": followerID, "following_id": followeeID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s - squirrel.Delete: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s - s.db.ExecContext: %w", op, err)
	}

	return nil
}

func (s *ProfileStorage) GetFollowedProfiles(ctx context.Context, id string) ([]model.Profile, error) {
	const op = "postgres.ProfileStorage.GetFollowedProfiles"

	return nil, fmt.Errorf("not implemented")
}
