package models

import (
	"time"
)

type Assessment struct {
	ID        uint64    `json:"id" db:"id"`
	CourseID  uint64    `json:"course_id" db:"course_id"`
	Type      string    `json:"type" db:"type"` // e.g., "multiple-choice", "essay"
	Question  string    `json:"question" db:"question"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
