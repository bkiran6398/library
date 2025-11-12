//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks
package service

import (
	"context"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/google/uuid"
)

// Service defines the interface for book business logic operations.
// Consumers should depend on this interface, not on concrete implementations.
type Service interface {
	Create(ctx context.Context, createRequest domain.CreateBookRequest) (domain.Book, error)
	Get(ctx context.Context, bookID uuid.UUID) (domain.Book, error)
	Update(ctx context.Context, bookID uuid.UUID, updateRequest domain.UpdateBookRequest) (domain.Book, error)
	Delete(ctx context.Context, bookID uuid.UUID) error
	List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, error)
}
