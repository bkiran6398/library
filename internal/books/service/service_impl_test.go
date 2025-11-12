package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/bkiran6398/library/internal/books/repository/mocks"
	intErr "github.com/bkiran6398/library/internal/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate_DefaultCopiesAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	createRequest := domain.CreateBookRequest{
		Title:       "T",
		Author:      "A",
		ISBN:        "123",
		CopiesTotal: 5,
	}

	expectedBook := domain.Book{
		ID:              uuid.New(),
		Title:           "T",
		Author:          "A",
		ISBN:            "123",
		CopiesTotal:     5,
		CopiesAvailable: 5,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, book domain.Book) (domain.Book, error) {
			require.Equal(t, "T", book.Title)
			require.Equal(t, "A", book.Author)
			require.Equal(t, "123", book.ISBN)
			require.Equal(t, 5, book.CopiesTotal)
			require.Equal(t, 5, book.CopiesAvailable)
			return expectedBook, nil
		}).
		Times(1)

	got, err := service.Create(context.Background(), createRequest)
	require.NoError(t, err)
	require.Equal(t, 5, got.CopiesAvailable)
	require.Equal(t, expectedBook.ID, got.ID)
}

func TestCreate_ValidateCopiesAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

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

func TestCreate_WithExplicitCopiesAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	avail := 3
	createRequest := domain.CreateBookRequest{
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: &avail,
	}

	expectedBook := domain.Book{
		ID:              uuid.New(),
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, book domain.Book) (domain.Book, error) {
			require.Equal(t, 3, book.CopiesAvailable)
			return expectedBook, nil
		}).
		Times(1)

	got, err := service.Create(context.Background(), createRequest)
	require.NoError(t, err)
	require.Equal(t, 3, got.CopiesAvailable)
}

func TestCreate_ValidationError_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	createRequest := domain.CreateBookRequest{
		Title:       "", // Empty title
		Author:      "Author",
		ISBN:        "ISBN-123",
		CopiesTotal: 5,
	}

	_, err := service.Create(context.Background(), createRequest)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")
}

func TestCreate_ValidationError_EmptyAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	createRequest := domain.CreateBookRequest{
		Title:       "Title",
		Author:      "", // Empty author
		ISBN:        "ISBN-123",
		CopiesTotal: 5,
	}

	_, err := service.Create(context.Background(), createRequest)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")
}

func TestCreate_ValidationError_EmptyISBN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	createRequest := domain.CreateBookRequest{
		Title:       "Title",
		Author:      "Author",
		ISBN:        "", // Empty ISBN
		CopiesTotal: 5,
	}

	_, err := service.Create(context.Background(), createRequest)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")
}

func TestCreate_ValidationError_NegativeCopiesTotal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	createRequest := domain.CreateBookRequest{
		Title:       "Title",
		Author:      "Author",
		ISBN:        "ISBN-123",
		CopiesTotal: -1, // Negative
	}

	_, err := service.Create(context.Background(), createRequest)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")
}

func TestCreate_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	createRequest := domain.CreateBookRequest{
		Title:       "Title",
		Author:      "Author",
		ISBN:        "ISBN-123",
		CopiesTotal: 5,
	}

	repoError := intErr.ErrConflict
	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(domain.Book{}, repoError).
		Times(1)

	_, err := service.Create(context.Background(), createRequest)
	require.Error(t, err)
	require.ErrorIs(t, err, intErr.ErrConflict)
}

func TestGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()
	expectedBook := domain.Book{
		ID:              bookID,
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockRepo.EXPECT().
		Get(gomock.Any(), bookID).
		Return(expectedBook, nil).
		Times(1)

	got, err := service.Get(context.Background(), bookID)
	require.NoError(t, err)
	require.Equal(t, expectedBook.ID, got.ID)
	require.Equal(t, expectedBook.Title, got.Title)
}

func TestGet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()

	mockRepo.EXPECT().
		Get(gomock.Any(), bookID).
		Return(domain.Book{}, intErr.ErrNotFound).
		Times(1)

	_, err := service.Get(context.Background(), bookID)
	require.Error(t, err)
	require.ErrorIs(t, err, intErr.ErrNotFound)
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()
	existingBook := domain.Book{
		ID:              bookID,
		Title:           "Old Title",
		Author:          "Old Author",
		ISBN:            "ISBN-OLD",
		CopiesTotal:     5,
		CopiesAvailable: 3,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	updateRequest := domain.UpdateBookRequest{
		Title:           "New Title",
		Author:          "New Author",
		ISBN:            "ISBN-NEW",
		CopiesTotal:     10,
		CopiesAvailable: 8,
	}

	updatedBook := domain.Book{
		ID:              bookID,
		Title:           "New Title",
		Author:          "New Author",
		ISBN:            "ISBN-NEW",
		CopiesTotal:     10,
		CopiesAvailable: 8,
		CreatedAt:       existingBook.CreatedAt,
		UpdatedAt:       time.Now(),
	}

	mockRepo.EXPECT().
		Get(gomock.Any(), bookID).
		Return(existingBook, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, book domain.Book) (domain.Book, error) {
			require.Equal(t, "New Title", book.Title)
			require.Equal(t, "New Author", book.Author)
			require.Equal(t, "ISBN-NEW", book.ISBN)
			require.Equal(t, 10, book.CopiesTotal)
			require.Equal(t, 8, book.CopiesAvailable)
			return updatedBook, nil
		}).
		Times(1)

	got, err := service.Update(context.Background(), bookID, updateRequest)
	require.NoError(t, err)
	require.Equal(t, "New Title", got.Title)
	require.Equal(t, "New Author", got.Author)
}

func TestUpdate_ValidationError_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()
	updateRequest := domain.UpdateBookRequest{
		Title:           "", // Empty title
		Author:          "Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
	}

	_, err := service.Update(context.Background(), bookID, updateRequest)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")
}

func TestUpdate_CopiesAvailableExceedsTotal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()
	updateRequest := domain.UpdateBookRequest{
		Title:           "Title",
		Author:          "Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 6, // Exceeds total
	}

	_, err := service.Update(context.Background(), bookID, updateRequest)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")
	require.Contains(t, err.Error(), "copies_available must be <= copies_total")
}

func TestUpdate_BookNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()
	updateRequest := domain.UpdateBookRequest{
		Title:           "Title",
		Author:          "Author",
		ISBN:            "ISBN-123",
		CopiesTotal:     5,
		CopiesAvailable: 3,
	}

	mockRepo.EXPECT().
		Get(gomock.Any(), bookID).
		Return(domain.Book{}, intErr.ErrNotFound).
		Times(1)

	_, err := service.Update(context.Background(), bookID, updateRequest)
	require.Error(t, err)
	require.ErrorIs(t, err, intErr.ErrNotFound)
}

func TestUpdate_RepositoryUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()
	existingBook := domain.Book{
		ID:              bookID,
		Title:           "Old Title",
		Author:          "Old Author",
		ISBN:            "ISBN-OLD",
		CopiesTotal:     5,
		CopiesAvailable: 3,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	updateRequest := domain.UpdateBookRequest{
		Title:           "New Title",
		Author:          "New Author",
		ISBN:            "ISBN-NEW",
		CopiesTotal:     10,
		CopiesAvailable: 8,
	}

	mockRepo.EXPECT().
		Get(gomock.Any(), bookID).
		Return(existingBook, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(domain.Book{}, intErr.ErrConflict).
		Times(1)

	_, err := service.Update(context.Background(), bookID, updateRequest)
	require.Error(t, err)
	require.ErrorIs(t, err, intErr.ErrConflict)
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()

	mockRepo.EXPECT().
		Delete(gomock.Any(), bookID).
		Return(nil).
		Times(1)

	err := service.Delete(context.Background(), bookID)
	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	bookID := uuid.New()

	mockRepo.EXPECT().
		Delete(gomock.Any(), bookID).
		Return(intErr.ErrNotFound).
		Times(1)

	err := service.Delete(context.Background(), bookID)
	require.Error(t, err)
	require.ErrorIs(t, err, intErr.ErrNotFound)
}

func TestList_Success_NoFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	expectedBooks := []domain.Book{
		{
			ID:              uuid.New(),
			Title:           "Book 1",
			Author:          "Author 1",
			ISBN:            "ISBN-1",
			CopiesTotal:     5,
			CopiesAvailable: 3,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			Title:           "Book 2",
			Author:          "Author 2",
			ISBN:            "ISBN-2",
			CopiesTotal:     10,
			CopiesAvailable: 8,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	filter := domain.ListFilter{}

	mockRepo.EXPECT().
		List(gomock.Any(), filter).
		Return(expectedBooks, nil).
		Times(1)

	got, err := service.List(context.Background(), filter)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, expectedBooks[0].Title, got[0].Title)
}

func TestList_Success_WithFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	titleFilter := "Gatsby"
	expectedBooks := []domain.Book{
		{
			ID:              uuid.New(),
			Title:           "The Great Gatsby",
			Author:          "F. Scott Fitzgerald",
			ISBN:            "ISBN-1",
			CopiesTotal:     5,
			CopiesAvailable: 3,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	filter := domain.ListFilter{
		Title: &titleFilter,
	}

	mockRepo.EXPECT().
		List(gomock.Any(), filter).
		Return(expectedBooks, nil).
		Times(1)

	got, err := service.List(context.Background(), filter)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "The Great Gatsby", got[0].Title)
}

func TestList_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := NewService(mockRepo)

	filter := domain.ListFilter{}
	repoError := errors.New("database error")

	mockRepo.EXPECT().
		List(gomock.Any(), filter).
		Return(nil, repoError).
		Times(1)

	_, err := service.List(context.Background(), filter)
	require.Error(t, err)
	require.Equal(t, repoError, err)
}
