package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/example/csvtojsonconverter/config"
	"github.com/example/csvtojsonconverter/logger"
	"github.com/example/csvtojsonconverter/validator"
	"go.uber.org/zap"
)

// ConvertHandler accepts multipart/form-data with a CSV file and returns JSON.
// It re‑uses the same validation logic as the CLI.
func ConvertHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form (max 10MB)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			logger.Logger.Error("parsing multipart form", zap.Error(err))
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			logger.Logger.Error("retrieving file", zap.Error(err))
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		records, err := parseCSV(file, cfg)
		if err != nil {
			logger.Logger.Error("parsing CSV", zap.Error(err))
			http.Error(w, "failed to parse CSV", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(records); err != nil {
			logger.Logger.Error("encoding JSON", zap.Error(err))
		}
	}
}

// parseCSV reads CSV data from an io.Reader, validates, and returns []map[string]interface{}
func parseCSV(r io.Reader, cfg *config.Config) ([]map[string]interface{}, error) {
	csvReader := csv.NewReader(r)
	csvReader.TrimLeadingSpace = true

	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	var out []map[string]interface{}
	rowIdx := 1
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		recordMap := make(map[string]string)
		for name, idx := range headerMap {
			if idx < len(row) {
				recordMap[name] = row[idx]
			}
		}
		if vErr := validator.ValidateRow(rowIdx, recordMap, cfg.Schema); vErr != nil {
			// Skip invalid rows but continue
			logger.Logger.Warn("validation error", zap.Int("row", rowIdx), zap.String("column", vErr.Column), zap.String("reason", vErr.Reason))
			rowIdx++
			continue
		}
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
		out = append(out, converted)
		rowIdx++
	}
	return out, nil
}