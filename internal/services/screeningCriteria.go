package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
)

type ScreeningCriteriaService interface {
	GetAllScreeningCriteria(ctx context.Context) ([]models.ScreeningCriteria, error)
	CreateScreeningCriteria(ctx context.Context, input models.ScreeningCriteriaCreate) (models.ScreeningCriteria, error)
	GetScreeningCriteria(ctx context.Context, id string) (models.ScreeningCriteria, error)
}

type ScreeningCriteriaServiceImpl struct {
	store db.Store
}

func NewScreeningCriteriaService(store db.Store) *ScreeningCriteriaServiceImpl {
	return &ScreeningCriteriaServiceImpl{store: store}
}

// Type conversion helpers
func toScreeningCriteriaDTO(dbScreeningCriteria repositories.ScreeningCriterium) models.ScreeningCriteria {
	matchedSkills := make([]string, len(dbScreeningCriteria.MatchedSkills))
	for i, f := range dbScreeningCriteria.MatchedSkills {
		matchedSkills[i] = f.String()
	}
	missingSkills := make([]string, len(dbScreeningCriteria.MissingSkills))
	for i, f := range dbScreeningCriteria.MissingSkills {
		missingSkills[i] = f.String()
	}

	return models.ScreeningCriteria{
		ID:                 dbScreeningCriteria.ID.String(),
		ScreeningResultsID: dbScreeningCriteria.ScreeningResultID.String(),
		Decision:           dbScreeningCriteria.Decision,
		Reasoning:          dbScreeningCriteria.Reasoning,
		MatchedSkills:      matchedSkills,
		MissingSkills:      missingSkills,
	}
}

func toScreeningCriteriaCreateParams(input models.ScreeningCriteriaCreate) repositories.CreateScreeningCriteriaParams {
	matchedSkills := make([]uuid.UUID, len(input.MatchedSkills))
	for i, f := range input.MatchedSkills {
		matchedSkills[i] = uuid.MustParse(f)
	}
	missingSkills := make([]uuid.UUID, len(input.MissingSkills))
	for i, f := range input.MissingSkills {
		missingSkills[i] = uuid.MustParse(f)
	}
	return repositories.CreateScreeningCriteriaParams{
		ScreeningResultID: uuid.MustParse(input.ScreeningResultsID),
		Decision:          input.Decision,
		Reasoning:         input.Reasoning,
		MatchedSkills:     matchedSkills,
		MissingSkills:     missingSkills,
	}
}

// Service method implementations
func (s *ScreeningCriteriaServiceImpl) GetAllScreeningCriteria(ctx context.Context) ([]models.ScreeningCriteria, error) {
	dbScreeningCriteria, err := s.store.ListScreeningCriteria(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list screeningCriteria: %w", err)
	}

	screeningCriteria := make([]models.ScreeningCriteria, len(dbScreeningCriteria))
	for i, f := range dbScreeningCriteria {
		screeningCriteria[i] = toScreeningCriteriaDTO(f)
	}
	return screeningCriteria, nil
}

func (s *ScreeningCriteriaServiceImpl) CreateScreeningCriteria(ctx context.Context, input models.ScreeningCriteriaCreate) (models.ScreeningCriteria, error) {
	params := toScreeningCriteriaCreateParams(input)

	dbScreeningCriteria, err := s.store.CreateScreeningCriteria(ctx, params)
	if err != nil {
		return models.ScreeningCriteria{}, fmt.Errorf("failed to create screeningCriteria: %w", err)
	}

	return toScreeningCriteriaDTO(dbScreeningCriteria), nil
}

func (s *ScreeningCriteriaServiceImpl) GetScreeningCriteria(ctx context.Context, id string) (models.ScreeningCriteria, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.ScreeningCriteria{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbScreeningCriteria, err := s.store.GetScreeningCriteria(ctx, uuidID)
	if err != nil {
		return models.ScreeningCriteria{}, fmt.Errorf("failed to get screeningCriteria: %w", err)
	}

	return toScreeningCriteriaDTO(dbScreeningCriteria), nil
}

func (s *ScreeningCriteriaServiceImpl) CreateScreeningCriteriaWithTx(ctx context.Context, q repositories.Querier, input models.ScreeningCriteriaCreate) (models.ScreeningCriteria, error) {
	params := toScreeningCriteriaCreateParams(input)

	dbScreeningCriteria, err := q.CreateScreeningCriteria(ctx, params)
	if err != nil {
		return models.ScreeningCriteria{}, fmt.Errorf("failed to create screeningCriteria: %w", err)
	}

	screeningCriteriaDTO := toScreeningCriteriaDTO(dbScreeningCriteria)
	return screeningCriteriaDTO, nil
}
