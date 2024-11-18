package assessment

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/dapthehuman/learning-management-system/database/models"
	model "github.com/dapthehuman/learning-management-system/database/models"
)

type Assessment struct {
	DB *gorm.DB
}

func NewAssessment(db *gorm.DB) *Assessment {
	return &Assessment{
		DB: db,
	}
}

func (r *Assessment) Create(ctx context.Context, assessment *model.Assessment) (*model.Assessment, error) {
	query := `INSERT INTO assessments (course_id, type, question, created_at) 
	VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	rows, err := r.DB.Raw(query, assessment.CourseID, assessment.Type, assessment.Question, time.Now()).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id uint64
		var createdAt time.Time
		err := rows.Scan(&id, &createdAt)
		if err != nil {
			return nil, err
		}
		assessment.ID = id
		assessment.CreatedAt = createdAt
	}

	return assessment, nil
}

func (r *Assessment) GetAllByCourseID(ctx context.Context, courseID uint64) ([]*model.Assessment, error) {
	query := `SELECT id, course_id, type, question, created_at FROM assessments WHERE course_id = $1`
	rows, err := r.DB.Raw(query, courseID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assessments := make([]*model.Assessment, 0)
	for rows.Next() {
		var assessment model.Assessment
		err := rows.Scan(&assessment.ID, &assessment.CourseID, &assessment.Type, &assessment.Question, &assessment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, &assessment)
	}

	return assessments, nil
}

func (r *Assessment) GetByID(ctx context.Context, assessmentID uint64) (*model.Assessment, error) {
	query := `SELECT id, course_id, type, question, created_at FROM assessments WHERE id = $1`
	row := r.DB.Raw(query, assessmentID).Row()

	var assessment model.Assessment
	err := row.Scan(&assessment.ID, &assessment.CourseID, &assessment.Type, &assessment.Question, &assessment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &assessment, nil
}

func (r *Assessment) SubmitAnswer(ctx context.Context, submission *models.Submission) (*model.Submission, error) {
	grade := autoGrade(submission.AssessmentID, submission.Answer)

	query := `INSERT INTO submissions (user_id, assessment_id, answer, grade, submitted_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, submitted_at`
	rows := r.DB.Raw(query, submission.UserID, submission.AssessmentID, submission.Answer, grade, time.Now()).Row()
	err := rows.Scan(&submission.ID, &submission.SubmittedAt)
	if err != nil {
		return nil, err
	}

	submission.Grade = grade
	return submission, nil
}

func autoGrade(i int, s string) int {
	return 100
}
