package seed

import (
	"gorm.io/gorm"
	"goyave.dev/goyave/v5/database"
)

// Seeders generate random records for testing purposes with the help of factories.
//
// Learn more here: https://goyave.dev/advanced/testing.html#factories

// Seed migrate all the models and populates the database.
func Seed(db *gorm.DB) {
	// TODO implement seeds

	studentFactory := database.NewFactory(StudentGenerator)
	_ = studentFactory.Generate(10)
	studentFactory.Save(db, 10)

	instructorFactory := database.NewFactory(InstructorGenerator)
	_ = instructorFactory.Generate(10)
	instructorFactory.Save(db, 10)

	adminFactory := database.NewFactory(AdminGenerator)
	_ = adminFactory.Generate(1)
	adminFactory.Save(db, 1)

	courseFactory := database.NewFactory(CourseGenerator)
	_ = courseFactory.Generate(10)
	courseFactory.Save(db, 10)

	curriculumFactory := database.NewFactory(CurriculumGenerator)
	_ = curriculumFactory.Generate(10)
	curriculumFactory.Save(db, 10)

	materialFactory := database.NewFactory(MaterialGenerator)
	_ = materialFactory.Generate(10)
	materialFactory.Save(db, 10)
}
