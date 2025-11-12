package repository

import (
	"context"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/google/uuid"
)

// Repository defines the interface for book data access operations.
// Consumers should depend on this interface, not on concrete implementations.
type Repository interface {
	Create(ctx context.Context, book domain.Book) (domain.Book, error)
	Get(ctx context.Context, bookID uuid.UUID) (domain.Book, error)
	Update(ctx context.Context, book domain.Book) (domain.Book, error)
	Delete(ctx context.Context, bookID uuid.UUID) error
	List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, error)
}
