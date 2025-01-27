package services

import (
	"context"
	"fmt"
	"math/big"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobRoleRequirementsService interface {
	GetAllJobRoleRequirements(ctx context.Context) ([]models.JobRoleRequirements, error)
	CreateJobRoleRequirements(ctx context.Context, input models.JobRoleRequirementsCreate) (models.JobRoleRequirements, error)
	GetJobRoleRequirements(ctx context.Context, jobRoleID string, skillID string) (models.JobRoleRequirements, error)
	UpdateJobRoleRequirements(ctx context.Context, jobRoleID string, skillID string, input models.JobRoleRequirementsUpdate) (models.JobRoleRequirements, error)
	DeleteJobRoleRequirements(ctx context.Context, jobRoleID string, skillID string) (any, error)
}

type JobRoleRequirementServiceImpl struct {
	store db.Store
}

func NewJobRoleRequirementService(store db.Store) *JobRoleRequirementServiceImpl {
	return &JobRoleRequirementServiceImpl{store: store}
}

// Convert pgtype.Numeric to float64
func numericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid {
		return 0, fmt.Errorf("numeric value is not valid")
	}

	// Convert the numeric value to a big.Float
	bigFloat := new(big.Float).SetInt(n.Int)

	// Apply the exponent to adjust the decimal places
	if n.Exp != 0 {
		exp := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil))
		bigFloat.Quo(bigFloat, exp)
	}

	// Convert big.Float to float64
	floatVal, _ := bigFloat.Float64()
	return floatVal, nil
}

// Type conversion helpers
func toJobRoleRequirementsDTO(dbJobRoleRequirement repositories.JobRoleRequirement) (models.JobRoleRequirements, error) {
	minExperienceYears, err := numericToFloat64(dbJobRoleRequirement.MinExperienceYears)
	if err != nil {
		return models.JobRoleRequirements{}, err
	}

	return models.JobRoleRequirements{
		JobRoleID:          dbJobRoleRequirement.JobRoleID.String(),
		SkillID:            dbJobRoleRequirement.SkillID.String(),
		Required:           dbJobRoleRequirement.Required,
		MinExperienceYears: minExperienceYears,
		Importance:         int(dbJobRoleRequirement.Importance),
	}, nil
}

func toJobRoleRequirementsCreateParams(input models.JobRoleRequirementsCreate) (repositories.CreateJobRoleRequirementParams, error) {
	jobRoleID, err := uuid.Parse(input.JobRoleID)
	if err != nil {
		return repositories.CreateJobRoleRequirementParams{}, err
	}

	skillID, err := uuid.Parse(input.SkillID)
	if err != nil {
		return repositories.CreateJobRoleRequirementParams{}, err
	}

	var minExperienceYears pgtype.Numeric
	if err := minExperienceYears.Scan(strconv.FormatFloat(input.MinExperienceYears, 'f', -1, 64)); err != nil {
		return repositories.CreateJobRoleRequirementParams{}, err
	}

	return repositories.CreateJobRoleRequirementParams{
		JobRoleID:          jobRoleID,
		SkillID:            skillID,
		Required:           input.Required,
		MinExperienceYears: minExperienceYears,
		Importance:         int32(input.Importance),
	}, nil
}

// Service method implementations
func (s *JobRoleRequirementServiceImpl) GetAllJobRoleRequirements(ctx context.Context) ([]models.JobRoleRequirements, error) {
	dbJobRoleRequirements, err := s.store.ListJobRoleRequirements(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobRoleRequirements: %w", err)
	}

	jobRoleRequirements := make([]models.JobRoleRequirements, len(dbJobRoleRequirements))
	for i, f := range dbJobRoleRequirements {
		if jobRoleRequirement, err := toJobRoleRequirementsDTO(f); err != nil {
			return nil, err
		} else {
			jobRoleRequirements[i] = jobRoleRequirement
		}
	}
	return jobRoleRequirements, nil
}

func (s *JobRoleRequirementServiceImpl) CreateJobRoleRequirements(ctx context.Context, input models.JobRoleRequirementsCreate) (models.JobRoleRequirements, error) {
	params, err := toJobRoleRequirementsCreateParams(input)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("failed to create jobRoleRequirement: %w", err)
	}

	dbJobRoleRequirement, err := s.store.CreateJobRoleRequirement(ctx, params)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("failed to create jobRoleRequirement: %w", err)
	}

	jobRoleRequirementsDTO, err := toJobRoleRequirementsDTO(dbJobRoleRequirement)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("failed to create jobRoleRequirement: %w", err)
	}

	return jobRoleRequirementsDTO, nil
}

func (s *JobRoleRequirementServiceImpl) GetJobRoleRequirements(ctx context.Context, jobRoleID string, skillID string) (models.JobRoleRequirements, error) {
	jobRoleIDUUID, err := uuid.Parse(jobRoleID)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("invalid jobRoleID format: %w", err)
	}
	skillIDUUID, err := uuid.Parse(skillID)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("invalid skillID format: %w", err)
	}

	dbJobRoleRequirement, err := s.store.GetJobRoleRequirement(ctx, repositories.GetJobRoleRequirementParams{
		JobRoleID: jobRoleIDUUID,
		SkillID:   skillIDUUID,
	})
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("failed to get jobRoleRequirement: %w", err)
	}

	jobRoleRequirementsDTO, err := toJobRoleRequirementsDTO(dbJobRoleRequirement)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("failed to create jobRoleRequirement: %w", err)
	}

	return jobRoleRequirementsDTO, nil
}

func (s *JobRoleRequirementServiceImpl) UpdateJobRoleRequirements(ctx context.Context, jobRoleID string, skillID string, input models.JobRoleRequirementsUpdate) (models.JobRoleRequirements, error) {
	jobRoleIDUUID, err := uuid.Parse(jobRoleID)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("invalid jobRoleID format: %w", err)
	}
	skillIDUUID, err := uuid.Parse(skillID)
	if err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("invalid skillID format: %w", err)
	}
	var minExperienceYears pgtype.Numeric
	if err := minExperienceYears.Scan(input.MinExperienceYears); err != nil {
		return models.JobRoleRequirements{}, err
	}
	params := repositories.UpdateJobRoleRequirementParams{
		JobRoleID:          jobRoleIDUUID,
		SkillID:            skillIDUUID,
		Required:           *input.Required,
		MinExperienceYears: minExperienceYears,
		Importance:         int32(*input.Importance),
	}

	if err := s.store.UpdateJobRoleRequirement(ctx, params); err != nil {
		return models.JobRoleRequirements{}, fmt.Errorf("failed to update jobRoleRequirement: %w", err)
	}

	// Fetch updated entity
	return s.GetJobRoleRequirements(ctx, jobRoleID, skillID)
}

func (s *JobRoleRequirementServiceImpl) DeleteJobRoleRequirements(ctx context.Context, jobRoleID string, skillID string) (any, error) {
	jobRoleIDUUID, err := uuid.Parse(jobRoleID)
	if err != nil {
		return nil, fmt.Errorf("invalid jobRoleID format: %w", err)
	}
	skillIDUUID, err := uuid.Parse(skillID)
	if err != nil {
		return nil, fmt.Errorf("invalid skillID format: %w", err)
	}

	if err := s.store.DeleteJobRoleRequirement(ctx, repositories.DeleteJobRoleRequirementParams{
		JobRoleID: jobRoleIDUUID,
		SkillID:   skillIDUUID,
	}); err != nil {
		return nil, fmt.Errorf("failed to delete jobRoleRequirement: %w", err)
	}

	return nil, nil
}
