package converter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/example/csvtojsonconverter/config"
	"github.com/example/csvtojsonconverter/logger"
	"github.com/example/csvtojsonconverter/validator"
)

// ConvertFile reads a CSV file, validates each row, and writes a JSON array to outPath.
// Returns an error if any step fails.
func ConvertFile(csvPath, outPath string, cfg *config.Config) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("opening CSV file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading CSV header: %w", err)
	}
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	// Verify that required columns exist in header
	for _, col := range cfg.Schema {
		if col.Required {
			if _, ok := headerMap[col.Name]; !ok {
				return fmt.Errorf("required column %s missing from CSV header", col.Name)
			}
		}
	}
// small cleanup

	var records []map[string]interface{}
	rowIdx := 1 // CSV rows start after header

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Logger.Error("reading CSV row", zap.Int("row", rowIdx), zap.Error(err))
			return fmt.Errorf("reading CSV row %d: %w", rowIdx, err)
		}

		// Build a map of column name -> raw string
		recordMap := make(map[string]string)
		for name, idx := range headerMap {
			if idx < len(row) {
				recordMap[name] = row[idx]
			}
		}

		// Validate
		if vErr := validator.ValidateRow(rowIdx, recordMap, cfg.Schema); vErr != nil {
			logger.Logger.Warn("validation error", zap.Int("row", rowIdx), zap.String("column", vErr.Column), zap.String("reason", vErr.Reason))
			// Skip invalid rows but continue processing
			rowIdx++
			continue
		}

		// Convert types
		converted := make(map[string]interface{})
		for _, col := range cfg.Schema {
			raw := strings.TrimSpace(recordMap[col.Name])
			if raw == "" {
				converted[col.Name] = nil
				continue
			}
			switch strings.ToLower(col.Type) {
			case "int":
				val, _ := strconv.Atoi(raw)
				converted[col.Name] = val
			case "float":
				val, _ := strconv.ParseFloat(raw, 64)
				converted[col.Name] = val
			case "bool":
				val, _ := strconv.ParseBool(raw)
				converted[col.Name] = val
			case "string":
				converted[col.Name] = raw
			}
		}
		records = append(records, converted)
		rowIdx++
	}

	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling JSON: %w", err)
	}

	// Write output
	if err := os.WriteFile(outPath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("writing JSON output: %w", err)
	}
	logger.Logger.Info("conversion completed", zap.Int("rows_written", len(records)), zap.String("output", outPath))
	return nil
}