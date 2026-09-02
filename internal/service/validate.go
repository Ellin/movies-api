package service

import (
	"errors"
	"fmt"
	"movies-api/internal/errs"
	"time"

	"github.com/go-playground/validator/v10"
)

// handleValidationError checks all the struct level validation errors from validate.Struct
// and returns a new user-friendly error intended for outputting to the client
func handleValidationError(validationErrors validator.ValidationErrors) error {
	var allErrs []error

	for _, fieldError := range validationErrors {
		field := fieldError.Field() // JSON field name

		switch fieldError.Tag() { // Validation tag name
		case "required":
			allErrs = append(allErrs, fmt.Errorf("%w: %s missing", errs.ErrInvalidUserInput, field))

		case "gte":

			allErrs = append(allErrs, fmt.Errorf("%w: %s must be greater than or equal to %s", errs.ErrInvalidUserInput, field, fieldError.Param()))

		case "lte":
			allErrs = append(allErrs, fmt.Errorf("%w: %s must be less than or equal to %s", errs.ErrInvalidUserInput, field, fieldError.Param()))

		// Assume min & max are only used for strings
		case "min":
			allErrs = append(allErrs, fmt.Errorf("%w: %s must have at least %v characters", errs.ErrInvalidUserInput, field, fieldError.Param()))

		case "max":
			allErrs = append(allErrs, fmt.Errorf("%w: %s cannot exceed %v characters", errs.ErrInvalidUserInput, field, fieldError.Param()))
		}
	}

	return errors.Join(allErrs...)
}

func validateActorMovieDates(birthDate string, releaseYear int) error {
	birthDateParsed, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return fmt.Errorf("invalid birth date: %w", err)
	}

	if birthDateParsed.Year() > releaseYear {
		return fmt.Errorf(
			"actor cannot participate in a movie released in %d before their birth in %d",
			releaseYear,
			birthDateParsed.Year(),
		)
	}

	return nil
}
