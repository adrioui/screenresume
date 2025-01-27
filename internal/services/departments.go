package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
)

type DepartmentsService interface {
	GetAllDepartments(ctx context.Context) ([]models.Departments, error)
	CreateDepartments(ctx context.Context, input models.DepartmentsCreate) (models.Departments, error)
	GetDepartments(ctx context.Context, id string) (models.Departments, error)
	UpdateDepartments(ctx context.Context, id string, input models.DepartmentsUpdate) (models.Departments, error)
	DeleteDepartments(ctx context.Context, id string) (any, error)
}

type DepartmentServiceImpl struct {
	store db.Store
}

func NewDepartmentService(store db.Store) *DepartmentServiceImpl {
	return &DepartmentServiceImpl{store: store}
}

// Type conversion helpers
func toDepartmentsDTO(dbDepartment repositories.Department) models.Departments {
	return models.Departments{
		ID:   dbDepartment.ID.String(),
		Name: dbDepartment.Name,
	}
}

// Service method implementations
func (s *DepartmentServiceImpl) GetAllDepartments(ctx context.Context) ([]models.Departments, error) {
	dbDepartments, err := s.store.ListDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list departments: %w", err)
	}

	departments := make([]models.Departments, len(dbDepartments))
	for i, f := range dbDepartments {
		departments[i] = toDepartmentsDTO(f)
	}
	return departments, nil
}

func (s *DepartmentServiceImpl) CreateDepartments(ctx context.Context, input models.DepartmentsCreate) (models.Departments, error) {
	dbDepartment, err := s.store.CreateDepartment(ctx, input.Name)
	if err != nil {
		return models.Departments{}, fmt.Errorf("failed to create department: %w", err)
	}

	return toDepartmentsDTO(dbDepartment), nil
}

func (s *DepartmentServiceImpl) GetDepartments(ctx context.Context, id string) (models.Departments, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Departments{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbDepartment, err := s.store.GetDepartment(ctx, uuidID)
	if err != nil {
		return models.Departments{}, fmt.Errorf("failed to get department: %w", err)
	}

	return toDepartmentsDTO(dbDepartment), nil
}

func (s *DepartmentServiceImpl) UpdateDepartments(ctx context.Context, id string, input models.DepartmentsUpdate) (models.Departments, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Departments{}, fmt.Errorf("invalid ID format: %w", err)
	}

	params := repositories.UpdateDepartmentParams{
		ID:   uuidID,
		Name: input.Name,
	}

	if err := s.store.UpdateDepartment(ctx, params); err != nil {
		return models.Departments{}, fmt.Errorf("failed to update department: %w", err)
	}

	// Fetch updated entity
	return s.GetDepartments(ctx, id)
}

func (s *DepartmentServiceImpl) DeleteDepartments(ctx context.Context, id string) (any, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.store.DeleteDepartment(ctx, uuidID); err != nil {
		return nil, fmt.Errorf("failed to delete department: %w", err)
	}

	return nil, nil
}
