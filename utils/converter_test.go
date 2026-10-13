package converter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/csvtojsonconverter/config"
)

func TestConvertFile(t *testing.T) {
	// Prepare temporary CSV file
	csvContent := `id,name,price,available
1,Apple,0.99,true
2,Banana,0.59,false
invalid,Cherry,1.20,true
3,Date,,true
`
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "input.csv")
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("write csv: %v", err)
	}
	jsonPath := filepath.Join(tmpDir, "output.json")

	// Minimal config
	cfg := &config.Config{
		LogLevel: "error",
		Schema: []config.ColumnSpec{
			{Name: "id", Type: "int", Required: true},
			{Name: "name", Type: "string", Required: true},
			{Name: "price", Type: "float", Required: false},
			{Name: "available", Type: "bool", Required: false},
		},
	}

	// Run conversion
	if err := ConvertFile(csvPath, jsonPath, cfg); err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	// Verify JSON output
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	expected := `[{
  "available": true,
  "id": 1,
  "name": "Apple",
  "price": 0.99
},
{
  "available": false,
  "id": 2,
  "name": "Banana",
  "price": 0.59
},
{
  "available": true,
  "id": 3,
  "name": "Date",
  "price": null
}]`
	if string(data) != expected {
		t.Errorf("unexpected JSON output.\nGot:\n%s\nExpected:\n%s", string(data), expected)
	}
}