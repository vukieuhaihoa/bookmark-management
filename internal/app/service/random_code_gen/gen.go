package random_code_gen

// GeneratePassword generates a cryptographically secure random password.
// It creates a password of fixed length (12 characters) using alphanumeric characters
// (a-z, A-Z, 0-9). The function uses crypto/rand for secure random number generation.
//
// Returns:
//   - string: The generated password
//   - error: An error if random number generation fails, nil otherwise
func (s *passwordService) GeneratePassword() (string, error) {
	return s.randomCodeGen.GenerateCode("", defaultPassLength)
}
