package table

import "strings"

// -------------------------------------------------------
// Table formatting.
// -------------------------------------------------------

// Format a two-dimensional array into a markdown table.
//
// Example:
//
//   [["head1", "head2"], ["val1", "val2"]]
//   =>
//   "
//     head1|head2
//     -|-
//     val1|val2
//   "
func Format(table [][]string) string {
	result := ""

	for i, row := range table {
		result += formatRow(row)

		if i == 0 {
			result += formatHeader(row)
		}
	}

	return strings.TrimSpace(result)
}

// Format the header line of a markdown table.
// For each col, put an equal number of "-" as the col length
// and separate each col with a "|".
//
// Example:
//
//   ["aaa", "bb"] => "---|--\n"
func formatHeader(row []string) string {
	result := "\n"

	for i, col := range row {
		if i != 0 {
			result += "|"
		}

		result += strings.Repeat("-", len(col))
	}

	return result
}

// Format a row - separate each value by a "|".
//
// Example:
//
//   ["a", "b"] => "a|b\n"
func formatRow(row []string) string {
	// Pass the line if empty.
	if isEmptyLine(row) {
		return ""
	}

	// Otherwise, create the line col by col.
	result := "\n"

	for i, col := range row {
		if i != 0 {
			result += "|"
		}

		result += col
	}

	return result
}

func isEmptyLine(row []string) bool {
	for _, col := range row {
		if strings.TrimSpace(col) != "" {
			return false
		}
	}

	return true
}
