package student

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	model "github.com/dapthehuman/learning-management-system/database/models"
)

type Student struct {
	DB *gorm.DB
}

func NewStudent(db *gorm.DB) *Student {
	return &Student{
		DB: db,
	}
}

func (r *Student) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var student model.User
	result := r.DB.First(&student, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &student, nil
}

func (r *Student) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?`
	row := r.DB.Raw(query, email).Row()

	var student model.User
	err := row.Scan(&student.ID, &student.Name, &student.Email, &student.CreatedAt, &student.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &student, nil
}

func (r *Student) Get(ctx context.Context, id uint64) (*model.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`
	row := r.DB.Raw(query, id).Row()

	var student model.User
	err := row.Scan(&student.ID, &student.Name, &student.Email, &student.CreatedAt, &student.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &student, nil
}

func (r *Student) GetAll(ctx context.Context) ([]*model.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users`
	rows, err := r.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]*model.User, 0)
	for rows.Next() {
		var student model.User
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.CreatedAt, &student.UpdatedAt)
		if err != nil {
			return nil, err
		}
		courses = append(courses, &student)
	}

	return courses, nil
}

func (r *Student) Update(ctx context.Context, student *model.User) (*model.User, error) {
	query := `UPDATE users SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING updated_at`
	rows := r.DB.Raw(query, student.Name, student.Email, student.ID).Row()
	err := rows.Scan(&student.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return student, nil
}
