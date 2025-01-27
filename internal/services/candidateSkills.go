package services

import (
	"context"
	"database/sql"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type CandidateSkillsService interface {
	GetAllCandidateSkills(ctx context.Context) ([]models.CandidateSkills, error)
	CreateCandidateSkills(ctx context.Context, input models.CandidateSkillsCreate) (models.CandidateSkills, error)
	GetCandidateSkills(ctx context.Context, candidateID string, skillID string) (models.CandidateSkills, error)
	UpdateCandidateSkills(ctx context.Context, candidateID string, skillID string, input models.CandidateSkillsUpdate) (models.CandidateSkills, error)
	DeleteCandidateSkills(ctx context.Context, candidateID string, skillID string) (any, error)
}

type CandidateSkillServiceImpl struct {
	store db.Store
}

func NewCandidateSkillService(store db.Store) *CandidateSkillServiceImpl {
	return &CandidateSkillServiceImpl{store: store}
}

func nullTimeToString(nt sql.NullTime) (string, error) {
	if !nt.Valid {
		return "", fmt.Errorf("null time is not valid")
	}
	return nt.Time.Format("2006-01-02"), nil
}

func stringToNullTime(timeStr string) (sql.NullTime, error) {
	if timeStr == "" {
		// Return a NullTime with Valid = false for empty strings
		return sql.NullTime{Valid: false}, nil
	}

	// Parse the time string into a time.Time object
	// Adjust the layout to match your input format
	parsedTime, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		return sql.NullTime{}, err
	}

	// Return a valid NullTime
	return sql.NullTime{
		Time:  parsedTime,
		Valid: true,
	}, nil
}

// Type conversion helpers
func toCandidateSkillsDTO(dbCandidateSkill repositories.CandidateSkill) (models.CandidateSkills, error) {
	var yearsExperience float64
	if err := dbCandidateSkill.YearsExperience.AssignTo(&yearsExperience); err != nil {
		return models.CandidateSkills{}, err
	}

	// Convert sql.NullTime to string
	lastUsed, err := nullTimeToString(dbCandidateSkill.LastUsed)
	if err != nil {
		return models.CandidateSkills{}, err
	}

	return models.CandidateSkills{
		CandidateID:     dbCandidateSkill.CandidateID.String(),
		SkillID:         dbCandidateSkill.SkillID.String(),
		YearsExperience: yearsExperience,
		LastUsed:        lastUsed,
	}, nil
}

func toCandidateSkillsCreateParams(input models.CandidateSkillsCreate) (repositories.CreateCandidateSkillsParams, error) {
	candidateID, err := uuid.Parse(input.CandidateID)
	if err != nil {
		return repositories.CreateCandidateSkillsParams{}, err
	}

	skillID, err := uuid.Parse(input.SkillID)
	if err != nil {
		return repositories.CreateCandidateSkillsParams{}, err
	}

	var yearsExperience pgtype.Numeric
	if err := yearsExperience.Scan(strconv.FormatFloat(input.YearsExperience, 'f', -1, 64)); err != nil {
		return repositories.CreateCandidateSkillsParams{}, err
	}

	lastUsed, err := stringToNullTime(input.LastUsed)
	if err != nil {
		return repositories.CreateCandidateSkillsParams{}, err
	}

	return repositories.CreateCandidateSkillsParams{
		CandidateID:     candidateID,
		SkillID:         skillID,
		YearsExperience: yearsExperience,
		LastUsed:        lastUsed,
	}, nil
}

// Service method implementations
func (s *CandidateSkillServiceImpl) GetAllCandidateSkills(ctx context.Context) ([]models.CandidateSkills, error) {
	dbCandidateSkills, err := s.store.ListCandidateSkills(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list candidateSkills: %w", err)
	}

	candidateSkills := make([]models.CandidateSkills, len(dbCandidateSkills))
	for i, f := range dbCandidateSkills {
		if candidateSkill, err := toCandidateSkillsDTO(f); err != nil {
			return nil, err
		} else {
			candidateSkills[i] = candidateSkill
		}
	}
	return candidateSkills, nil
}

func (s *CandidateSkillServiceImpl) CreateCandidateSkills(ctx context.Context, input models.CandidateSkillsCreate) (models.CandidateSkills, error) {
	params, err := toCandidateSkillsCreateParams(input)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("failed to create candidateSkill: %w", err)
	}

	dbCandidateSkill, err := s.store.CreateCandidateSkills(ctx, params)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("failed to create candidateSkill: %w", err)
	}

	candidateSkillsDTO, err := toCandidateSkillsDTO(dbCandidateSkill)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("failed to create candidateSkill: %w", err)
	}

	return candidateSkillsDTO, nil
}

func (s *CandidateSkillServiceImpl) GetCandidateSkills(ctx context.Context, candidateID string, skillID string) (models.CandidateSkills, error) {
	candidateIDUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("invalid jobRoleID format: %w", err)
	}
	skillIDUUID, err := uuid.Parse(skillID)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("invalid skillID format: %w", err)
	}

	dbCandidateSkill, err := s.store.GetCandidateSkills(ctx, repositories.GetCandidateSkillsParams{
		CandidateID: candidateIDUUID,
		SkillID:     skillIDUUID,
	})
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("failed to get candidateSkill: %w", err)
	}

	candidateSkillsDTO, err := toCandidateSkillsDTO(dbCandidateSkill)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("failed to create candidateSkill: %w", err)
	}

	return candidateSkillsDTO, nil
}

func (s *CandidateSkillServiceImpl) UpdateCandidateSkills(ctx context.Context, candidateID string, skillID string, input models.CandidateSkillsUpdate) (models.CandidateSkills, error) {
	candidateIDUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("invalid jobRoleID format: %w", err)
	}
	skillIDUUID, err := uuid.Parse(skillID)
	if err != nil {
		return models.CandidateSkills{}, fmt.Errorf("invalid skillID format: %w", err)
	}
	var yearsExperience pgtype.Numeric
	if err := yearsExperience.Scan(input.YearsExperience); err != nil {
		return models.CandidateSkills{}, err
	}

	lastUsed, err := stringToNullTime(input.LastUsed)
	if err != nil {
		return models.CandidateSkills{}, err
	}

	params := repositories.UpdateCandidateSkillsParams{
		CandidateID:     candidateIDUUID,
		SkillID:         skillIDUUID,
		YearsExperience: yearsExperience,
		LastUsed:        lastUsed,
	}

	if err := s.store.UpdateCandidateSkills(ctx, params); err != nil {
		return models.CandidateSkills{}, fmt.Errorf("failed to update candidateSkill: %w", err)
	}

	// Fetch updated entity
	return s.GetCandidateSkills(ctx, candidateID, skillID)
}

func (s *CandidateSkillServiceImpl) DeleteCandidateSkills(ctx context.Context, candidateID string, skillID string) (any, error) {
	candidateIDUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidateID format: %w", err)
	}
	skillIDUUID, err := uuid.Parse(skillID)
	if err != nil {
		return nil, fmt.Errorf("invalid skillID format: %w", err)
	}

	if err := s.store.DeleteCandidateSkills(ctx, repositories.DeleteCandidateSkillsParams{
		CandidateID: candidateIDUUID,
		SkillID:     skillIDUUID,
	}); err != nil {
		return nil, fmt.Errorf("failed to delete candidateSkill: %w", err)
	}

	return nil, nil
}
