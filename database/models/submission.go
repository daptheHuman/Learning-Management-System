package models

import (
	"time"
)

type Submission struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	AssessmentID int       `json:"assessment_id" db:"assessment_id"`
	Answer       string    `json:"answer" db:"answer"`
	Grade        int       `json:"grade" db:"grade"`
	SubmittedAt  time.Time `json:"submitted_at" db:"submitted_at"`
}
