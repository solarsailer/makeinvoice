package parser

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/solarsailer/makeinvoice/table"
)

const defaultSeparator = ','

// -------------------------------------------------------
// CSV.
// -------------------------------------------------------

// parseRawCSV returns a 2-dimensional string array from a CSV file.
func parseRawCSV(filename string, separator string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("cannot open `" + filename + "`: no such file")
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = convertSeparator(separator)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("cannot read `" + filename + "`: invalid CSV file")
	}

	return records, nil
}

// ParseCSVFiles parse a list of files and turn them into a simple
// markdown table. Reference the table with the filename.
func ParseCSVFiles(files []string, separator string) (map[string]string, error) {
	if len(files) == 0 {
		return nil, errors.New("no file to parse")
	}

	data := map[string]string{}

	for _, filename := range files {
		fileData, err := parseRawCSV(filename, separator)
		if err != nil {
			return nil, err
		}

		data[prepareKey(filename)] = table.Format(fileData)
	}

	return data, nil
}

// -------------------------------------------------------
// Utils.
// -------------------------------------------------------

func convertSeparator(separator string) rune {
	runes := []rune(separator)
	if len(runes) == 0 {
		return defaultSeparator
	}

	return runes[0]
}

// -------------------------------------------------------
// Filename.
// -------------------------------------------------------

func prepareKey(filename string) string {
	result := filename

	result = filepath.Base(result)
	result = removeExtension(result)
	result = toUpperFirst(result)

	return result
}

// toUpperFirst returns the string with its first letter in uppercase.
func toUpperFirst(filename string) string {
	runes := []rune(filename)
	if len(runes) == 0 {
		return ""
	}

	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// removeExtension deletes the extension of a filename.
func removeExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.Replace(filename, ext, "", 1)
}
