package services

import (
	"context"
	"fmt"
	"log"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	pb "screenresume/internal/repositories"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ScreenResumeService interface {
	ProcessScreening(
		request *pb.ScreenResumeRequest,
		input models.ProcessScreeningCreate,
		fileInput models.FilesCreate,
	) (models.ScreenResume, error)
}

type ScreenResumeServiceImpl struct {
	client                   pb.ResumeScreenerClient
	skillsService            SkillsService
	screeningResultsService  ScreeningResultsService
	screeningCriteriaService ScreeningCriteriaService
	filesService             FilesService
}

func NewScreenResumeService(
	conn *grpc.ClientConn,
	skillsService SkillsService,
	screeningResultsService ScreeningResultsService,
	screeningCriteriaService ScreeningCriteriaService,
	filesService FilesService,
) *ScreenResumeServiceImpl {
	client := repositories.NewResumeScreenerClient(conn)
	return &ScreenResumeServiceImpl{
		client:                   client,
		skillsService:            skillsService,
		screeningResultsService:  screeningResultsService,
		screeningCriteriaService: screeningCriteriaService,
		filesService:             filesService,
	}
}

func (s *ScreenResumeServiceImpl) ProcessScreening(
	request *pb.ScreenResumeRequest,
	input models.ProcessScreeningCreate,
	fileInput models.FilesCreate,
) (models.ScreenResume, error) {
	// Call gRPC service
	response, err := s.client.ScreenResume(context.Background(), request)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("error: %v", st)
		return models.ScreenResume{}, err
	}

	// Ensure response length matches criteria length
	if len(response.CriteriaDecisions) != len(input.Criteria) {
		log.Printf("Warning: Number of criteria decisions (%d) does not match input criteria (%d)", len(response.CriteriaDecisions), len(input.Criteria))
	}

	// Convert response to model
	criteriaDecisions := make([]*models.CriteriaDecision, len(response.CriteriaDecisions))
	for i, decision := range response.CriteriaDecisions {
		criteriaDecisions[i] = &models.CriteriaDecision{
			Reasoning: decision.Reasoning,
			Decision:  decision.Decision,
		}
	}

	// Create a screening result
	var screeningResultCreate models.ScreeningResultsCreate
	screeningResultCreate.ApplicationID = input.ApplicationID
	screeningResultCreate.ModelVersion = "gpt-4o mini"
	screeningResultCreate.RawResponse.CriteriaDecisions = criteriaDecisions
	screeningResultCreate.RawResponse.OverallDecision = response.OverallDecision
	screeningResultCreate.RawResponse.OverallReasoning = response.OverallReasoning
	screeningResultCreate.RawResponse.ResumeName = response.ResumeName

	screeningResult, err := s.screeningResultsService.CreateScreeningResults(context.Background(), screeningResultCreate)
	if err != nil {
		return models.ScreenResume{}, fmt.Errorf("failed to create screening result: %w", err)
	}

	log.Printf("Created screening result: %+v", screeningResult)

	var screeningCriteriaCreate models.ScreeningCriteriaCreate
	screeningCriteriaCreate.ScreeningResultsID = screeningResult.ID
	screeningCriteriaCreate.Decision = response.OverallDecision
	screeningCriteriaCreate.Reasoning = response.OverallReasoning

	for i, criteria := range input.Criteria {
		skill, err := s.skillsService.GetSkillByName(context.Background(), criteria)
		if err != nil {
			log.Printf("Error fetching skill for criteria %s: %v", criteria, err)
			continue // Skip errors instead of failing the entire process
		}

		log.Printf("Fetched skill: %+v", skill)

		// Get the corresponding decision (if available)
		if i >= len(response.CriteriaDecisions) {
			log.Printf("Skipping criteria %s as there is no matching decision", criteria)
			continue
		}

		decision := response.CriteriaDecisions[i]

		if decision.Decision {
			// If true, add to MatchedSkills
			screeningCriteriaCreate.MatchedSkills = append(screeningCriteriaCreate.MatchedSkills, skill.ID)
		} else {
			// If false, add to MissingSkills
			screeningCriteriaCreate.MissingSkills = append(screeningCriteriaCreate.MissingSkills, skill.ID)
		}
	}
	screeningCriteria, err := s.screeningCriteriaService.CreateScreeningCriteria(context.Background(), screeningCriteriaCreate)
	if err != nil {
		return models.ScreenResume{}, fmt.Errorf("failed to create screening criteria: %w", err)
	}

	fmt.Printf("Created screening criteria: %+v", screeningCriteria)

	// Create file record
	if _, err := s.filesService.CreateFiles(context.Background(), fileInput); err != nil {
		return models.ScreenResume{}, fmt.Errorf("failed to create file: %w", err)
	}

	result := models.ScreenResume{
		CriteriaDecisions: criteriaDecisions,
		OverallReasoning:  response.OverallReasoning,
		OverallDecision:   response.OverallDecision,
		ResumeName:        response.ResumeName,
	}

	return result, nil
}
