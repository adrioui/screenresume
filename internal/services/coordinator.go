package services

import (
	"context"
	"fmt"
	"screenresume/internal/models"
	"screenresume/pkg/db"
)

type CoordinatorService struct {
	store            db.Store
	screeningService *ScreeningResultServiceImpl
}

func (c *CoordinatorService) CreateScreeningWithRelatedData(ctx context.Context, screeningInput models.ScreeningResultsCreate) error {
	// Begin transaction
	txStore, err := c.store.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback - will be no-op if transaction is committed
	defer c.store.RollbackTx(ctx)

	// Execute operations in different services
	screeningResult, err := c.screeningService.CreateScreeningResultsWithTx(ctx, txStore, screeningInput)
	if err != nil {
		return fmt.Errorf("failed to create screening results: %w", err)
	}
	fmt.Println(screeningResult)

	// If everything succeeded, commit the transaction
	if err := c.store.CommitTx(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
