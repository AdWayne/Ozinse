package repository

import (
	"database/sql"
	"time"

	"ozinse/internal/model"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) Save(userID int, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *TokenRepo) FindByToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.QueryRow(
		`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = $1`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *TokenRepo) DeleteByUserID(userID int) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

func (r *TokenRepo) DeleteByToken(token string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}