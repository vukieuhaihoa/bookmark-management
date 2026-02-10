package common

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

var (
	InternalErrorResponse      = Message{Message: "Internal server error"}
	InputErrorResponse         = Message{Message: "Invalid input"}
	InvalidTokenResponse       = Message{Message: "Invalid token"}
	UnauthorizedResponse       = Message{Message: "Unauthorized"}
	InvalidSortedFieldResponse = Message{Message: "Invalid sorted field"}
)

// InputFieldError processes an error and returns a structured Message response.
// If the error is a validation error, it extracts field-specific error messages.
// Otherwise, it returns a generic input error message.
//
// Parameters:
//   - err: The error to be processed
//
// Returns:
//   - Message: A structured message containing error details
func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrorResponse
	}

	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}

	return Message{
		Message: "Invalid input fields",
		Details: errs,
	}
}

type SuccessResponse[Data any] struct {
	Data       Data    `json:"data,omitempty"`
	Message    string  `json:"message,omitempty"`
	Pagination *Paging `json:"pagination,omitempty"`
}
