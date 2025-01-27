package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
)

type SkillsService interface {
	GetAllSkills(ctx context.Context) ([]models.Skills, error)
	CreateSkills(ctx context.Context, input models.SkillsCreate) (models.Skills, error)
	GetSkills(ctx context.Context, id string) (models.Skills, error)
	UpdateSkills(ctx context.Context, id string, input models.SkillsUpdate) (models.Skills, error)
	DeleteSkills(ctx context.Context, id string) (any, error)
}

type SkillServiceImpl struct {
	store db.Store
}

func NewSkillService(store db.Store) *SkillServiceImpl {
	return &SkillServiceImpl{store: store}
}

// Type conversion helpers
func toSkillsDTO(dbSkill repositories.Skill) models.Skills {
	return models.Skills{
		ID:       dbSkill.ID.String(),
		Name:     dbSkill.Name,
		Category: dbSkill.Category,
	}
}

func toSkillsCreateParams(input models.SkillsCreate) repositories.CreateSkillParams {
	return repositories.CreateSkillParams{
		Name:     input.Name,
		Category: input.Category,
	}
}

// Service method implementations
func (s *SkillServiceImpl) GetAllSkills(ctx context.Context) ([]models.Skills, error) {
	dbSkills, err := s.store.ListSkills(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}

	skills := make([]models.Skills, len(dbSkills))
	for i, f := range dbSkills {
		skills[i] = toSkillsDTO(f)
	}
	return skills, nil
}

func (s *SkillServiceImpl) CreateSkills(ctx context.Context, input models.SkillsCreate) (models.Skills, error) {
	params := toSkillsCreateParams(input)

	dbSkill, err := s.store.CreateSkill(ctx, params)
	if err != nil {
		return models.Skills{}, fmt.Errorf("failed to create skill: %w", err)
	}

	return toSkillsDTO(dbSkill), nil
}

func (s *SkillServiceImpl) GetSkills(ctx context.Context, id string) (models.Skills, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Skills{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbSkill, err := s.store.GetSkill(ctx, uuidID)
	if err != nil {
		return models.Skills{}, fmt.Errorf("failed to get skill: %w", err)
	}

	return toSkillsDTO(dbSkill), nil
}

func (s *SkillServiceImpl) UpdateSkills(ctx context.Context, id string, input models.SkillsUpdate) (models.Skills, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Skills{}, fmt.Errorf("invalid ID format: %w", err)
	}

	params := repositories.UpdateSkillParams{
		ID:       uuidID,
		Name:     input.Name,
		Category: input.Category,
	}

	if err := s.store.UpdateSkill(ctx, params); err != nil {
		return models.Skills{}, fmt.Errorf("failed to update skill: %w", err)
	}

	// Fetch updated entity
	return s.GetSkills(ctx, id)
}

func (s *SkillServiceImpl) DeleteSkills(ctx context.Context, id string) (any, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.store.DeleteSkill(ctx, uuidID); err != nil {
		return nil, fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil, nil
}
