package repository

import (
	"database/sql"

	"ozinse/internal/model"
)

type ReferenceRepo struct {
	db *sql.DB
}

func NewReferenceRepo(db *sql.DB) *ReferenceRepo {
	return &ReferenceRepo{db: db}
}

func (r *ReferenceRepo) GetCategories() ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Category
	for rows.Next() {
		var c model.Category
		rows.Scan(&c.ID, &c.Name)
		list = append(list, c)
	}
	return list, nil
}

func (r *ReferenceRepo) GetGenres() ([]model.Genre, error) {
	rows, err := r.db.Query(`SELECT id, name FROM genres ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Genre
	for rows.Next() {
		var g model.Genre
		rows.Scan(&g.ID, &g.Name)
		list = append(list, g)
	}
	return list, nil
}

func (r *ReferenceRepo) GetAgeRatings() ([]model.AgeRating, error) {
	rows, err := r.db.Query(`SELECT id, range FROM age_ratings ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AgeRating
	for rows.Next() {
		var a model.AgeRating
		rows.Scan(&a.ID, &a.Range)
		list = append(list, a)
	}
	return list, nil
}