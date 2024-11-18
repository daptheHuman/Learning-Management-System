package curriculum

import (
	"context"

	"gorm.io/gorm"

	model "github.com/dapthehuman/learning-management-system/database/models"
)

type Curriculum struct {
	DB *gorm.DB
}

func NewCurriculum(db *gorm.DB) *Curriculum {
	return &Curriculum{
		DB: db,
	}
}

func (r *Curriculum) CreateCurriculum(ctx context.Context, courseID uint64, createDTO *model.Curriculum) error {
	query := `INSERT INTO curriculums (course_id, 
	section_name, 'order', created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	err := r.DB.Exec(query, courseID, createDTO.SectionName, createDTO.SectionOrder, createDTO.CreatedAt, createDTO.UpdatedAt).Error
	return err
}

func (r *Curriculum) GetCurriculum(ctx context.Context, id uint64) ([]*model.Curriculum, error) {
	query := `SELECT id, course_id, section_name, 'order', created_at, updated_at FROM curriculums WHERE course_id = ?`
	rows, err := r.DB.Raw(query, id).Rows()
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
}

func (r *Curriculum) UpdateCurriculum(ctx context.Context, curriculum *model.Curriculum) error {
	query := `UPDATE curriculums SET section_name = ?, 'order' = ?, updated_at = ? WHERE id = ?`
	err := r.DB.Exec(query, curriculum.SectionName, curriculum.SectionOrder, curriculum.UpdatedAt, curriculum.ID).Error
	return err
}

func (r *Curriculum) DeleteCurriculum(ctx context.Context, id uint64) error {
	query := `DELETE FROM curriculums WHERE id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}
