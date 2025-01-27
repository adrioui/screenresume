package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	// Handle Location conversion (pgtype.Text to *string)
	var location *string
	if dbJobRole.Location.Valid {
		location = &dbJobRole.Location.String
	}

	var salaryRange *string
	if dbJobRole.SalaryRange.Valid {
		salaryRange = &dbJobRole.SalaryRange.String
	}

	return models.JobRoles{
		ID:           dbJobRole.ID.String(),
		Title:        dbJobRole.Title,
		DepartmentID: dbJobRole.DepartmentID.String(),
		Level:        string(dbJobRole.Level),
		SalaryRange:  salaryRange,
		Location:     location,
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

	// Convert Location and SalaryRange from *string to pgtype.Text
	var location pgtype.Text
	if input.Location != nil {
		location = pgtype.Text{String: *input.Location, Valid: true}
	} else {
		location = pgtype.Text{Valid: false}
	}

	var salaryRange pgtype.Text
	if input.SalaryRange != nil {
		salaryRange = pgtype.Text{String: *input.SalaryRange, Valid: true}
	} else {
		salaryRange = pgtype.Text{Valid: false}
	}

	return repositories.CreateJobRoleParams{
		Title:        input.Title,
		DepartmentID: departmentID,
		Level:        level,
		SalaryRange:  salaryRange,
		Location:     location,
		IsActive:     input.IsActive,
	}, nil
}

func toUpdateJobRoleParams(id uuid.UUID, input models.JobRolesUpdate) (repositories.UpdateJobRoleParams, error) {
	params := repositories.UpdateJobRoleParams{
		ID: id,
	}

	// Handle nullable fields
	if input.Title != nil {
		params.Title = *input.Title
	}

	if input.DepartmentID != nil {
		departmentID, err := uuid.Parse(*input.DepartmentID)
		if err != nil {
			return repositories.UpdateJobRoleParams{}, fmt.Errorf("invalid department ID: %w", err)
		}
		params.DepartmentID = departmentID
	}

	if input.Level != nil {
		params.Level = repositories.ExperienceLevel(*input.Level)
	}

	// Handle pgtype.Text conversion
	if input.Location != nil {
		params.Location = pgtype.Text{
			String: *input.Location,
			Valid:  true,
		}
	} else {
		params.Location = pgtype.Text{Valid: false}
	}

	if input.SalaryRange != nil {
		params.SalaryRange = pgtype.Text{
			String: *input.SalaryRange,
			Valid:  true,
		}
	} else {
		params.SalaryRange = pgtype.Text{Valid: false}
	}

	if input.IsActive != nil {
		params.IsActive = *input.IsActive
	}

	return params, nil
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
