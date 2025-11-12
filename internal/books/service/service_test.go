package service

import (
	"context"
	"testing"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	store map[uuid.UUID]domain.Book
}

func (f *fakeRepo) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	f.store[book.ID] = book
	return book, nil
}

func (f *fakeRepo) Get(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	return f.store[bookID], nil
}

func (f *fakeRepo) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	f.store[book.ID] = book
	return book, nil
}

func (f *fakeRepo) Delete(ctx context.Context, bookID uuid.UUID) error {
	delete(f.store, bookID)
	return nil
}

func (f *fakeRepo) List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, error) {
	var out []domain.Book
	for _, v := range f.store {
		out = append(out, v)
	}
	return out, nil
}

func TestCreate_DefaultCopiesAvailable(t *testing.T) {
	service := NewService(&fakeRepo{store: map[uuid.UUID]domain.Book{}})
	createRequest := domain.CreateBookRequest{
		Title:       "T",
		Author:      "A",
		ISBN:        "123",
		CopiesTotal: 5,
	}
	got, err := service.Create(context.Background(), createRequest)
	require.NoError(t, err)
	require.Equal(t, 5, got.CopiesAvailable)
}

func TestCreate_ValidateCopiesAvailable(t *testing.T) {
	service := NewService(&fakeRepo{store: map[uuid.UUID]domain.Book{}})
	avail := 6
	createRequest := domain.CreateBookRequest{
		Title:           "T",
		Author:          "A",
		ISBN:            "123",
		CopiesTotal:     5,
		CopiesAvailable: &avail,
	}
	_, err := service.Create(context.Background(), createRequest)
	require.Error(t, err)
}
