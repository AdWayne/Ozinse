package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"ozinse/internal/model"

	"github.com/lib/pq"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) GetAll(search *string, categoryID, genreID, ageRatingID *int, projectType *string, page, limit int) ([]model.Project, int, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if search != nil && *search != "" {
		where = append(where, fmt.Sprintf("(p.title ILIKE $%d OR p.keywords ILIKE $%d)", idx, idx+1))
		args = append(args, "%"+*search+"%", "%"+*search+"%")
		idx += 2
	}
	if categoryID != nil {
		where = append(where, fmt.Sprintf("p.category_id = $%d", idx))
		args = append(args, *categoryID)
		idx++
	}
	if ageRatingID != nil {
		where = append(where, fmt.Sprintf("p.age_rating_id = $%d", idx))
		args = append(args, *ageRatingID)
		idx++
	}
	if projectType != nil {
		where = append(where, fmt.Sprintf("p.project_type = $%d", idx))
		args = append(args, *projectType)
		idx++
	}
	if genreID != nil {
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM project_genres pg WHERE pg.project_id = p.id AND pg.genre_id = $%d)", idx))
		args = append(args, *genreID)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	err := r.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM projects p WHERE %s", whereClause), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT p.id, p.title, p.description, p.release_year, p.director, p.producer,
			   p.cover_image_url, p.banner_image_url, p.keywords, p.category_id,
			   COALESCE(c.name, ''), p.age_rating_id, COALESCE(a.range, ''),
			   p.project_type, p.duration_minutes, p.youtube_video_id,
			   p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN age_ratings a ON p.age_rating_id = a.id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, idx, idx+1)

	args = append(args, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ReleaseYear, &p.Director,
			&p.Producer, &p.CoverImageURL, &p.BannerImageURL, &p.Keywords,
			&p.CategoryID, &p.CategoryName, &p.AgeRatingID, &p.AgeRatingRange,
			&p.ProjectType, &p.DurationMinutes, &p.YouTubeVideoID,
			&p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}

	// Подгружаем жанры для всех проектов
	if len(projects) > 0 {
		ids := make([]int, len(projects))
		for i, p := range projects {
			ids[i] = p.ID
		}

		genreRows, err := r.db.Query(`
			SELECT pg.project_id, g.id, g.name
			FROM project_genres pg
			JOIN genres g ON pg.genre_id = g.id
			WHERE pg.project_id = ANY($1)
		`, pq.Array(ids))
		if err == nil {
			defer genreRows.Close()
			genreMap := make(map[int][]model.Genre)
			for genreRows.Next() {
				var projectID, genreID int
				var genreName string
				genreRows.Scan(&projectID, &genreID, &genreName)
				genreMap[projectID] = append(genreMap[projectID], model.Genre{ID: genreID, Name: genreName})
			}
			for i := range projects {
				projects[i].Genres = genreMap[projects[i].ID]
			}
		}
	}

	return projects, total, nil
}

func (r *ProjectRepo) GetByID(id int) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(`
		SELECT p.id, p.title, p.description, p.release_year, p.director, p.producer,
			   p.cover_image_url, p.banner_image_url, p.keywords, p.category_id,
			   COALESCE(c.name, ''), p.age_rating_id, COALESCE(a.range, ''),
			   p.project_type, p.duration_minutes, p.youtube_video_id,
			   p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN age_ratings a ON p.age_rating_id = a.id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.Title, &p.Description, &p.ReleaseYear, &p.Director,
		&p.Producer, &p.CoverImageURL, &p.BannerImageURL, &p.Keywords,
		&p.CategoryID, &p.CategoryName, &p.AgeRatingID, &p.AgeRatingRange,
		&p.ProjectType, &p.DurationMinutes, &p.YouTubeVideoID,
		&p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// жанры
	rows, err := r.db.Query(`
		SELECT g.id, g.name FROM genres g
		JOIN project_genres pg ON g.id = pg.genre_id
		WHERE pg.project_id = $1
	`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var g model.Genre
			rows.Scan(&g.ID, &g.Name)
			p.Genres = append(p.Genres, g)
		}
	}

	return &p, nil
}

func (r *ProjectRepo) GetSeasonsWithEpisodes(projectID int) ([]model.Season, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.project_id, s.season_number,
			   e.id, e.season_id, e.episode_number, e.title, e.youtube_video_id, e.duration
		FROM seasons s
		LEFT JOIN episodes e ON s.id = e.season_id
		WHERE s.project_id = $1
		ORDER BY s.season_number, e.episode_number
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seasonMap := make(map[int]*model.Season)
	var seasonOrder []int

	for rows.Next() {
		var sID, sProjectID, sNumber int
		var eID, eSeasonID, eNumber sql.NullInt64
		var eTitle, eYouTubeID sql.NullString
		var eDuration sql.NullInt64

		err := rows.Scan(&sID, &sProjectID, &sNumber, &eID, &eSeasonID, &eNumber, &eTitle, &eYouTubeID, &eDuration)
		if err != nil {
			return nil, err
		}

		if _, ok := seasonMap[sID]; !ok {
			seasonMap[sID] = &model.Season{
				ID:           sID,
				ProjectID:    sProjectID,
				SeasonNumber: sNumber,
				Episodes:     []model.Episode{},
			}
			seasonOrder = append(seasonOrder, sID)
		}

		if eID.Valid {
			ep := model.Episode{
				ID:             int(eID.Int64),
				SeasonID:       int(eSeasonID.Int64),
				EpisodeNumber:  int(eNumber.Int64),
				YouTubeVideoID: eYouTubeID.String,
			}
			if eTitle.Valid {
				ep.Title = &eTitle.String
			}
			if eDuration.Valid {
				d := int(eDuration.Int64)
				ep.Duration = &d
			}
			seasonMap[sID].Episodes = append(seasonMap[sID].Episodes, ep)
		}
	}

	var seasons []model.Season
	for _, id := range seasonOrder {
		seasons = append(seasons, *seasonMap[id])
	}
	return seasons, nil
}

func (r *ProjectRepo) GetFeatured() (map[string][]model.Project, error) {
	rows, err := r.db.Query(`
		SELECT fc.block_type, fc.sort_order, p.id, p.title, p.description, p.cover_image_url,
			   p.banner_image_url, p.project_type, p.release_year, p.duration_minutes,
			   p.created_at, p.updated_at
		FROM featured_content fc
		JOIN projects p ON fc.project_id = p.id
		ORDER BY fc.block_type, fc.sort_order
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]model.Project)
	for rows.Next() {
		var blockType string
		var sortOrder int
		var p model.Project
		err := rows.Scan(&blockType, &sortOrder, &p.ID, &p.Title, &p.Description,
			&p.CoverImageURL, &p.BannerImageURL, &p.ProjectType,
			&p.ReleaseYear, &p.DurationMinutes, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result[blockType] = append(result[blockType], p)
	}
	return result, nil
}