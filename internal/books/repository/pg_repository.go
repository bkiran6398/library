package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/bkiran6398/library/internal/books/domain"
	intErr "github.com/bkiran6398/library/internal/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgRepository is the PostgreSQL implementation of Repository.
type pgRepository struct {
	dbPool *pgxpool.Pool
}

// NewRepository creates a new PostgreSQL-based Repository implementation.
// This constructor is the only place where consumers should depend on the concrete type.
func NewPgRepository(dbPool *pgxpool.Pool) *pgRepository {
	return &pgRepository{dbPool: dbPool}
}

func (repository *pgRepository) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	const insertQuery = `
INSERT INTO books (id, title, author, isbn, published_year, copies_total, copies_available, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
RETURNING created_at, updated_at;
`
	row := repository.dbPool.QueryRow(ctx, insertQuery, book.ID, book.Title, book.Author, book.ISBN, book.PublishedYear, book.CopiesTotal, book.CopiesAvailable)
	if err := row.Scan(&book.CreatedAt, &book.UpdatedAt); err != nil {
		if isUniqueViolationError(err) {
			return domain.Book{}, intErr.ErrConflict
		}
		return domain.Book{}, fmt.Errorf("insert book: %w", err)
	}
	return book, nil
}

// isUniqueViolationError checks if the error is a PostgreSQL unique violation error.
func isUniqueViolationError(err error) bool {
	var pgError *pgconn.PgError
	return errors.As(err, &pgError) && pgError.Code == "23505"
}

func (repository *pgRepository) Get(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	const selectQuery = `
SELECT id, title, author, isbn, published_year, copies_total, copies_available, created_at, updated_at
FROM books WHERE id=$1;
`
	var book domain.Book
	row := repository.dbPool.QueryRow(ctx, selectQuery, bookID)
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.CopiesTotal, &book.CopiesAvailable, &book.CreatedAt, &book.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, intErr.ErrNotFound
		}
		return domain.Book{}, fmt.Errorf("get book: %w", err)
	}
	return book, nil
}

func (repository *pgRepository) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	const updateQuery = `
UPDATE books SET title=$2, author=$3, isbn=$4, published_year=$5, copies_total=$6, copies_available=$7, updated_at=NOW()
WHERE id=$1
RETURNING created_at, updated_at;
`
	row := repository.dbPool.QueryRow(ctx, updateQuery, book.ID, book.Title, book.Author, book.ISBN, book.PublishedYear, book.CopiesTotal, book.CopiesAvailable)
	if err := row.Scan(&book.CreatedAt, &book.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, intErr.ErrNotFound
		}
		if isUniqueViolationError(err) {
			return domain.Book{}, intErr.ErrConflict
		}
		return domain.Book{}, fmt.Errorf("update book: %w", err)
	}
	return book, nil
}

func (repository *pgRepository) Delete(ctx context.Context, bookID uuid.UUID) error {
	const deleteQuery = `DELETE FROM books WHERE id=$1;`
	result, err := repository.dbPool.Exec(ctx, deleteQuery, bookID)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	if result.RowsAffected() == 0 {
		return intErr.ErrNotFound
	}
	return nil
}

func (repository *pgRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Book, error) {
	query, queryArguments := buildListQuery(filter)
	rows, err := repository.dbPool.Query(ctx, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	defer rows.Close()

	books, err := scanBooksFromRows(rows)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// scanBooksFromRows scans database rows into Book entities.
func scanBooksFromRows(rows pgx.Rows) ([]domain.Book, error) {
	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublishedYear, &book.CopiesTotal, &book.CopiesAvailable, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan book row: %w", err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return books, nil
}
