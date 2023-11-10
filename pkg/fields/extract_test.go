package fields

import (
	"reflect"
	"testing"
)

func TestExtractField(t *testing.T) {
	tests := []struct {
		name   string
		m      any
		fields []any
		want   any
	}{
		{
			"map with slice",
			map[string]any{
				"field": []string{"a", "b"},
			},
			[]any{"field", 1},
			"b",
		},
		{
			"struct with slice",
			struct {
				Field []string
			}{
				Field: []string{"a", "b"},
			},
			[]any{"Field", 0},
			"a",
		},
		{
			"map with nested map",
			map[string]any{
				"field": map[string]any{
					"subfield": "value",
				},
			},
			[]any{"field", "subfield"},
			"value",
		},
		{
			"struct with nested struct",
			struct {
				Field struct {
					Subfield string
				}
			}{
				Field: struct {
					Subfield string
				}{
					Subfield: "value",
				},
			},
			[]any{"Field", "Subfield"},
			"value",
		},
		{
			"map with missing key",
			map[string]any{
				"field": "value",
			},
			[]any{"missing"},
			nil,
		},
		{
			"struct with missing field",
			struct {
				Field string
			}{
				Field: "value",
			},
			[]any{"Missing"},
			nil,
		},
		{
			"slice out of bounds",
			[]string{"a", "b"},
			[]any{2},
			nil,
		},
		{
			"non-map/slice/struct",
			"not a map or slice or struct",
			[]any{"field"},
			nil,
		},
		{
			"map with slice with map",
			map[string]any{
				"field": []map[string]int{
					{"x": 1, "y": 2},
					{"z": 3, "g": 4},
				},
			},
			[]any{"field", "*", "g"},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractdAny(tt.m, tt.fields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractField() = %v, want %v", got, tt.want)
			}
		})
	}
}
