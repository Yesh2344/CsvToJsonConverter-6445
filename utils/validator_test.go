package validator

import (
	"testing"

	"github.com/example/csvtojsonconverter/config"
)

func TestValidateRow(t *testing.T) {
	schema := []config.ColumnSpec{
		{Name: "id", Type: "int", Required: true},
		{Name: "name", Type: "string", Required: true},
		{Name: "price", Type: "float", Required: false},
		{Name: "available", Type: "bool", Required: false},
	}

	tests := []struct {
		row   map[string]string
		valid bool
	}{
		{row: map[string]string{"id": "10", "name": "Widget", "price": "12.5", "available": "true"}, valid: true},
		{row: map[string]string{"id": "abc", "name": "Widget"}, valid: false}, // invalid int
		{row: map[string]string{"id": "5", "name": ""}, valid: false},        // empty required string
		{row: map[string]string{"id": "7", "name": "Gadget"}, valid: true},  // optional fields missing
	}

	for i, tt := range tests {
		err := ValidateRow(i+1, tt.row, schema)
		if tt.valid && err != nil {
			t.Errorf("expected valid row %d, got error %v", i+1, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("expected invalid row %d, got nil error", i+1)
		}
	}
}