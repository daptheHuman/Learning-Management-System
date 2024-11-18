package seed

import (
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/go-faker/faker/v4"
)

func MaterialGenerator() *model.Material {
	a := &model.Material{}
	a.CurriculumID = 1
	a.MaterialType = "video"
	a.Content = faker.URL()
	a.Order = 0

	return a
}
