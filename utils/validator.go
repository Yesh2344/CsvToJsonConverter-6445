package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/example/csvtojsonconverter/config"
)

// ValidationError describes a single validation problem.
type ValidationError struct {
	Row    int
	Column string
	Reason string
}

// Error implements the error interface.
func (ve ValidationError) Error() string {
	return fmt.Sprintf("row %d, column '%s': %s", ve.Row, ve.Column, ve.Reason)
}

// ValidateRow validates a CSV record according to the schema.
// record is a map of column name -> raw string value.
// Returns nil if the row is valid, otherwise a ValidationError.
func ValidateRow(rowIdx int, record map[string]string, schema []config.ColumnSpec) error {
	for _, col := range schema {
		val, exists := record[col.Name]
		if !exists {
			if col.Required {
				return ValidationError{Row: rowIdx, Column: col.Name, Reason: "missing required column"}
			}
			continue // optional column missing is fine
		}
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			if col.Required {
				return ValidationError{Row: rowIdx, Column: col.Name, Reason: "empty required field"}
			}
			continue
		}
		switch strings.ToLower(col.Type) {
		case "int":
			if _, err := strconv.Atoi(trimmed); err != nil {
				return ValidationError{Row: rowIdx, Column: col.Name, Reason: "invalid integer"}
			}
		case "float":
			if _, err := strconv.ParseFloat(trimmed, 64); err != nil {
				return ValidationError{Row: rowIdx, Column: col.Name, Reason: "invalid float"}
			}
		case "bool":
			if _, err := strconv.ParseBool(trimmed); err != nil {
				return ValidationError{Row: rowIdx, Column: col.Name, Reason: "invalid boolean"}
			}
		case "string":
			// no validation needed
		default:
			return ValidationError{Row: rowIdx, Column: col.Name, Reason: "unknown type in schema"}
		}
	}
	return nil
}