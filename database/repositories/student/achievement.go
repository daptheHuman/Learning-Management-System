package student

import (
	"context"

	"github.com/dapthehuman/learning-management-system/database/models"
)

func (r *Student) ListAchievementByUserID(ctx context.Context, userID uint64) ([]*models.Achievement, error) {
	query := `SELECT id, user_id, course_id, achievement_type, description, awarded_at FROM achievements WHERE user_id = ?`
	rows, err := r.DB.Raw(query, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*models.Achievement
	for rows.Next() {
		var achievement models.Achievement
		if err := rows.Scan(&achievement.ID, &achievement.UserID, &achievement.CourseID, &achievement.AchievementType, &achievement.Description, &achievement.AwardedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, &achievement)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *Student) GetAchievementByID(ctx context.Context, id uint64) (*models.Achievement, error) {
	query := `SELECT id, user_id, course_id, achievement_type, description, awarded_at FROM achievements WHERE id = ?`
	row := r.DB.Raw(query, id).Row()

	var achievement models.Achievement
	err := row.Scan(&achievement.ID, &achievement.UserID, &achievement.CourseID, &achievement.AchievementType, &achievement.Description, &achievement.AwardedAt)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *Student) CreateAchievement(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error) {
	query := `INSERT INTO achievements (user_id, course_id, achievement_type, description, awarded_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, user_id, course_id, achievement_type, description, awarded_at`
	row := r.DB.Raw(query, achievement.UserID, achievement.CourseID, achievement.AchievementType, achievement.Description, achievement.AwardedAt).Row()

	var newAchievement models.Achievement
	err := row.Scan(&newAchievement.ID, &newAchievement.UserID, &newAchievement.CourseID, &newAchievement.AchievementType, &newAchievement.Description, &newAchievement.AwardedAt)
	if err != nil {
		return nil, err
	}

	return &newAchievement, nil
}

func (r *Student) UpdateAchievement(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error) {
	query := `UPDATE achievements SET course_id = $1, achievement_type = $2, description = $3, awarded_at = $4 WHERE id = $5 RETURNING id, user_id, course_id, achievement_type, description, awarded_at`
	row := r.DB.Raw(query, achievement.CourseID, achievement.AchievementType, achievement.Description, achievement.AwardedAt, achievement.ID).Row()

	var updatedAchievement models.Achievement
	err := row.Scan(&updatedAchievement.ID, &updatedAchievement.UserID, &updatedAchievement.CourseID, &updatedAchievement.AchievementType, &updatedAchievement.Description, &updatedAchievement.AwardedAt)
	if err != nil {
		return nil, err
	}

	return &updatedAchievement, nil
}

func (r *Student) DeleteAchievement(ctx context.Context, id uint64) error {
	query := `DELETE FROM achievements WHERE id = $1`
	result := r.DB.Exec(query, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
