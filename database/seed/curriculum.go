package seed

import (
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/go-faker/faker/v4"
)

func CurriculumGenerator() *model.Curriculum {
	a := &model.Curriculum{}
	a.CourseID = uint64(1)
	a.SectionName = faker.Sentence()
	a.SectionOrder = 0
	return a
}
