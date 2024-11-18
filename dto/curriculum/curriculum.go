package dto

import "time"

type Curriculum struct {
	ID           int    `json:"id"`
	CourseID     int    `json:"course_id"`
	SectionName  string `json:"section_name"`
	SectionOrder int    `json:"section_order"`
}

type CreateCurriculumRequest struct {
	CourseID     int    `json:"course_id"`
	SectionName  string `json:"section_name"`
	SectionOrder int    `json:"section_order"`
}

type CurriculumResponse struct {
	ID           int       `json:"id"`
	CourseID     int       `json:"course_id"`
	SectionName  string    `json:"section_name"`
	SectionOrder int       `json:"section_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateCurriculumRequest struct {
	SectionName  string `json:"section_name"`
	SectionOrder int    `json:"section_order"`
}
