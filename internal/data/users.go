package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  password   `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
	Version   int        `json:"version"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

type UserModel struct {
	DB *sql.DB
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (m *UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version`

	args := []interface{}{user.Username, user.Email, user.Password.hash}

	err := m.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at
	FROM users
	WHERE email = $1 AND deleted_at IS NULL`

	var user User
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password.hash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetByID(id int64) (*User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at
	FROM users
	WHERE id = $1`

	var user User

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(context, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password.hash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET username = $1, email = $2, password_hash = $3, version = version + 1, updated_at = NOW()
	WHERE id = $4 AND version = $5 AND deleted_at IS NULL`

	args := []interface{}{user.Username, user.Email, user.Password.hash, user.ID, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *UserModel) Delete(context context.Context, user *User) error {
	query := `
	UPDATE users
	SET deleted_at = NOW(), version = version + 1
	WHERE id = $1 AND version = $2 AND deleted_at IS NULL`

	result, err := m.DB.ExecContext(context, query, user.ID, user.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found or edit conflict")
	}

	return nil
}
