package dto

import "time"

type CreateMaterialRequest struct {
	CurriculumID uint64 `json:"curriculum_id" binding:"required"`
	MaterialType string `json:"material_type" binding:"required,oneof=text video quiz"`
	Content      string `json:"content" binding:"required"`
	Order        int    `json:"order" binding:"required"`
}

type CreateMaterialResponse struct {
	ID           int       `json:"id"`
	CurriculumID int       `json:"curriculum_id"`
	MaterialType string    `json:"material_type"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateMaterialRequest struct {
	MaterialType string `json:"material_type" binding:"omitempty,oneof=text video quiz"`
	Content      string `json:"content" binding:"omitempty"`
	Order        int    `json:"order" binding:"omitempty"`
}

type MaterialResponse struct {
	ID           int       `json:"id"`
	CurriculumID int       `json:"curriculum_id"`
	MaterialType string    `json:"material_type"`
	Content      string    `json:"content"`
	Order        int       `json:"order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
