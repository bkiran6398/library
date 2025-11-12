package service

import (
	"github.com/bkiran6398/library/internal/books/domain"
	"github.com/google/uuid"
)

// mapCreateRequestToBook converts a CreateBookRequest to a Book domain entity.
// If CopiesAvailable is not provided, it defaults to CopiesTotal.
func mapCreateRequestToBook(request domain.CreateBookRequest) domain.Book {
	copiesAvailable := request.CopiesTotal
	if request.CopiesAvailable != nil {
		copiesAvailable = *request.CopiesAvailable
	}

	return domain.Book{
		ID:              uuid.New(),
		Title:           request.Title,
		Author:          request.Author,
		ISBN:            request.ISBN,
		PublishedYear:   request.PublishedYear,
		CopiesTotal:     request.CopiesTotal,
		CopiesAvailable: copiesAvailable,
	}
}

// applyUpdateRequestToBook applies update request fields to an existing book.
func applyUpdateRequestToBook(existingBook domain.Book, updateRequest domain.UpdateBookRequest) domain.Book {
	existingBook.Title = updateRequest.Title
	existingBook.Author = updateRequest.Author
	existingBook.ISBN = updateRequest.ISBN
	existingBook.PublishedYear = updateRequest.PublishedYear
	existingBook.CopiesTotal = updateRequest.CopiesTotal
	existingBook.CopiesAvailable = updateRequest.CopiesAvailable
	return existingBook
}
