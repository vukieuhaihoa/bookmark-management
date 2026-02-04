package common

import (
	"errors"
	"reflect"
	"testing"
)

func Test_ParseSortParams(t *testing.T) {

	allowed := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}

	tests := []struct {
		name string

		sort    string
		allowed map[string]string
		want    []SortedField
		wantErr error
	}{
		{
			name:    "empty sort string",
			sort:    "",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
			},
			wantErr: nil,
		},
		{
			name:    "single ascending field",
			sort:    "updatedAt",
			allowed: allowed,
			want: []SortedField{
				{Field: "updated_at", Direction: SortAsc},
			},
			wantErr: nil,
		},
		{
			name:    "single descending field",
			sort:    "-createdAt",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
			},
			wantErr: nil,
		},
		{
			name:    "multiple fields",
			sort:    "-createdAt,updatedAt",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
				{Field: "updated_at", Direction: SortAsc},
			},
			wantErr: nil,
		},
		{
			name:    "invalid field",
			sort:    "invalidField",
			allowed: allowed,
			want:    nil,
			wantErr: ErrInvalidSortField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSortParams(tt.sort, tt.allowed)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseSortParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSortParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
