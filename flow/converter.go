package flow

import "github.com/russross/blackfriday"

// ConvertAllMarkdownToHTML returns a map of HTML data from a map of
// markdown data.
func ConvertAllMarkdownToHTML(data map[string]string) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		result[key] = ConvertMarkdownToHTML(value)
	}

	return result
}

// ConvertMarkdownToHTML returns an HTML for a given markdown.
func ConvertMarkdownToHTML(markdown string) string {
	return string(
		blackfriday.MarkdownCommon([]byte(markdown)),
	)
}
