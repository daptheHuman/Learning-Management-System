package admin

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/dapthehuman/learning-management-system/cache"
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/redis/go-redis/v9"
)

type Admin struct {
	DB    *gorm.DB
	redis redis.UniversalClient
}

func NewAdmin(db *gorm.DB, redis redis.UniversalClient) *Admin {
	return &Admin{
		DB:    db,
		redis: redis,
	}
}

func (r *Admin) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	key := "all-users"
	query := `SELECT id, name, email, role, created_at, updated_at FROM users`

	return cache.Cache(ctx, r.redis, key, func() ([]*model.User, error) {
		rows, err := r.DB.Raw(query).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var users []*model.User
		for rows.Next() {
			var user model.User
			err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				return nil, err
			}
			users = append(users, &user)
		}

		return users, nil
	})
}

func (r *Admin) UpdateUser(ctx context.Context, userID uint64, user *model.User) (*model.User, error) {
	query := `UPDATE users SET name = ?, email = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING updated_at`
	err := r.DB.Raw(query, user.Name, user.Email, user.Role, userID).
		Scan(&user.UpdatedAt)

	if err.Error != nil {
		return nil, err.Error
	}

	return user, nil
}

func (r *Admin) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	key := fmt.Sprintf("user:%d", id)
	query := `SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = ?`
	return cache.Cache(ctx, r.redis, key, func() (*model.User, error) {
		row := r.DB.Raw(query, id).Row()

		var user model.User
		err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		return &user, nil
	})
}

func (r *Admin) DeleteUser(ctx context.Context, id uint64) error {
	query := `DELETE FROM users WHERE id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}
