package repository

import (
	"database/sql"

	"ozinse/internal/model"
)

type FavoriteRepo struct {
	db *sql.DB
}

func NewFavoriteRepo(db *sql.DB) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (r *FavoriteRepo) GetByUserID(userID int) ([]model.Favorite, error) {
	rows, err := r.db.Query(`
		SELECT f.user_id, f.project_id, f.added_at,
			   p.id, p.title, p.cover_image_url, p.project_type, p.release_year, p.duration_minutes
		FROM favorites f
		JOIN projects p ON f.project_id = p.id
		WHERE f.user_id = $1
		ORDER BY f.added_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Favorite
	for rows.Next() {
		var f model.Favorite
		var p model.Project
		err := rows.Scan(&f.UserID, &f.ProjectID, &f.AddedAt,
			&p.ID, &p.Title, &p.CoverImageURL, &p.ProjectType, &p.ReleaseYear, &p.DurationMinutes)
		if err != nil {
			return nil, err
		}
		f.Project = &p
		list = append(list, f)
	}
	return list, nil
}

func (r *FavoriteRepo) Add(userID, projectID int) error {
	_, err := r.db.Exec(
		`INSERT INTO favorites (user_id, project_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, projectID,
	)
	return err
}

func (r *FavoriteRepo) Remove(userID, projectID int) error {
	_, err := r.db.Exec(
		`DELETE FROM favorites WHERE user_id = $1 AND project_id = $2`,
		userID, projectID,
	)
	return err
}