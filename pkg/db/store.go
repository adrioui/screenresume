package db

import (
	"context"
	"fmt"
	"screenresume/internal/repositories"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// In db/store.go
type Store interface {
	repositories.Querier
	BeginTx(ctx context.Context) (repositories.Querier, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type SQLStore struct {
	connPool *pgxpool.Pool
	tx       pgx.Tx
	*repositories.Queries
}

func (s *SQLStore) BeginTx(ctx context.Context) (repositories.Querier, error) {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	s.tx = tx
	return repositories.New(tx), nil
}

func (s *SQLStore) CommitTx(ctx context.Context) error {
	if s.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}
	err := s.tx.Commit(ctx)
	s.tx = nil
	return err
}

func (s *SQLStore) RollbackTx(ctx context.Context) error {
	if s.tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}
	err := s.tx.Rollback(ctx)
	s.tx = nil
	return err
}
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  repositories.New(connPool),
		tx:       nil, // Initialize tx as nil
	}
}
