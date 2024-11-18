package dto

import "time"

type SubmissionResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	AssessmentID int       `json:"assessment_id"`
	Answer       string    `json:"answer"`
	Grade        int       `json:"grade"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

type SubmissionRequest struct {
	UserID       int    `json:"user_id" binding:"required"`
	AssessmentID int    `json:"assessment_id" binding:"required"`
	Answer       string `json:"answer" binding:"required"`
}
