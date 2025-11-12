package service

import (
	"fmt"

	"github.com/bkiran6398/library/internal/books/domain"
	intErr "github.com/bkiran6398/library/internal/errors"
	"github.com/go-playground/validator/v10"
)

// validateCreateRequest validates a CreateBookRequest and returns an error if validation fails.
func validateCreateRequest(validatorInstance *validator.Validate, request domain.CreateBookRequest) error {
	if err := validatorInstance.Struct(request); err != nil {
		return fmt.Errorf("%w: %v", intErr.ErrBadRequest, err)
	}
	return nil
}

// validateUpdateRequest validates an UpdateBookRequest and returns an error if validation fails.
func validateUpdateRequest(validatorInstance *validator.Validate, request domain.UpdateBookRequest) error {
	if err := validatorInstance.Struct(request); err != nil {
		return fmt.Errorf("%w: %v", intErr.ErrBadRequest, err)
	}
	return nil
}

// validateCopiesAvailable validates that copies_available does not exceed copies_total.
func validateCopiesAvailable(copiesAvailable, copiesTotal int) error {
	if copiesAvailable > copiesTotal {
		return fmt.Errorf("%w: copies_available must be <= copies_total", intErr.ErrBadRequest)
	}
	return nil
}
