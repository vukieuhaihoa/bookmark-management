package dbutils

import (
	"errors"
	"strings"
)

// errorFilter is a list of functions that filter and categorize database errors.
var errorFilter = []func(err error) (bool, error){
	filterDuplicationError,
	filterRecordNotFoundError,
	filterForeignKeyError,
	filterInvalidSortFieldError,
}

var (
	ErrDuplicationType    = errors.New("duplication type error")
	ErrRecordNotFoundType = errors.New("not found type error")
	ErrForeignKeyType     = errors.New("foreign key constraint error")
	ErrInvalidSortField   = errors.New("invalid sort field error")
)

// CatchDBError inspects the provided error and categorizes it into predefined error types.
// It returns a corresponding error type if a match is found; otherwise, it returns the original error.
//
// Parameters:
//   - err: The error to be inspected.
//
// Returns:
//   - error: A categorized error type or the original error.
func CatchDBError(err error) error {
	if err == nil {
		return nil
	}

	for _, filter := range errorFilter {
		match, filteredErr := filter(err)
		if match {
			return filteredErr
		}
	}

	return err
}

// filterDuplicationError checks if the error is a duplication error.
func filterDuplicationError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint"), ErrDuplicationType
}

// filterRecordNotFoundError checks if the error is a record not found error.
func filterRecordNotFoundError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "record not found"), ErrRecordNotFoundType
}

func filterForeignKeyError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "foreign key constraint"), ErrForeignKeyType
}

func filterInvalidSortFieldError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "no such column"), ErrInvalidSortField
}
