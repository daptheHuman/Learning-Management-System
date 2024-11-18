package material

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/dapthehuman/learning-management-system/cache"
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/redis/go-redis/v9"
)

type Material struct {
	DB    *gorm.DB
	Redis redis.UniversalClient
}

func NewMaterial(db *gorm.DB, redis redis.UniversalClient) *Material {
	return &Material{
		DB:    db,
		Redis: redis,
	}
}

func (r *Material) Create(ctx context.Context, material *model.Material) (*model.Material, error) {
	query := `INSERT INTO materials (curriculum_id, material_type, content, "order", created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	row := r.DB.Raw(query, material.CurriculumID, material.MaterialType, material.Content, material.Order, material.CreatedAt, material.UpdatedAt).Row()
	err := row.Scan(&material.ID, &material.CreatedAt, &material.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return material, nil
}

func (r *Material) GetByID(ctx context.Context, id uint64) (*model.Material, error) {
	cacheKey := fmt.Sprintf("material:%d", id)

	return cache.Cache(ctx, r.Redis, cacheKey, func() (*model.Material, error) {
		query := `SELECT id, curriculum_id, material_type, content, "order", created_at, updated_at FROM materials WHERE id = $1`
		row := r.DB.Raw(query, id).Row()

		material := &model.Material{}
		err := row.Scan(&material.ID, &material.CurriculumID, &material.MaterialType, &material.Content, &material.Order, &material.CreatedAt, &material.UpdatedAt)
		if err != nil {
			return nil, err
		}

		return material, nil
	})
}

func (r *Material) GetByCurriculumID(ctx context.Context, curriculumID uint64) ([]*model.Material, error) {
	cacheKey := fmt.Sprintf("materials:curriculum:%d", curriculumID)

	return cache.Cache(ctx, r.Redis, cacheKey, func() ([]*model.Material, error) {
		query := `SELECT id, curriculum_id, material_type, content, "order", created_at, updated_at FROM materials WHERE curriculum_id = $1`
		rows, err := r.DB.Raw(query, curriculumID).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		materials := []*model.Material{}
		for rows.Next() {
			material := &model.Material{}
			err := rows.Scan(&material.ID, &material.CurriculumID, &material.MaterialType, &material.Content, &material.Order, &material.CreatedAt, &material.UpdatedAt)
			if err != nil {
				return nil, err
			}
			materials = append(materials, material)
		}

		return materials, nil
	})
}

func (r *Material) Update(ctx context.Context, material *model.Material) (*model.Material, error) {
	query := `UPDATE materials SET curriculum_id = $1, material_type = $2, content = $3, "order" = $4, updated_at = $5 WHERE id = $6 RETURNING id, created_at, updated_at`
	row := r.DB.Raw(query, material.CurriculumID, material.MaterialType, material.Content, material.Order, material.UpdatedAt, material.ID).Row()
	err := row.Scan(&material.ID, &material.CreatedAt, &material.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return material, nil
}

func (r *Material) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM materials WHERE id = $1`
	result := r.DB.Exec(query, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
