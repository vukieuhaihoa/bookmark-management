package common

import (
	"errors"
	"reflect"
	"testing"
)

func Test_ParseSortParams(t *testing.T) {

	allowed := map[string]struct{}{
		"created_at": {},
		"updated_at": {},
	}

	tests := []struct {
		name string

		sort    string
		allowed map[string]struct{}
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
			sort:    "updated_at",
			allowed: allowed,
			want: []SortedField{
				{Field: "updated_at", Direction: SortAsc},
			},
			wantErr: nil,
		},
		{
			name:    "single descending field",
			sort:    "-created_at",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
			},
			wantErr: nil,
		},
		{
			name:    "multiple fields",
			sort:    "-created_at,updated_at",
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
