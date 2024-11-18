package models

import (
	"time"
)

type Material struct {
	ID           int       `json:"id"`
	CurriculumID int       `json:"curriculum_id"`
	MaterialType string    `json:"material_type"` // e.g., "text", "video", "quiz"
	Content      string    `json:"content"`
	Order        int       `json:"order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
