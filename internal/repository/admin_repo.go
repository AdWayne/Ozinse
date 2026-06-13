package repository

import (
	"database/sql"
	"ozinse/internal/model"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

// PROJECTS

func (r *AdminRepo) CreateProject(req model.CreateProjectRequest) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(`
		INSERT INTO projects (title, description, release_year, director, producer,
			cover_image_url, banner_image_url, keywords, category_id, age_rating_id,
			project_type, duration_minutes, youtube_video_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, title, description, release_year, director, producer,
			cover_image_url, banner_image_url, keywords, category_id, age_rating_id,
			project_type, duration_minutes, youtube_video_id, created_at, updated_at
	`, req.Title, nn(req.Description), req.ReleaseYear, nn(req.Director), nn(req.Producer),
		nn(req.CoverImageURL), nn(req.BannerImageURL), nn(req.Keywords),
		req.CategoryID, req.AgeRatingID, req.ProjectType, req.DurationMinutes, nn(req.YouTubeVideoID),
	).Scan(&p.ID, &p.Title, &p.Description, &p.ReleaseYear, &p.Director,
		&p.Producer, &p.CoverImageURL, &p.BannerImageURL, &p.Keywords,
		&p.CategoryID, &p.AgeRatingID, &p.ProjectType, &p.DurationMinutes,
		&p.YouTubeVideoID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// жанры
	if len(req.GenreIDs) > 0 {
		for _, gID := range req.GenreIDs {
			r.db.Exec(`INSERT INTO project_genres (project_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, p.ID, gID)
		}
	}

	return &p, nil
}

func (r *AdminRepo) UpdateProject(id int, req model.CreateProjectRequest) error {
	_, err := r.db.Exec(`
		UPDATE projects SET title=$1, description=$2, release_year=$3, director=$4, producer=$5,
			cover_image_url=$6, banner_image_url=$7, keywords=$8, category_id=$9, age_rating_id=$10,
			project_type=$11, duration_minutes=$12, youtube_video_id=$13, updated_at=NOW()
		WHERE id=$14
	`, req.Title, nn(req.Description), req.ReleaseYear, nn(req.Director), nn(req.Producer),
		nn(req.CoverImageURL), nn(req.BannerImageURL), nn(req.Keywords),
		req.CategoryID, req.AgeRatingID, req.ProjectType, req.DurationMinutes, nn(req.YouTubeVideoID),
		id)
	if err != nil {
		return err
	}

	// обновляем жанры: удаляем старые, вставляем новые
	r.db.Exec(`DELETE FROM project_genres WHERE project_id = $1`, id)
	for _, gID := range req.GenreIDs {
		r.db.Exec(`INSERT INTO project_genres (project_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, gID)
	}
	return nil
}

func (r *AdminRepo) DeleteProject(id int) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = $1`, id)
	return err
}

// SEASONS

func (r *AdminRepo) CreateSeason(projectID int, req model.CreateSeasonRequest) (*model.Season, error) {
	var s model.Season
	err := r.db.QueryRow(`
		INSERT INTO seasons (project_id, season_number) VALUES ($1, $2)
		RETURNING id, project_id, season_number
	`, projectID, req.SeasonNumber).Scan(&s.ID, &s.ProjectID, &s.SeasonNumber)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// EPISODES

func (r *AdminRepo) CreateEpisode(seasonID int, req model.CreateEpisodeRequest) (*model.Episode, error) {
	var e model.Episode
	err := r.db.QueryRow(`
		INSERT INTO episodes (season_id, episode_number, title, youtube_video_id, duration)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, season_id, episode_number, title, youtube_video_id, duration
	`, seasonID, req.EpisodeNumber, nn(req.Title), req.YouTubeVideoID, req.Duration).
		Scan(&e.ID, &e.SeasonID, &e.EpisodeNumber, &e.Title, &e.YouTubeVideoID, &e.Duration)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// FEATURED ORDER

func (r *AdminRepo) UpdateFeaturedOrder(blockType string, items []struct {
	ProjectID int
	SortOrder int
}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM featured_content WHERE block_type = $1`, blockType)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.Exec(
			`INSERT INTO featured_content (project_id, block_type, sort_order) VALUES ($1, $2, $3)`,
			item.ProjectID, blockType, item.SortOrder,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// REFERENCES CRUD

func (r *AdminRepo) CreateCategory(name string) (*model.Category, error) {
	var c model.Category
	err := r.db.QueryRow(`INSERT INTO categories (name) VALUES ($1) RETURNING id, name`, name).Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *AdminRepo) UpdateCategory(id int, name string) error {
	_, err := r.db.Exec(`UPDATE categories SET name = $1 WHERE id = $2`, name, id)
	return err
}

func (r *AdminRepo) DeleteCategory(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}

func (r *AdminRepo) CreateGenre(name string) (*model.Genre, error) {
	var g model.Genre
	err := r.db.QueryRow(`INSERT INTO genres (name) VALUES ($1) RETURNING id, name`, name).Scan(&g.ID, &g.Name)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *AdminRepo) UpdateGenre(id int, name string) error {
	_, err := r.db.Exec(`UPDATE genres SET name = $1 WHERE id = $2`, name, id)
	return err
}

func (r *AdminRepo) DeleteGenre(id int) error {
	_, err := r.db.Exec(`DELETE FROM genres WHERE id = $1`, id)
	return err
}

func (r *AdminRepo) CreateAgeRating(rng string) (*model.AgeRating, error) {
	var a model.AgeRating
	err := r.db.QueryRow(`INSERT INTO age_ratings (range) VALUES ($1) RETURNING id, range`, rng).Scan(&a.ID, &a.Range)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepo) UpdateAgeRating(id int, rng string) error {
	_, err := r.db.Exec(`UPDATE age_ratings SET range = $1 WHERE id = $2`, rng, id)
	return err
}

func (r *AdminRepo) DeleteAgeRating(id int) error {
	_, err := r.db.Exec(`DELETE FROM age_ratings WHERE id = $1`, id)
	return err
}

// USERS

func (r *AdminRepo) GetAllUsers(page, limit int) ([]model.User, int, error) {
	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT u.id, u.email, u.full_name, u.phone, u.birth_date, u.user_avatar_url,
			   u.role_id, COALESCE(r.name, ''), u.created_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Phone, &u.BirthDate,
			&u.UserAvatarURL, &u.RoleID, &u.RoleName, &u.CreatedAt)
		users = append(users, u)
	}
	return users, total, nil
}

func (r *AdminRepo) AssignRole(userID, roleID int) error {
	_, err := r.db.Exec(`UPDATE users SET role_id = $1 WHERE id = $2`, roleID, userID)
	return err
}

// ROLES

func (r *AdminRepo) GetRoles() ([]model.Role, error) {
	rows, err := r.db.Query(`SELECT id, name, permissions FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		var perms []byte
		rows.Scan(&role.ID, &role.Name, &perms)
		role.Permissions = string(perms)
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *AdminRepo) CreateRole(name, permissions string) (*model.Role, error) {
	var role model.Role
	err := r.db.QueryRow(
		`INSERT INTO roles (name, permissions) VALUES ($1, $2) RETURNING id, name, permissions`,
		name, permissions,
	).Scan(&role.ID, &role.Name, &role.Permissions)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *AdminRepo) UpdateRole(id int, name, permissions string) error {
	_, err := r.db.Exec(`UPDATE roles SET name=$1, permissions=$2 WHERE id=$3`, name, permissions, id)
	return err
}

func (r *AdminRepo) DeleteRole(id int) error {
	_, err := r.db.Exec(`DELETE FROM roles WHERE id = $1`, id)
	return err
}

// хелпер: nil -> NULL
func nn(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}