package student

import (
	"context"
	"database/sql"
	"time"

	"github.com/dapthehuman/learning-management-system/database/models"
)

func (r *Student) EnrollCourse(ctx context.Context, studentID uint64, courseID uint64) (*models.Enrollment, error) {
	query := `INSERT INTO enrollments (user_id, course_id, enrolled_at) VALUES ($1, $2, $3) RETURNING id, course_id, user_id, enrolled_at`
	row := r.DB.Raw(query, studentID, courseID, time.Now()).Row()

	var enrollment models.Enrollment
	err := row.Scan(&enrollment.ID, &enrollment.CourseID, &enrollment.UserID, &enrollment.EnrolledAt)
	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (r *Student) GetEnrollmentsByUserID(ctx context.Context, studentID uint64) ([]*models.Enrollment, error) {
	query := `SELECT id, user_id, course_id, enrolled_at FROM enrollments WHERE user_id = $1`
	rows, err := r.DB.Raw(query, studentID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		var enrollment models.Enrollment
		if err := rows.Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.EnrolledAt); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, &enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *Student) TrackProgress(ctx context.Context, progress *models.ProgressTracking) (*models.ProgressTracking, error) {
	var existingID int
	query := `SELECT id FROM progress_tracking WHERE user_id = $1 AND curriculum_id = $2 AND material_id = $3`
	row := r.DB.Raw(query, progress.UserID, progress.CurriculumID, progress.MaterialID).Row()
	err := row.Scan(&existingID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Create new progress record
			query := `INSERT INTO progress_tracking (user_id, curriculum_id, material_id, status, progress_percentage, updated_at) 
					  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, updated_at`
			row := r.DB.Raw(query, progress.UserID, progress.CurriculumID, progress.MaterialID, progress.Status, progress.ProgressPercentage, time.Now()).Row()
			err = row.Scan(&progress.ID, &progress.UpdatedAt)
			if err != nil {
				return nil, err
			}

			return progress, nil
		}
	}

	query = `UPDATE progress_tracking SET status = $1, progress_percentage = $2, updated_at = $3 WHERE id = $4 RETURNING id, updated_at`
	row = r.DB.Raw(query, progress.Status, progress.ProgressPercentage, time.Now(), existingID).Row()
	err = row.Scan(&progress.ID, &progress.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *Student) GetProgressByStudentAndCurriculum(ctx context.Context, studentID uint64, curriculumID uint64) ([]*models.ProgressTracking, error) {
	query := `SELECT id, user_id, curriculum_id, material_id, status, progress_percentage, updated_at FROM
	 progress_tracking WHERE user_id = $1 AND curriculum_id = $2`
	rows, err := r.DB.Raw(query, studentID, curriculumID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []*models.ProgressTracking
	for rows.Next() {
		var p models.ProgressTracking
		if err := rows.Scan(&p.ID, &p.UserID, &p.CurriculumID, &p.MaterialID, &p.Status, &p.ProgressPercentage, &p.UpdatedAt); err != nil {
			return nil, err
		}
		progress = append(progress, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return progress, nil
}
