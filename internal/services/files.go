package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/internal/repositories"
	"screenresume/pkg/db"

	"github.com/google/uuid"
)

type FilesService interface {
	GetAllFiles(ctx context.Context) ([]models.Files, error)
	CreateFiles(ctx context.Context, input models.FilesCreate) (models.Files, error)
	GetFiles(ctx context.Context, id string) (models.Files, error)
	UpdateFiles(ctx context.Context, id string, input models.FilesUpdate) (models.Files, error)
	DeleteFiles(ctx context.Context, id string) (any, error)
}

type FileServiceImpl struct {
	store db.Store
}

func NewFileService(store db.Store) *FileServiceImpl {
	return &FileServiceImpl{store: store}
}

// Type conversion helpers
func toFilesDTO(dbFile repositories.File) models.Files {
	return models.Files{
		ID:       dbFile.ID.String(),
		Path:     dbFile.Path,
		FileType: dbFile.FileType,
		Checksum: dbFile.Checksum,
	}
}

func toFilesCreateParams(input models.FilesCreate) repositories.CreateFileParams {
	return repositories.CreateFileParams{
		Path:     input.Path,
		FileType: input.FileType,
		Checksum: input.Checksum,
	}
}

// Service method implementations
func (s *FileServiceImpl) GetAllFiles(ctx context.Context) ([]models.Files, error) {
	dbFiles, err := s.store.ListFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]models.Files, len(dbFiles))
	for i, f := range dbFiles {
		files[i] = toFilesDTO(f)
	}
	return files, nil
}

func (s *FileServiceImpl) CreateFiles(ctx context.Context, input models.FilesCreate) (models.Files, error) {
	params := toFilesCreateParams(input)

	dbFile, err := s.store.CreateFile(ctx, params)
	if err != nil {
		return models.Files{}, fmt.Errorf("failed to create file: %w", err)
	}

	return toFilesDTO(dbFile), nil
}

func (s *FileServiceImpl) GetFiles(ctx context.Context, id string) (models.Files, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Files{}, fmt.Errorf("invalid ID format: %w", err)
	}

	dbFile, err := s.store.GetFile(ctx, uuidID)
	if err != nil {
		return models.Files{}, fmt.Errorf("failed to get file: %w", err)
	}

	return toFilesDTO(dbFile), nil
}

func (s *FileServiceImpl) UpdateFiles(ctx context.Context, id string, input models.FilesUpdate) (models.Files, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return models.Files{}, fmt.Errorf("invalid ID format: %w", err)
	}

	params := repositories.UpdateFileParams{
		ID:       uuidID,
		Path:     input.Path,
		FileType: input.FileType,
		Checksum: input.Checksum,
	}

	if err := s.store.UpdateFile(ctx, params); err != nil {
		return models.Files{}, fmt.Errorf("failed to update file: %w", err)
	}

	// Fetch updated entity
	return s.GetFiles(ctx, id)
}

func (s *FileServiceImpl) DeleteFiles(ctx context.Context, id string) (any, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	if err := s.store.DeleteFile(ctx, uuidID); err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return nil, nil
}
