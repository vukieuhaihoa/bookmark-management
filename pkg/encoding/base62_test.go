package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoding_Encode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		input     int64
		expected  string
		expectErr error
	}{
		{"zero", 0, "0", nil},
		{"negative number", -1, "", ErrNegativeNumber},
		{"one", 1, "1", nil},
		{"ten", 10, "A", nil},
		{"sixty-one", 61, "z", nil},
		{"sixty-two", 62, "10", nil},
		{"one hundred", 100, "1c", nil},
		{"twelve thousand", 12345, "3D7", nil},
		{"one million", 1000000, "4C92", nil},
		{"max int64", 9223372036854775807, "AzL8n0Y58m7", nil},
		{"large number", 1000000000, "15ftgG", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := StdEncoding.EncodeInt64ToString(tc.input)
			if err != nil {
				assert.Equal(t, tc.expectErr, err)
				assert.Equal(t, tc.expected, result)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEncoding_Decode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		input     string
		expected  int64
		expectErr error
	}{
		{"zero", "0", 0, nil},

		{"one", "1", 1, nil},
		{"ten", "A", 10, nil},
		{"sixty-one", "z", 61, nil},
		{"sixty-two", "10", 62, nil},
		{"one hundred", "1c", 100, nil},
		{"twelve thousand", "3D7", 12345, nil},
		{"one million", "4C92", 1000000, nil},
		{"max int64", "AzL8n0Y58m7", 9223372036854775807, nil},
		{"large number", "15ftgG", 1000000000, nil},
		{"invalid character", "abc!", 0, CorruptInputError(3)},
		{"empty string", "", 0, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := StdEncoding.DecodeStringToInt64(tc.input)
			if err != nil {
				assert.Equal(t, tc.expectErr, err)
				assert.Equal(t, tc.expected, result)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEncoding_RoundTrip(t *testing.T) {
	t.Parallel()

	testCases := []int64{0, 1, 61, 62, 100, 1000, 12345, 1000000, 9223372036854775807}

	for _, num := range testCases {
		t.Run("round trip", func(t *testing.T) {
			t.Parallel()
			encoded, err := StdEncoding.EncodeInt64ToString(num)
			assert.NoError(t, err)
			decoded, err := StdEncoding.DecodeStringToInt64(encoded)
			assert.NoError(t, err)
			assert.Equal(t, num, decoded)
		})
	}
}

func TestEncoding_Decode_InvalidInput(t *testing.T) {
	t.Parallel()

	testCases := []string{"abc!", "test ", "AB-CD"}

	for _, input := range testCases {
		t.Run("invalid input", func(t *testing.T) {
			t.Parallel()
			_, err := StdEncoding.DecodeStringToInt64(input)
			assert.Error(t, err)
		})
	}
}
