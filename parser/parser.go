package parser

import (
	"encoding/csv"
	"errors"
	"os"
)

// -------------------------------------------------------
// CSV.
// -------------------------------------------------------

// ParseCSV returns a 2-dimensional string array from a CSV file.
func ParseCSV(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("cannot open `" + filename + "`: no such file")
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, errors.New("invalid CSV file")
	}

	return records, nil
}
