package seed

import (
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/go-faker/faker/v4"
)

func CourseGenerator() *model.Course {
	a := &model.Course{}
	a.Title = faker.Sentence()
	a.Description = faker.Paragraph()
	return a
}
