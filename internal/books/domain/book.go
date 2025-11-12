package domain

import (
	"time"

	"github.com/google/uuid"
)

// Book represents a book entity in the library system.
type Book struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title" validate:"required,min=1"`
	Author          string    `json:"author" validate:"required,min=1"`
	ISBN            string    `json:"isbn" validate:"required"`
	PublishedYear   *int      `json:"published_year,omitempty"`
	CopiesTotal     int       `json:"copies_total" validate:"gte=0"`
	CopiesAvailable int       `json:"copies_available" validate:"gte=0,ltefield=CopiesTotal"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateBookRequest represents the request payload for creating a new book.
type CreateBookRequest struct {
	Title           string `json:"title" validate:"required,min=1"`
	Author          string `json:"author" validate:"required,min=1"`
	ISBN            string `json:"isbn" validate:"required"`
	PublishedYear   *int   `json:"published_year,omitempty"`
	CopiesTotal     int    `json:"copies_total" validate:"gte=0"`
	CopiesAvailable *int   `json:"copies_available,omitempty" validate:"omitempty,gte=0"`
}

// UpdateBookRequest represents the request payload for updating an existing book.
type UpdateBookRequest struct {
	Title           string `json:"title" validate:"required,min=1"`
	Author          string `json:"author" validate:"required,min=1"`
	ISBN            string `json:"isbn" validate:"required"`
	PublishedYear   *int   `json:"published_year,omitempty"`
	CopiesTotal     int    `json:"copies_total" validate:"gte=0"`
	CopiesAvailable int    `json:"copies_available" validate:"gte=0"`
}

// ListFilter represents filtering options for listing books.
type ListFilter struct {
	Title  *string
	Author *string
	ISBN   *string
	Limit  int
	Offset int
}
