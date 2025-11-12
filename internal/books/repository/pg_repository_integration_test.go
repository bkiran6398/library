package repository

import (
	"context"
	"testing"
	"time"

	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/bkiran6398/library/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	defaultBook domain.Book = domain.Book{
		ID:              uuid.New(),
		Title:           "Test Book",
		Author:          "Author",
		ISBN:            "ISBN-123",
		PublishedYear:   nil,
		CopiesTotal:     3,
		CopiesAvailable: 2,
	}
)

func setupTestDB(t *testing.T) (repository *pgRepository, cleanup func()) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("library"),
		postgres.WithUsername("library"),
		postgres.WithPassword("secret"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	pool, err := db.ConnectAndMigrate(ctx, host, port.Int(), "library", "secret", "library", "disable", 5, 1)
	require.NoError(t, err)

	cleanup = func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return NewPgRepository(pool), cleanup
}

func TestPgRepository_CreateGet(t *testing.T) {
	repository, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	createContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	created, err := repository.Create(createContext, defaultBook)
	require.NoError(t, err)
	require.Equal(t, defaultBook.ID, created.ID)

	getContext, getCancel := context.WithTimeout(ctx, 5*time.Second)
	defer getCancel()
	got, err := repository.Get(getContext, defaultBook.ID)
	require.NoError(t, err)
	require.Equal(t, "Test Book", got.Title)
	require.Equal(t, "ISBN-123", got.ISBN)
	require.Equal(t, "Author", got.Author)
	require.Nil(t, got.PublishedYear)
	require.Equal(t, 3, got.CopiesTotal)
	require.Equal(t, 2, got.CopiesAvailable)
}

func TestPgRepository_CreateUpdate(t *testing.T) {
	repository, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	createContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	created, err := repository.Create(createContext, defaultBook)
	require.NoError(t, err)
	require.Equal(t, defaultBook.ID, created.ID)

	updatedBook := defaultBook
	updatedBook.Title = "Updated Title"

	updateContext, getCancel := context.WithTimeout(ctx, 5*time.Second)
	defer getCancel()
	got, err := repository.Update(updateContext, updatedBook)
	require.NoError(t, err)
	require.Equal(t, "Updated Title", got.Title)
	require.Equal(t, "ISBN-123", got.ISBN)
	require.Equal(t, "Author", got.Author)
	require.Nil(t, got.PublishedYear)
	require.Equal(t, 3, got.CopiesTotal)
	require.Equal(t, 2, got.CopiesAvailable)
}

func TestPgRepository_CreateList(t *testing.T) {
	repository, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	createContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := repository.Create(createContext, defaultBook)
	require.NoError(t, err)

	defaultBook2 := defaultBook
	defaultBook2.ID = uuid.New()
	defaultBook2.ISBN = "UniqueISBN"
	_, err = repository.Create(createContext, defaultBook2)
	require.NoError(t, err)

	listContext, getCancel := context.WithTimeout(ctx, 5*time.Second)
	defer getCancel()

	filter := domain.ListFilter{}
	got, err := repository.List(listContext, filter)
	require.NoError(t, err)
	require.Len(t, got, 2)
}
