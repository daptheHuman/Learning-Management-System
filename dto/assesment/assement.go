package dto

type Assessment struct {
	ID        int    `json:"id"`
	CourseID  int    `json:"course_id"`
	Type      string `json:"type"` // e.g., "multiple-choice", "essay"
	Question  string `json:"question"`
	CreatedAt string `json:"created_at"`
}

type CreateAssessmentRequest struct {
	CourseID int    `json:"course_id" binding:"required"`
	Type     string `json:"type" binding:"required"` // e.g., "multiple-choice", "essay"
	Question string `json:"question" binding:"required"`
}

type UpdateAssessmentRequest struct {
	Type     string `json:"type" binding:"required"` // e.g., "multiple-choice", "essay"
	Question string `json:"question" binding:"required"`
}

type AssessmentResponse struct {
	ID        int    `json:"id"`
	CourseID  int    `json:"course_id"`
	Type      string `json:"type"` // e.g., "multiple-choice", "essay"
	Question  string `json:"question"`
	CreatedAt string `json:"created_at"`
}
