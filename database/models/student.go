package models

import (
	"time"
)

type LearningPath struct {
	ID           int   `json:"id"`
	UserID       int   `json:"user_id"`
	CourseID     int   `json:"course_id"`
	ModulesOrder []int `json:"modules_order"`
}

type Enrollment struct {
	ID         int       `json:"id"`
	UserID     uint64    `json:"user_id"`
	CourseID   uint64    `json:"course_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

type ProgressTracking struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	CurriculumID       int       `json:"curriculum_id"`
	MaterialID         int       `json:"material_id"`
	Status             string    `json:"status"` // e.g., "in-progress", "completed"
	ProgressPercentage int       `json:"progress_percentage"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Achievement struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	CourseID        int       `json:"course_id"`
	AchievementType string    `json:"achievement_type"` // e.g., "course completion"
	Description     string    `json:"description"`
	AwardedAt       time.Time `json:"awarded_at"`
}

type StudentPath struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	CurriculumID       int       `json:"curriculum_id"`
	Status             string    `json:"status"` // e.g., "not started", "in-progress", "completed"
	ProgressPercentage int       `json:"progress_percentage"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
