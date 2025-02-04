package services

import (
	"context"
	"database/sql"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
)

type CandidatesService interface {
	GetAllCandidates(ctx context.Context) ([]models.Candidates, error)
	CreateCandidates(ctx context.Context, input models.CandidatesCreate) (models.Candidates, error)
	GetCandidates(ctx context.Context, id string) (models.Candidates, error)
	UpdateCandidates(ctx context.Context, id string, input models.CandidatesUpdate) (models.Candidates, error)
	DeleteCandidates(ctx context.Context, id string) (any, error)
	CandidateAndJobRoles(ctx context.Context, nameSearch string, limit int, page int) ([]models.CandidateAndJobRoles, error)
}

type CandidateServiceImpl struct {
	store db.Store
}

func NewCandidateService(store db.Store) *CandidateServiceImpl {
	return &CandidateServiceImpl{store: store}
}

// Type conversion helpers
func toCandidatesDTO(dbCandidate repositories.Candidate) models.Candidates {
	return models.Candidates{
		ID:       dbCandidate.ID.String(),
		FullName: dbCandidate.FullName,
		Email:    dbCandidate.Email,
		Phone:    dbCandidate.Phone,
		FileID:   dbCandidate.FileID.String(),
		Status:   dbCandidate.Status,
	}
}

func toCandidatesCreateParams(input models.CandidatesCreate) repositories.CreateCandidateParams {
	return repositories.CreateCandidateParams{
		FullName: input.FullName,
		Email:    input.Email,
		Phone:    input.Phone,
		FileID:   uuid.MustParse(input.FileID),
		Status:   input.Status,
	}
}

func toCandidateAndJobRolesDTO(dbCandidate repositories.CandidateAndJobRolesRow) models.CandidateAndJobRoles {
	return models.CandidateAndJobRoles{
		Candidate: toCandidatesDTO(dbCandidate.Candidate),
		JobRole:   toJobRolesDTO(dbCandidate.JobRole),
		AppliedAt: dbCandidate.AppliedAt,
	}
}

// Service method implementations
func (s *CandidateServiceImpl) GetAllCandidates(ctx context.Context) ([]models.Candidates, error) {
	dbCandidates, err := s.store.ListCandidates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list candidates: %w", err)
	}

	candidates := make([]models.Candidates, len(dbCandidates))
	for i, f := range dbCandidates {
		candidates[i] = toCandidatesDTO(f)
	}
	return candidates, nil
}

func (s *CandidateServiceImpl) CreateCandidates(ctx context.Context, input models.CandidatesCreate) (models.Candidates, error) {
	params := toCandidatesCreateParams(input)

	dbCandidate, err := s.store.CreateCandidate(ctx, params)
	if err != nil {
		return models.Candidates{}, fmt.Errorf("failed to create candidate: %w", err)
	}

	return toCandidatesDTO(dbCandidate), nil
}

func (s *CandidateServiceImpl) GetCandidates(ctx context.Context, id string) (models.Candidates, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Candidates{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbCandidate, err := s.store.GetCandidate(ctx, uuidID)
	if err != nil {
		return models.Candidates{}, fmt.Errorf("failed to get candidate: %w", err)
	}

	return toCandidatesDTO(dbCandidate), nil
}

func (s *CandidateServiceImpl) UpdateCandidates(ctx context.Context, id string, input models.CandidatesUpdate) (models.Candidates, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Candidates{}, fmt.Errorf("invalid ID format: %w", err)
	}

	params := repositories.UpdateCandidateParams{
		ID:       uuidID,
		FullName: input.FullName,
		Email:    input.Email,
		Phone:    input.Phone,
		FileID:   uuid.MustParse(input.FileID),
		Status:   input.Status,
	}

	if err := s.store.UpdateCandidate(ctx, params); err != nil {
		return models.Candidates{}, fmt.Errorf("failed to update candidate: %w", err)
	}

	// Fetch updated entity
	return s.GetCandidates(ctx, id)
}

func (s *CandidateServiceImpl) DeleteCandidates(ctx context.Context, id string) (any, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.store.DeleteCandidate(ctx, uuidID); err != nil {
		return nil, fmt.Errorf("failed to delete candidate: %w", err)
	}

	return nil, nil
}

func (s *CandidateServiceImpl) CandidateAndJobRoles(ctx context.Context, nameSearch string, limit int, page int) ([]models.CandidateAndJobRoles, error) {
	dbCandidates, err := s.store.CandidateAndJobRoles(ctx, repositories.CandidateAndJobRolesParams{
		NameSearch: sql.NullString{String: nameSearch, Valid: nameSearch != ""},
		Limitquery: int32(limit),
		Pagequery:  int32(page),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate and job roles: %w", err)
	}

	candidates := make([]models.CandidateAndJobRoles, len(dbCandidates))
	for i, f := range dbCandidates {
		candidates[i] = toCandidateAndJobRolesDTO(f)
	}
	return candidates, nil
}
