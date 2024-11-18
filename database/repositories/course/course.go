package course

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/dapthehuman/learning-management-system/cache"
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/redis/go-redis/v9"
)

type Course struct {
	DB    *gorm.DB
	Redis redis.UniversalClient
}

func NewCourse(db *gorm.DB, redis redis.UniversalClient) *Course {
	return &Course{
		DB:    db,
		Redis: redis,
	}
}

func (r *Course) First(ctx context.Context, id uint64) (*model.Course, error) {
	key := fmt.Sprintf("course:%d", id)
	query := `SELECT id, title, description, created_at, updated_at FROM courses WHERE id = ?`

	return cache.Cache(ctx, r.Redis, key, func() (*model.Course, error) {
		var course model.Course
		err := r.DB.Raw(query, id).Scan(&course).Error
		if err != nil {
			return nil, err
		}
		return &course, nil
	})
}

func (r *Course) Create(ctx context.Context, course *model.Course) (*model.Course, error) {
	query := `INSERT INTO courses (title, description) VALUES (?, ?) RETURNING id`
	rows, err := r.DB.Raw(query, course.Title, course.Description).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the ID of the newly created course
	if rows.Next() {
		var id uint64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		course.ID = id
	}

	return course, nil
}

func (r *Course) GetAll(ctx context.Context) ([]*model.Course, error) {
	key := "courses:all"
	query := `SELECT id, title, description, created_at, updated_at FROM courses`

	return cache.Cache(ctx, r.Redis, key, func() ([]*model.Course, error) {
		rows, err := r.DB.Raw(query).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		courses := make([]*model.Course, 0)
		for rows.Next() {
			var course model.Course
			err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.CreatedAt, &course.UpdatedAt)
			if err != nil {
				return nil, err
			}
			courses = append(courses, &course)
		}

		return courses, nil
	})
}

func (r *Course) Update(ctx context.Context, course *model.Course) (*model.Course, error) {
	query := `UPDATE courses SET title = $1, description = $2, updated_at = $3 WHERE id = $4 
	          RETURNING updated_at`
	err := r.DB.Raw(query, course.Title, course.Description, time.Now(), course.ID).
		Scan(&course.UpdatedAt)

	if err.Error != nil {
		return nil, err.Error
	}

	return course, nil
}

func (r *Course) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM courses WHERE id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}

func (r *Course) CreateCurriculum(ctx context.Context, courseID uint64, createDTO *model.Curriculum) (*model.Curriculum, error) {
	query := `INSERT INTO curriculums (course_id, 
	section_name, section_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?) returning id`
	err := r.DB.Raw(query, courseID, createDTO.SectionName, createDTO.SectionOrder, createDTO.CreatedAt, createDTO.UpdatedAt).Scan(&createDTO.ID)

	if err.Error != nil {
		return nil, err.Error
	}

	return createDTO, nil
}

func (r *Course) GetCurriculum(ctx context.Context, courseID uint64) ([]*model.Curriculum, error) {
	key := fmt.Sprintf("curriculum:course:%d", courseID)
	query := `SELECT id, course_id, section_name, section_order, created_at, updated_at FROM curriculums WHERE course_id = ?`

	return cache.Cache(ctx, r.Redis, key, func() ([]*model.Curriculum, error) {
		rows, err := r.DB.Raw(query, courseID).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		curriculums := make([]*model.Curriculum, 0)
		for rows.Next() {
			var curriculum model.Curriculum
			err := rows.Scan(&curriculum.ID, &curriculum.CourseID, &curriculum.SectionName, &curriculum.SectionOrder, &curriculum.CreatedAt, &curriculum.UpdatedAt)
			if err != nil {
				return nil, err
			}
			curriculums = append(curriculums, &curriculum)
		}

		return curriculums, nil
	})
}

func (r *Course) GetCurriculumByID(ctx context.Context, id uint64) (*model.Curriculum, error) {
	key := fmt.Sprintf("curriculum:%d", id)
	query := `SELECT id, course_id, section_name, section_order, created_at, updated_at FROM curriculums WHERE id = ?`

	return cache.Cache(ctx, r.Redis, key, func() (*model.Curriculum, error) {
		var curriculum model.Curriculum
		err := r.DB.Raw(query, id).Scan(&curriculum).Error
		if err != nil {
			return nil, err
		}
		return &curriculum, nil
	})
}

func (r *Course) UpdateCurriculum(ctx context.Context, curriculum *model.Curriculum) (*model.Curriculum, error) {
	query := `UPDATE curriculums SET section_name = $1, section_order = $2, updated_at = $3 WHERE id = $4 
	          RETURNING updated_at`
	err := r.DB.Raw(query, curriculum.SectionName, curriculum.SectionOrder, time.Now(), curriculum.ID).Scan(&curriculum.UpdatedAt)

	if err.Error != nil {
		return nil, err.Error
	}

	return curriculum, nil
}

func (r *Course) DeleteCurriculum(ctx context.Context, id uint64) error {
	query := `DELETE FROM curriculums WHERE id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}
