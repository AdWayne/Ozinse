package repository

import (
	"database/sql"

	"ozinse/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(email, passwordHash, fullName string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash, full_name) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, email, password_hash, full_name, phone, birth_date, role_id, created_at`,
		email, passwordHash, fullName,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.BirthDate, &user.RoleID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	var roleName *string
	err := r.db.QueryRow(
		`SELECT u.id, u.email, u.password_hash, u.user_avatar_url, u.full_name, 
		        u.phone, u.birth_date, u.role_id, u.created_at, r.name
		 FROM users u
		 LEFT JOIN roles r ON u.role_id = r.id
		 WHERE u.email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.UserAvatarURL,
		&user.FullName, &user.Phone, &user.BirthDate, &user.RoleID,
		&user.CreatedAt, &roleName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if roleName != nil {
		user.RoleName = *roleName
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id int) (*model.User, error) {
	var user model.User
	var roleName *string
	err := r.db.QueryRow(
		`SELECT u.id, u.email, u.password_hash, u.user_avatar_url, u.full_name, 
		        u.phone, u.birth_date, u.role_id, u.created_at, r.name
		 FROM users u
		 LEFT JOIN roles r ON u.role_id = r.id
		 WHERE u.id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.UserAvatarURL,
		&user.FullName, &user.Phone, &user.BirthDate, &user.RoleID,
		&user.CreatedAt, &roleName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if roleName != nil {
		user.RoleName = *roleName
	}
	return &user, nil
}

func (r *UserRepo) UpdateProfile(id int, fullName, phone, birthDate string) error {
	_, err := r.db.Exec(
		`UPDATE users SET full_name = $1, phone = $2, birth_date = $3 WHERE id = $4`,
		fullName, phone, birthDate, id,
	)
	return err
}

func (r *UserRepo) UpdatePassword(id int, passwordHash string) error {
	_, err := r.db.Exec(
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		passwordHash, id,
	)
	return err
}

func (r *UserRepo) GetRoleByName(name string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM roles WHERE name = $1`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) CreateWithRole(email, passwordHash, fullName, avatarURL string, roleID int) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash, full_name, user_avatar_url, role_id) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, email, password_hash, user_avatar_url, full_name, phone, birth_date, role_id, created_at`,
		email, passwordHash, fullName, avatarURL, roleID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.UserAvatarURL,
		&user.FullName, &user.Phone, &user.BirthDate, &user.RoleID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
