package services

import (
	"context"
	"encoding/json"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type ScreeningResultsService interface {
	GetAllScreeningResults(ctx context.Context) ([]models.ScreeningResults, error)
	CreateScreeningResults(ctx context.Context, input models.ScreeningResultsCreate) (models.ScreeningResults, error)
	GetScreeningResults(ctx context.Context, id string) (models.ScreeningResults, error)
}

type ScreeningResultServiceImpl struct {
	store db.Store
}

func NewScreeningResultService(store db.Store) *ScreeningResultServiceImpl {
	return &ScreeningResultServiceImpl{store: store}
}

// Type conversion helpers
func toScreeningResultsDTO(dbScreeningResult repositories.ScreeningResult) (models.ScreeningResults, error) {
	// Convert RawResponse from pgtype.JSONB to ScreenResume
	var rawResponse models.ScreenResume
	// Check if RawResponse has data
	if dbScreeningResult.RawResponse.Status == pgtype.Present {
		// Unmarshal the JSONB data into ScreenResume
		err := json.Unmarshal(dbScreeningResult.RawResponse.Bytes, &rawResponse)
		if err != nil {
			return models.ScreeningResults{}, fmt.Errorf("failed to unmarshal raw response: %w", err)
		}
	}

	return models.ScreeningResults{
		ID:            dbScreeningResult.ID.String(),
		ApplicationID: dbScreeningResult.ApplicationID.String(),
		ModelVersion:  dbScreeningResult.ModelVersion,
		RawResponse:   rawResponse,
		ProcessedAt:   dbScreeningResult.ProcessedAt,
	}, nil
}

func toCreateScreeningResultParams(input models.ScreeningResultsCreate) (repositories.CreateScreeningResultsParams, error) {
	applicationID, err := uuid.Parse(input.ApplicationID)
	if err != nil {
		return repositories.CreateScreeningResultsParams{}, fmt.Errorf("invalid ApplicationID: %w", err)
	}

	// Convert ScreenResume to JSONB
	rawJSON, err := json.Marshal(input.RawResponse)
	if err != nil {
		return repositories.CreateScreeningResultsParams{}, fmt.Errorf("failed to marshal raw response: %w", err)
	}

	// Create JSONB object
	rawResponseJSONB := pgtype.JSONB{
		Bytes:  rawJSON,
		Status: pgtype.Present,
	}

	return repositories.CreateScreeningResultsParams{
		ApplicationID: applicationID,
		ModelVersion:  input.ModelVersion,
		RawResponse:   rawResponseJSONB,
	}, nil
}

// Service method implementations
func (s *ScreeningResultServiceImpl) GetAllScreeningResults(ctx context.Context) ([]models.ScreeningResults, error) {
	dbScreeningResults, err := s.store.ListScreeningResults(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobRoles: %w", err)
	}

	jobRoles := make([]models.ScreeningResults, len(dbScreeningResults))
	for i, f := range dbScreeningResults {
		jobRoles[i], err = toScreeningResultsDTO(f)
		if err != nil {
			return nil, fmt.Errorf("failed to convert screening result: %w", err)
		}
	}
	return jobRoles, nil
}

func (s *ScreeningResultServiceImpl) CreateScreeningResults(ctx context.Context, input models.ScreeningResultsCreate) (models.ScreeningResults, error) {
	params, err := toCreateScreeningResultParams(input)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to create jobRole: %w", err)
	}

	dbScreeningResult, err := s.store.CreateScreeningResults(ctx, params)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to create jobRole: %w", err)
	}

	screeningResultsDTO, err := toScreeningResultsDTO(dbScreeningResult)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to create jobRole: %w", err)
	}
	return screeningResultsDTO, nil
}

func (s *ScreeningResultServiceImpl) CreateScreeningResultsWithTx(ctx context.Context, q repositories.Querier, input models.ScreeningResultsCreate) (models.ScreeningResults, error) {
	params, err := toCreateScreeningResultParams(input)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to create screening result params: %w", err)
	}

	dbScreeningResult, err := q.CreateScreeningResults(ctx, params)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to create screening result: %w", err)
	}

	screeningResultsDTO, err := toScreeningResultsDTO(dbScreeningResult)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to convert to DTO: %w", err)
	}

	return screeningResultsDTO, nil
}

func (s *ScreeningResultServiceImpl) GetScreeningResults(ctx context.Context, id string) (models.ScreeningResults, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbScreeningResult, err := s.store.GetScreeningResults(ctx, uuidID)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to get jobRole: %w", err)
	}

	screeningResultsDTO, err := toScreeningResultsDTO(dbScreeningResult)
	if err != nil {
		return models.ScreeningResults{}, fmt.Errorf("failed to create jobRole: %w", err)
	}
	return screeningResultsDTO, nil
}
