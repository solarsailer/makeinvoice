package parser

import (
	"encoding/csv"
	"os"
)

// -------------------------------------------------------
// CSV.
// -------------------------------------------------------

// ParseCSV returns a 2-dimensional string array from a CSV file.
func ParseCSV(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
