package user

import (
	"context"

	"gorm.io/gorm"

	model "github.com/dapthehuman/learning-management-system/database/models"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		DB: db,
	}
}

func (r *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?) returning id`
	row := r.DB.Raw(query, user.Name, user.Email, user.PasswordHash, user.Role).Row()

	err := row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, role, password_hash, created_at, updated_at FROM users WHERE email = ?`
	row := r.DB.Raw(query, email).Row()

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *User) Update(ctx context.Context, user *model.User) (*model.User, error) {
	query := `UPDATE users SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? returning updated_at`
	err := r.DB.Raw(query, user.Name, user.Email, user.ID).
		Scan(&user.UpdatedAt)

	if err.Error != nil {
		return nil, err.Error
	}

	return user, nil
}

func (r *User) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM users WHERE id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}
