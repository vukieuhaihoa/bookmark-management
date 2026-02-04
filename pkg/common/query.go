package common

import (
	"errors"
	"strings"
)

var ErrInvalidSortField = errors.New("invalid sort field")

type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

// SortedField represents a field to sort by along with the direction of sorting.
// Field: The name of the field to sort by.
// Direction: The direction of sorting, either ascending (ASC) or descending (DESC).
type SortedField struct {
	Field     string
	Direction SortDirection
}

// ParseSortParams parses a sort string and maps it to a slice of SortedField structures.
// The sort string can contain multiple fields separated by commas, with an optional '-' prefix
// to indicate descending order.
//
// Parameters:
//   - sort: A comma-separated string of fields to sort by. A '-' prefix indicates descending order.
//   - allowed: A map of allowed field names to their corresponding database column names.
//
// Returns:
//   - []SortedField: A slice of SortedField structures representing the parsed sort parameters.
//   - error: An error if any of the specified fields are not allowed.
//
// Example:
//
//	allowed := map[string]string{
//	    "createdAt": "created_at",
//	    "updatedAt": "updated_at",
//	}
//	sort := "-createdAt,updatedAt"
//	parsed, err := ParseSortParams(sort, allowed)
//	// parsed will contain:
//	// []SortedField{
//	//     {Field: "created_at", Direction: SortDesc},
//	//     {Field: "updated_at", Direction: SortAsc},
//	// }
func ParseSortParams(sort string, allowed map[string]string) ([]SortedField, error) {
	if sort == "" {
		return []SortedField{
			{
				Field: allowed["createdAt"], Direction: SortDesc,
			},
		}, nil
	}

	parts := strings.Split(sort, ",")
	results := make([]SortedField, 0, len(parts))

	for _, part := range parts {
		direction := SortAsc
		field := part

		if strings.HasPrefix(part, "-") {
			direction = SortDesc
			field = strings.TrimPrefix(part, "-")
		}

		mappedField, ok := allowed[field]
		if !ok {
			return nil, ErrInvalidSortField
		}

		results = append(results, SortedField{
			Field:     mappedField,
			Direction: direction,
		})
	}

	return results, nil
}

// Paging represents pagination information for query results.
type Paging struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// Process normalizes the paging parameters to ensure they fall within acceptable ranges.
// It adjusts the Page and Limit fields as follows:
// - If Page is less than 1, it is set to 1.
// - If Limit is less than or equal to 1, it is set to 1.
// - If Limit is greater than or equal to 20, it is set to 20.
func (p *Paging) Process() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit <= 1 {
		p.Limit = 1
	}

	if p.Limit >= 20 {
		p.Limit = 20
	}

}

// QueryOptions encapsulates options for querying data, including sorting and pagination.
type QueryOptions struct {
	Sorting []SortedField
	Paging
}
