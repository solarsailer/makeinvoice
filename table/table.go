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
			result += formatHeader(len(row))
		}
	}

	return strings.TrimSpace(result)
}

// Format the header line of a markdown table.
// For each col, put a "-" and separate each col with a "|".
//
// Example:
//
//   ["a", "b"] => "-|-\n"
func formatHeader(length int) string {
	result := "\n"

	for i := 0; i < length; i++ {
		if i != 0 {
			result += "|"
		}

		result += "-"
	}

	return result
}

// Format a row - separate each value by a "|".
//
// Example:
//
//   ["a", "b"] => "a|b\n"
func formatRow(row []string) string {
	result := "\n"

	for i, col := range row {
		if i != 0 {
			result += "|"
		}

		result += col
	}

	return result
}
