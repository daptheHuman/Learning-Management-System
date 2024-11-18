package seed

import (
	model "github.com/dapthehuman/learning-management-system/database/models"
	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
)

func StudentGenerator() *model.User {
	a := &model.User{}
	a.Name = faker.Name()
	a.Email = faker.Email()
	a.Role = "student"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	a.PasswordHash = string(hashedPassword)
	return a
}

func InstructorGenerator() *model.User {
	a := &model.User{}
	a.Name = faker.Name()
	a.Email = faker.Email()
	a.Role = "instructor"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	a.PasswordHash = string(hashedPassword)
	return a
}

func AdminGenerator() *model.User {
	a := &model.User{}
	a.Name = faker.Name()
	a.Email = "admin@admin.com"
	a.Role = "admin"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	a.PasswordHash = string(hashedPassword)
	return a
}
