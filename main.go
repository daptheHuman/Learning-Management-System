package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	assessmentRepo "github.com/dapthehuman/learning-management-system/database/repositories/assessment"
	courseRepo "github.com/dapthehuman/learning-management-system/database/repositories/course"
	materialRepo "github.com/dapthehuman/learning-management-system/database/repositories/material"
	studentRepo "github.com/dapthehuman/learning-management-system/database/repositories/student"
	userRepo "github.com/dapthehuman/learning-management-system/database/repositories/user"

	assessmentService "github.com/dapthehuman/learning-management-system/service/assessment-service"
	courseService "github.com/dapthehuman/learning-management-system/service/course-service"
	materialService "github.com/dapthehuman/learning-management-system/service/material-service"
	"github.com/dapthehuman/learning-management-system/service/redis"
	studentService "github.com/dapthehuman/learning-management-system/service/student-service"
	userService "github.com/dapthehuman/learning-management-system/service/user-service"

	seeders "github.com/dapthehuman/learning-management-system/database/seed"
	"github.com/dapthehuman/learning-management-system/http/route"

	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/fsutil"

	// Import the appropriate GORM dialect for the database you're using.
	// _ "goyave.dev/goyave/v5/database/dialect/mysql"
	_ "goyave.dev/goyave/v5/database/dialect/postgres"
	// _ "goyave.dev/goyave/v5/database/dialect/sqlite"
	// _ "goyave.dev/goyave/v5/database/dialect/mssql"
	// _ "goyave.dev/goyave/v5/database/dialect/clickhouse"
)

//go:embed resources
var resources embed.FS

func main() {
	var seed bool
	flag.BoolVar(&seed, "seed", false, "If true, the database will be seeded with random data.")
	flag.Parse()
	resources := fsutil.NewEmbed(resources)
	langFS, err := resources.Sub("resources/lang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.(*errors.Error).String())
		os.Exit(1)
	}

	opts := goyave.Options{
		LangFS: langFS,
	}

	server, err := goyave.New(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.(*errors.Error).String())
		os.Exit(1)
	}

	server.Logger.Info("Registering hooks")
	server.RegisterSignalHook()

	server.RegisterStartupHook(func(s *goyave.Server) {
		server.Logger.Info("Server is listening", "host", s.Host())
	})

	server.RegisterShutdownHook(func(s *goyave.Server) {
		s.Logger.Info("Server is shutting down")
	})

	registerServices(server)

	server.Logger.Info("Registering routes")
	server.RegisterRoutes(route.Register)

	if seed {
		server.Logger.Info("Seeding database...")
		seeders.Seed(server.DB())
		server.Logger.Info("Database seeded")
		os.Exit(0)
	}

	if err := server.Start(); err != nil {
		server.Logger.Error(err)
		os.Exit(2)
	}
}

func registerServices(server *goyave.Server) {
	server.Logger.Info("Registering services")

	// Services represent the Domain/Business layer.
	// This is where the core logic and value of your application resides.
	// This function is where you will register your services in the server's
	// service container to make them accessible to dependents.
	// https://goyave.dev/basics/services.html#service-container

	redis := redis.NewClient(
		fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)

	// TODO register services
	userRepository := userRepo.NewUser(server.DB())
	server.RegisterService(userService.NewService(userRepository))

	studentRepository := studentRepo.NewStudent(server.DB())
	server.RegisterService(studentService.NewService(studentRepository))

	courseRepository := courseRepo.NewCourse(server.DB(), redis)
	server.RegisterService(courseService.NewService(courseRepository))

	materialRepository := materialRepo.NewMaterial(server.DB(), redis)
	server.RegisterService(materialService.NewService(materialRepository))

	assessmentRepository := assessmentRepo.NewAssessment(server.DB())
	server.RegisterService(assessmentService.NewService(assessmentRepository))
}
