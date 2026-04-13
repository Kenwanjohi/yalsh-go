package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Kenwanjohi/yalsh-go/internal/validator"
)

type Link struct {
	ID        int64      `json:"link_id"`
	URL       string     `json:"url"`
	Key       string     `json:"key"`
	Clicks    int        `json:"clicks"`
	UserID    int64      `json:"user_id"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt *string    `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"_"`
	Version   int        `json:"version"`
}

type LinkModel struct {
	DB *sql.DB
}

func ValidateLink(v *validator.Validator, link *Link) {
	v.Check(link.URL != "", "url", "must be provided")
	v.Check(v.IsURL(link.URL), "url", "must be a valid absolute URL")

	v.Check(link.Key != "", "key", "must be provided")
	v.CheckLength(link.Key, 3, 20, "key", "must be between 3 and 20 characters long")
	v.Check(v.Matches(link.Key, validator.KeyRX), "key", "must only contain alphanumeric characters and underscores")
}

func (m *LinkModel) Insert(context context.Context, link *Link) error {
	query := `
	INSERT INTO links (url, key, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	args := []interface{}{link.URL, link.Key, link.UserID}

	return m.DB.QueryRowContext(context, query, args...).Scan(&link.ID, &link.CreatedAt)
}

func (m *LinkModel) GetByKey(context context.Context, key string) (*Link, error) {
	query := `
	SELECT id, url, key, clicks, user_id, created_at
	FROM links
	WHERE key = $1 AND deleted_at IS NULL`

	var link Link
	err := m.DB.QueryRowContext(context, query, key).Scan(&link.ID, &link.URL, &link.Key, &link.Clicks, &link.UserID, &link.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (m *LinkModel) GetAllForUser(context context.Context, userID int64) ([]*Link, error) {
	query := `
	SELECT id, url, key, clicks, user_id, created_at
	FROM links
	WHERE user_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC`

	rows, err := m.DB.QueryContext(context, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var links []*Link

	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.URL, &link.Key, &link.Clicks, &link.UserID, &link.CreatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (m *LinkModel) Delete(context context.Context, id int64, userID int64) error {
	query := `
	UPDATE links
	SET deleted_at = NOW(), version = version + 1
	WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	result, err := m.DB.ExecContext(context, query, id, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found or unauthorized")
	}
	return nil
}

func (m *LinkModel) GetAndIncrementClicks(context context.Context, key string) (string, error) {
	query := `
	UPDATE links
	SET clicks = clicks + 1
	WHERE key = $1 AND deleted_at IS NULL
	RETURNING url`

	var url string

	err := m.DB.QueryRowContext(context, query, key).Scan(&url)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}

	return url, nil
}

func (m *LinkModel) Update(ctx context.Context, link *Link) error {
	query := `
	UPDATE links
	SET url = $1, key = $2, version = version + 1, updated_at = NOW()
	WHERE id = $3 AND user_id = $4 AND version = $5 AND deleted_at IS NULL
	RETURNING version`

	args := []interface{}{link.URL, link.Key, link.ID, link.UserID, link.Version}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&link.Version)
	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"links_unique_key_active\"":
			return ErrDuplicateKey
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *LinkModel) GetByID(ctx context.Context, id int64, userID int64) (*Link, error) {
	query := `
	SELECT id, url, key, clicks, user_id, version, created_at, updated_at
	FROM links
	WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var link Link
	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(&link.ID, &link.URL, &link.Key, &link.Clicks, &link.UserID, &link.Version, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &link, nil
}
