package dto

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type Enrollment struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	CourseID   int    `json:"course_id"`
	EnrolledAt string `json:"enrolled_at"`
}

type Achievement struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	CourseID        int    `json:"course_id"`
	AchievementType string `json:"achievement_type"` // e.g., "course completion"
	Description     string `json:"description"`
	AwardedAt       string `json:"awarded_at"`
}

type CreateAchievementRequest struct {
	UserID          int    `json:"user_id"`
	CourseID        int    `json:"course_id"`
	AchievementType string `json:"achievement_type"`
	Description     string `json:"description"`
}

type UpdateAchievementRequest struct {
	AchievementType string `json:"achievement_type"`
	Description     string `json:"description"`
}

type ProgressTracking struct {
	ID                 int    `json:"id"`
	UserID             int    `json:"user_id"`
	CurriculumID       int    `json:"curriculum_id"`
	MaterialID         int    `json:"material_id"`
	Status             string `json:"status"` // e.g., "in-progress", "completed"
	ProgressPercentage int    `json:"progress_percentage"`
}

type TrackProgressRequest struct {
	UserID       int    `json:"user_id"`
	CurriculumID int    `json:"curriculum_id"`
	MaterialID   int    `json:"material_id"`
	Status       string `json:"status"`
	Progress     int    `json:"progress_percentage"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateStudentRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type EnrollStudentRequest struct {
	UserID   uint64 `json:"user_id"`
	CourseID uint64 `json:"course_id"`
}

type CreateLearningPathRequest struct {
	UserID       int   `json:"user_id" validate:"required"`
	CourseID     int   `json:"course_id" validate:"required"`
	ModulesOrder []int `json:"modules_order" validate:"required,dive,gt=0"`
}

type StartLearningPath struct {
	UserID   int `json:"user_id"`
	CourseID int `json:"course_id"`
}

type UpdateProgress struct {
	CurriculumID       int    `json:"curriculum_id"`
	Status             string `json:"status"` // e.g., "completed"
	ProgressPercentage int    `json:"progress_percentage"`
}

type LearningPathResponse struct {
	CurriculumID int    `json:"curriculum_id"`
	SectionName  string `json:"section_name"`
	Status       string `json:"status"`
	Progress     int    `json:"progress_percentage"`
}
