package service

import (
	"context"
	"time"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/bkiran6398/library/internal/books/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	createTimeout = 5 * time.Second
	getTimeout    = 5 * time.Second
	updateTimeout = 5 * time.Second
	deleteTimeout = 5 * time.Second
	listTimeout   = 10 * time.Second
)

// service is the implementation of Service.
type service struct {
	repository repository.Repository
	validator  *validator.Validate
}

// NewService creates a new Service implementation.
// This constructor is the only place where consumers should depend on the concrete type.
func NewService(repository repository.Repository) Service {
	return &service{
		repository: repository,
		validator:  validator.New(),
	}
}

func (serviceInstance *service) Create(ctx context.Context, createRequest domain.CreateBookRequest) (domain.Book, error) {
	if err := validateCreateRequest(serviceInstance.validator, createRequest); err != nil {
		return domain.Book{}, err
	}

	book := mapCreateRequestToBook(createRequest)
	if err := validateCopiesAvailable(book.CopiesAvailable, book.CopiesTotal); err != nil {
		return domain.Book{}, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()
	return serviceInstance.repository.Create(ctxWithTimeout, book)
}

func (serviceInstance *service) Get(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, getTimeout)
	defer cancel()
	return serviceInstance.repository.Get(ctxWithTimeout, bookID)
}

func (serviceInstance *service) Update(ctx context.Context, bookID uuid.UUID, updateRequest domain.UpdateBookRequest) (domain.Book, error) {
	if err := validateUpdateRequest(serviceInstance.validator, updateRequest); err != nil {
		return domain.Book{}, err
	}

	if err := validateCopiesAvailable(updateRequest.CopiesAvailable, updateRequest.CopiesTotal); err != nil {
		return domain.Book{}, err
	}

	existingBook, err := serviceInstance.Get(ctx, bookID)
	if err != nil {
		return domain.Book{}, err
	}

	updatedBook := applyUpdateRequestToBook(existingBook, updateRequest)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()
	return serviceInstance.repository.Update(ctxWithTimeout, updatedBook)
}

func (serviceInstance *service) Delete(ctx context.Context, bookID uuid.UUID) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()
	return serviceInstance.repository.Delete(ctxWithTimeout, bookID)
}

func (serviceInstance *service) List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, listTimeout)
	defer cancel()
	return serviceInstance.repository.List(ctxWithTimeout, filter)
}
