package common

// HandlerError is a utility function to handle errors.
// If the provided error is not nil, it panics with the error.
//
// Parameters:
//   - err: The error to be checked
func HandlerError(err error) {
	if err != nil {
		panic(err)
	}
}
