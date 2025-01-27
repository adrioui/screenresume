package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
)

type JobRolesService interface {
	GetAllJobRoles(ctx context.Context) ([]models.JobRoles, error)
	CreateJobRoles(ctx context.Context, input models.JobRolesCreate) (models.JobRoles, error)
	GetJobRoles(ctx context.Context, id string) (models.JobRoles, error)
	UpdateJobRoles(ctx context.Context, id string, input models.JobRolesUpdate) (models.JobRoles, error)
	DeleteJobRoles(ctx context.Context, id string) (any, error)
}

type JobRoleServiceImpl struct {
	store db.Store
}

func NewJobRoleService(store db.Store) *JobRoleServiceImpl {
	return &JobRoleServiceImpl{store: store}
}

// Type conversion helpers
func toJobRolesDTO(dbJobRole repositories.JobRole) models.JobRoles {
	return models.JobRoles{
		ID:           dbJobRole.ID.String(),
		Title:        dbJobRole.Title,
		DepartmentID: dbJobRole.DepartmentID.String(),
		Level:        string(dbJobRole.Level),
		SalaryRange:  dbJobRole.SalaryRange,
		Location:     dbJobRole.Location,
		IsActive:     dbJobRole.IsActive,
	}
}

func toCreateJobRoleParams(input models.JobRolesCreate) (repositories.CreateJobRoleParams, error) {
	// Convert DepartmentID from string to uuid.UUID
	departmentID, err := uuid.Parse(input.DepartmentID)
	if err != nil {
		return repositories.CreateJobRoleParams{}, fmt.Errorf("invalid DepartmentID: %w", err)
	}

	// Convert Level from string to ExperienceLevel
	level := repositories.ExperienceLevel(input.Level) // Assuming ExperienceLevel is a type alias for string

	return repositories.CreateJobRoleParams{
		Title:        input.Title,
		DepartmentID: departmentID,
		Level:        level,
		SalaryRange:  input.SalaryRange,
		Location:     input.Location,
		IsActive:     input.IsActive,
	}, nil
}

func toUpdateJobRoleParams(id uuid.UUID, input models.JobRolesUpdate) (repositories.UpdateJobRoleParams, error) {
	// Convert DepartmentID from string to uuid.UUID
	departmentID, err := uuid.Parse(input.DepartmentID)
	if err != nil {
		return repositories.UpdateJobRoleParams{}, fmt.Errorf("invalid DepartmentID: %w", err)
	}

	// Convert Level from string to ExperienceLevel
	level := repositories.ExperienceLevel(input.Level) // Assuming ExperienceLevel is a type alias for string

	return repositories.UpdateJobRoleParams{
		ID:           id,
		Title:        input.Title,
		DepartmentID: departmentID,
		Level:        level,
		SalaryRange:  input.SalaryRange,
		Location:     input.Location,
		IsActive:     input.IsActive,
	}, nil
}

// Service method implementations
func (s *JobRoleServiceImpl) GetAllJobRoles(ctx context.Context) ([]models.JobRoles, error) {
	dbJobRoles, err := s.store.ListJobRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobRoles: %w", err)
	}

	jobRoles := make([]models.JobRoles, len(dbJobRoles))
	for i, f := range dbJobRoles {
		jobRoles[i] = toJobRolesDTO(f)
	}
	return jobRoles, nil
}

func (s *JobRoleServiceImpl) CreateJobRoles(ctx context.Context, input models.JobRolesCreate) (models.JobRoles, error) {
	params, err := toCreateJobRoleParams(input)
	if err != nil {
		return models.JobRoles{}, fmt.Errorf("failed to create jobRole: %w", err)
	}

	dbJobRole, err := s.store.CreateJobRole(ctx, params)
	if err != nil {
		return models.JobRoles{}, fmt.Errorf("failed to create jobRole: %w", err)
	}

	fmt.Println(params)
	return toJobRolesDTO(dbJobRole), nil
}

func (s *JobRoleServiceImpl) GetJobRoles(ctx context.Context, id string) (models.JobRoles, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.JobRoles{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbJobRole, err := s.store.GetJobRole(ctx, uuidID)
	if err != nil {
		return models.JobRoles{}, fmt.Errorf("failed to get jobRole: %w", err)
	}

	return toJobRolesDTO(dbJobRole), nil
}

func (s *JobRoleServiceImpl) UpdateJobRoles(ctx context.Context, id string, input models.JobRolesUpdate) (models.JobRoles, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.JobRoles{}, fmt.Errorf("invalid ID format: %w", err)
	}

	params, err := toUpdateJobRoleParams(uuidID, input)

	if err := s.store.UpdateJobRole(ctx, params); err != nil {
		return models.JobRoles{}, fmt.Errorf("failed to update jobRole: %w", err)
	}

	// Fetch updated entity
	return s.GetJobRoles(ctx, id)
}

func (s *JobRoleServiceImpl) DeleteJobRoles(ctx context.Context, id string) (any, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.store.DeleteJobRole(ctx, uuidID); err != nil {
		return nil, fmt.Errorf("failed to delete jobRole: %w", err)
	}

	return nil, nil
}
