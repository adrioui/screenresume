package main

import (
	"log"
	"screenresume/internal/controller"
	"screenresume/internal/services"
	"screenresume/pkg/db"

	"github.com/go-fuego/fuego"
)

func main() {
	// Initialize the database connection pool
	dbPool, err := db.CreateDBConnectionPool()
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v", err)
	}

	// Create a new store and service
	store := db.NewStore(dbPool)

	filesService := services.NewFileService(store)
	filesController := controller.FilesResources{FilesService: filesService}

	jobRolesService := services.NewJobRoleService(store)
	jobRolesController := controller.JobRolesResources{JobRolesService: jobRolesService}

	departmentsService := services.NewDepartmentService(store)
	departmentsController := controller.DepartmentsResources{DepartmentsService: departmentsService}

	candidatesService := services.NewCandidateService(store)
	candidatesController := controller.CandidatesResources{CandidatesService: candidatesService}

	skillsService := services.NewSkillService(store)
	skillsController := controller.SkillsResources{SkillsService: skillsService}

	s := fuego.NewServer()

	filesController.Routes(s)
	jobRolesController.Routes(s)
	departmentsController.Routes(s)

	s.Run()
}
