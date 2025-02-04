package main

import (
	"log"
	"screenresume/internal/controller"
	"screenresume/internal/services"
	"screenresume/pkg/db"

	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create the client
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Initialize the database connection pool
	dbPool, err := db.CreateDBConnectionPool()
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v", err)
	}
	defer dbPool.Close()

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

	jobRoleRequirementsService := services.NewJobRoleRequirementService(store)
	jobRoleRequirementsController := controller.JobRoleRequirementsResources{JobRoleRequirementsService: jobRoleRequirementsService}

	candidateSkillsService := services.NewCandidateSkillService(store)
	candidateSkillsController := controller.CandidateSkillsResources{CandidateSkillsService: candidateSkillsService}

	screeningResultsService := services.NewScreeningResultService(store)
	screeningCriteriaService := services.NewScreeningCriteriaService(store)

	screenResumeService := services.NewScreenResumeService(
		conn,
		skillsService,
		screeningResultsService,
		screeningCriteriaService,
		filesService,
	)
	screenResumeController := controller.ScreenResumeResources{ScreenResumeService: screenResumeService}

	applicationService := services.NewApplicationService(store)
	applicationController := controller.ApplicationResources{ApplicationService: applicationService}

	s := fuego.NewServer(
		fuego.WithCorsMiddleware(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}).Handler),
	)

	filesController.Routes(s)
	jobRolesController.Routes(s)
	departmentsController.Routes(s)
	candidatesController.Routes(s)
	skillsController.Routes(s)
	jobRoleRequirementsController.Routes(s)
	candidateSkillsController.Routes(s)
	screenResumeController.Routes(s)
	applicationController.Routes(s)

	s.Run()
}
