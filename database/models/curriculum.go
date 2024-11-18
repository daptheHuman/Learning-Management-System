package models

import (
	"time"
)

type Curriculum struct {
	ID           uint64    `json:"id" db:"id"`
	CourseID     uint64    `json:"course_id" db:"course_id"`
	SectionName  string    `json:"section_name" db:"section_name"`
	SectionOrder int       `json:"section_order" db:"section_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
