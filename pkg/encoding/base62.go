package encoding

import (
	"errors"
	"strconv"
)

var (
	ErrNegativeNumber = errors.New("negative numbers are not supported")
)

// Encoding represents a base62 encoding scheme defined by a 62-character alphabet.
type Encoding struct {
	// The encoding table.
	// Purpose: Map an index(0..61) to the character byte that represents it.
	encode [62]byte

	// The decoding table.
	// Purpose: Reverse lookup - maps a character byte back to its index value (0-61).
	// Why size 256? Because a byte can have values 0-255, so we need space for all possible byte values. Invalid characters map to 255 (or another sentinel value) to indicate "not in alphabet".
	decodeMap [256]byte // 256 to cover all possible byte values
}

// CorruptInputError reports the presence of invalid base62 data at a given offset.
type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base62 data at input byte " + strconv.FormatInt(int64(e), 10)
}

// NewEncoding creates a new Encoding defined by the given alphabet string.
//
// Parameters:
//   - alphabet: A string of exactly 62 unique characters representing the encoding alphabet.
//
// Returns:
//   - An Encoding instance configured with the provided alphabet.
//
// Panics:
//   - If the alphabet length is not 62.
//   - If the alphabet contains duplicate characters.
func NewEncoding(alphabet string) *Encoding {
	if len(alphabet) != 62 {
		panic("encoding alphabet is not 62 bytes long")
	}

	enc := new(Encoding)

	// Initialize decode map with sentinel value 255 (indicating invalid character)
	for i := range enc.decodeMap {
		enc.decodeMap[i] = 255
	}

	// Build encode array and decode map
	seen := make(map[byte]bool)
	for i := 0; i < 62; i++ {
		c := alphabet[i]
		if seen[c] {
			panic("encoding alphabet contains duplicate character: " + string(c))
		}
		seen[c] = true
		enc.encode[i] = c
		enc.decodeMap[c] = byte(i)
	}

	return enc
}

// StdEncoding is the standard base62 encoding using 0-9, A-Z, a-z.
var StdEncoding = NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// EncodeInt64ToString encodes a non-negative int64 number into a base62 string using the encoding's alphabet.
//
// Parameters:
//   - num: The non-negative int64 number to encode.
//
// Returns:
//   - A base62 encoded string representing the input number.
//   - An error if the input number is negative.
func (enc *Encoding) EncodeInt64ToString(num int64) (string, error) {
	// Handle negative numbers
	if num < 0 {
		return "", ErrNegativeNumber
	}

	// Special case for zero
	if num == 0 {
		return string(enc.encode[0]), nil
	}

	var result []byte
	for num > 0 {
		remainder := num % 62
		result = append(result, enc.encode[remainder])
		num = num / 62
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}

// DecodeStringToInt64 decodes a base62 encoded string back to an int64.
//
// Parameters:
//   - s: The base62 encoded string to decode.
//
// Returns:
//   - The decoded int64 value.
//   - An error if the input string contains invalid characters or if the decoded value exceeds int64 limits.
func (enc *Encoding) DecodeStringToInt64(s string) (int64, error) {
	if len(s) == 0 {
		return 0, nil
	}

	var result int64

	for i, c := range []byte(s) {
		val := enc.decodeMap[c]
		if val == 255 {
			return 0, CorruptInputError(i)
		}

		if result > (1<<63-1-int64(val))/62 {
			return 0, CorruptInputError(i) // overflow
		}

		result = result*62 + int64(val)
	}

	return result, nil
}
