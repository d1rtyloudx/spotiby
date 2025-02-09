package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"user-service/internal/domain/model"
	"user-service/internal/storage"
)

type CredentialStorage struct {
	db *sqlx.DB
}

func NewCredentialStorage(db *sqlx.DB) *CredentialStorage {
	return &CredentialStorage{
		db: db,
	}
}

func (s *CredentialStorage) Create(ctx context.Context, cred model.Credential) (string, error) {
	const op = "postgres.CredentialStorage.Create"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("%s - s.db.BeginTx: %w", op, err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	query, args, err := squirrel.
		Insert("credentials").
		Columns("username, email, hash_pass").
		Values(cred.Username, cred.Email, cred.HashPass).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s - squirrel.Insert: %w", op, err)
	}

	var credentialID string
	row := tx.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&credentialID); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == storage.UniqueViolationCode {
				return "", storage.ErrAlreadyExists
			}
		}
		return "", fmt.Errorf("%s - tx.QueryRowContext: %w", op, err)
	}

	query, args, err = squirrel.
		Insert("profiles").
		Columns("credential_id").
		Values(credentialID).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s - squirrel.Insert: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("%s - tx.ExecContext: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("%s - tx.Commit: %w", op, err)
	}

	return credentialID, nil
}

func (s *CredentialStorage) Update(ctx context.Context, cred model.Credential) error {
	const op = "postgres.CredentialStorage.Update"

	builder := squirrel.
		Update("credentials").
		Where(squirrel.Eq{"id": cred.ID}).
		PlaceholderFormat(squirrel.Dollar)

	if cred.HashPass != nil {
		builder.Set("pass_hash", cred.HashPass)
	}

	if cred.Username != "" {
		builder.Set("username", cred.Username)
	}

	if cred.Email != "" {
		builder.Set("email", cred.Email)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s - builder.ToSql: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == storage.UniqueViolationCode {
				return storage.ErrAlreadyExists
			}
		}
		return fmt.Errorf("%s - s.db.ExecContext: %w", op, err)
	}

	return nil
}

func (s *CredentialStorage) Delete(ctx context.Context, id string) error {
	const op = "postgres.CredentialStorage.Delete"

	query, args, err := squirrel.
		Delete("credentials").
		Where(squirrel.Eq{"id": id}).
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

func (s *CredentialStorage) GetByID(ctx context.Context, id string) (model.Credential, error) {
	return s.getByField(ctx, "id", id)
}

func (s *CredentialStorage) GetByUsername(ctx context.Context, username string) (model.Credential, error) {
	return s.getByField(ctx, "username", username)
}

func (s *CredentialStorage) GetByEmail(ctx context.Context, email string) (model.Credential, error) {
	return s.getByField(ctx, "email", email)
}

func (s *CredentialStorage) getByField(ctx context.Context, field string, value interface{}) (model.Credential, error) {
	const op = "postgres.CredentialStorage.getByField"

	query, args, err := squirrel.
		Select(`id, username, email, hash_pass, role, created_at, updated_at`).
		From("credentials").
		Where(squirrel.Eq{field: value}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return model.Credential{}, fmt.Errorf("%s - squirrel.Select: %w", op, err)
	}

	var creds model.Credential
	row := s.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(
		&creds.ID,
		&creds.Username,
		&creds.Email,
		&creds.HashPass,
		&creds.Role,
		&creds.CreatedAt,
		&creds.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Credential{}, storage.ErrNotFound
		}
		return model.Credential{}, fmt.Errorf("%s - row.Scan: %w", op, err)
	}

	return creds, nil
}
