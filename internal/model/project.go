package model

import "time"

type ProjectType string

const (
	ProjectTypeMovie  ProjectType = "MOVIE"
	ProjectTypeSeries ProjectType = "SERIES"
)

type Project struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	ReleaseYear     *int       `json:"release_year"`
	Director        *string    `json:"director"`
	Producer        *string    `json:"producer"`
	CoverImageURL   *string    `json:"cover_image_url"`
	BannerImageURL  *string    `json:"banner_image_url"`
	Keywords        *string    `json:"keywords"`
	CategoryID      *int       `json:"category_id"`
	CategoryName    string     `json:"category_name,omitempty"`
	AgeRatingID     *int       `json:"age_rating_id"`
	AgeRatingRange  string     `json:"age_rating_range,omitempty"`
	ProjectType     string     `json:"project_type"`
	DurationMinutes *int       `json:"duration_minutes"`
	YouTubeVideoID  *string    `json:"youtube_video_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Genres          []Genre    `json:"genres,omitempty"`
	Seasons         []Season   `json:"seasons,omitempty"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AgeRating struct {
	ID    int    `json:"id"`
	Range string `json:"range"`
}

type Season struct {
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	SeasonNumber int       `json:"season_number"`
	Episodes     []Episode `json:"episodes,omitempty"`
}

type Episode struct {
	ID             int     `json:"id"`
	SeasonID       int     `json:"season_id"`
	EpisodeNumber  int     `json:"episode_number"`
	Title          *string `json:"title"`
	YouTubeVideoID string  `json:"youtube_video_id"`
	Duration       *int    `json:"duration"`
}

type Favorite struct {
	UserID    int       `json:"user_id"`
	ProjectID int       `json:"project_id"`
	AddedAt   time.Time `json:"added_at"`
	Project   *Project  `json:"project,omitempty"`
}

type FeaturedBlock struct {
	BlockType string    `json:"block_type"`
	Projects  []Project `json:"projects"`
}

type CreateProjectRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	ReleaseYear     *int   `json:"release_year"`
	Director        string `json:"director"`
	Producer        string `json:"producer"`
	CoverImageURL   string `json:"cover_image_url"`
	BannerImageURL  string `json:"banner_image_url"`
	Keywords        string `json:"keywords"`
	CategoryID      *int   `json:"category_id"`
	AgeRatingID     *int   `json:"age_rating_id"`
	ProjectType     string `json:"project_type" binding:"required"`
	DurationMinutes *int   `json:"duration_minutes"`
	YouTubeVideoID  string `json:"youtube_video_id"`
	GenreIDs        []int  `json:"genre_ids"`
}

type CreateSeasonRequest struct {
	SeasonNumber int `json:"season_number" binding:"required"`
}

type CreateEpisodeRequest struct {
	EpisodeNumber  int    `json:"episode_number" binding:"required"`
	Title          string `json:"title"`
	YouTubeVideoID string `json:"youtube_video_id" binding:"required"`
	Duration       *int   `json:"duration"`
}

type UpdateFeaturedOrderRequest struct {
	BlockType string `json:"block_type" binding:"required"`
	Items     []struct {
		ProjectID int `json:"project_id" binding:"required"`
		SortOrder int `json:"sort_order" binding:"required"`
	} `json:"items" binding:"required"`
}