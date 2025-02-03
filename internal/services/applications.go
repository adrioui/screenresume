package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type ApplicationService interface {
	GetAllApplication(ctx context.Context) ([]models.Application, error)
	CreateApplication(ctx context.Context, input models.ApplicationCreate) (models.Application, error)
	GetApplication(ctx context.Context, id string) (models.Application, error)
	UpdateApplication(ctx context.Context, id string, input models.ApplicationUpdate) (models.Application, error)
	DeleteApplication(ctx context.Context, id string) (any, error)
}

type ApplicationServiceImpl struct {
	store db.Store
}

func NewApplicationService(store db.Store) *ApplicationServiceImpl {
	return &ApplicationServiceImpl{store: store}
}

// Type conversion helpers
func toApplicationDTO(dbApplication repositories.Application) (models.Application, error) {
	var score float64
	if err := dbApplication.Score.AssignTo(&score); err != nil {
		return models.Application{}, err
	}

	return models.Application{
		ID:          dbApplication.ID.String(),
		CandidateID: dbApplication.CandidateID.String(),
		JobRoleID:   dbApplication.JobRoleID.String(),
		Stage:       dbApplication.Stage,
		Score:       score,
		AppliedAt:   dbApplication.AppliedAt,
		LastUpdated: dbApplication.LastUpdated,
	}, nil
}

func toApplicationCreateParams(input models.ApplicationCreate) (repositories.CreateApplicationParams, error) {
	var score pgtype.Numeric
	if err := score.Scan(strconv.FormatFloat(input.Score, 'f', -1, 64)); err != nil {
		return repositories.CreateApplicationParams{}, err
	}

	return repositories.CreateApplicationParams{
		CandidateID: uuid.MustParse(input.CandidateID),
		JobRoleID:   uuid.MustParse(input.JobRoleID),
		Stage:       input.Stage,
		Score:       score,
	}, nil
}

// Service method implementations
func (s *ApplicationServiceImpl) GetAllApplication(ctx context.Context) ([]models.Application, error) {
	dbApplication, err := s.store.ListApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]models.Application, len(dbApplication))
	for i, f := range dbApplication {
		dtoApplication, err := toApplicationDTO(f)
		if err != nil {
			return nil, err
		}
		files[i] = dtoApplication
	}
	return files, nil
}

func (s *ApplicationServiceImpl) CreateApplication(ctx context.Context, input models.ApplicationCreate) (models.Application, error) {
	params, err := toApplicationCreateParams(input)

	dbApplication, err := s.store.CreateApplication(ctx, params)
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to create file: %w", err)
	}
	applicationDTO, err := toApplicationDTO(dbApplication)
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to create file: %w", err)
	}
	return applicationDTO, nil
}

func (s *ApplicationServiceImpl) GetApplication(ctx context.Context, id string) (models.Application, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Application{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbApplication, err := s.store.GetApplication(ctx, uuidID)
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to get file: %w", err)
	}

	applicationDTO, err := toApplicationDTO(dbApplication)
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to create file: %w", err)
	}
	return applicationDTO, nil
}

func (s *ApplicationServiceImpl) UpdateApplication(ctx context.Context, id string, input models.ApplicationUpdate) (models.Application, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Application{}, fmt.Errorf("invalid ID format: %w", err)
	}

	var score pgtype.Numeric
	if err := score.Scan(strconv.FormatFloat(input.Score, 'f', -1, 64)); err != nil {
		return models.Application{}, err
	}

	params := repositories.UpdateApplicationParams{
		ID:          uuidID,
		CandidateID: uuid.MustParse(input.CandidateID),
		JobRoleID:   uuid.MustParse(input.JobRoleID),
		Stage:       input.Stage,
		Score:       score,
	}

	if err := s.store.UpdateApplication(ctx, params); err != nil {
		return models.Application{}, fmt.Errorf("failed to update file: %w", err)
	}

	// Fetch updated entity
	return s.GetApplication(ctx, id)
}

func (s *ApplicationServiceImpl) DeleteApplication(ctx context.Context, id string) (any, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.store.DeleteApplication(ctx, uuidID); err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return nil, nil
}
