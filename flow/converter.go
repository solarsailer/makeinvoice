package flow

import "github.com/russross/blackfriday"

// ConvertAllMarkdownToHTML returns a map of HTML data from a map of
// markdown data.
func ConvertAllMarkdownToHTML(data map[string]string) map[string]HTML {
	result := make(map[string]HTML)

	for key, value := range data {
		result[key] = ConvertMarkdownToHTML(Markdown(value))
	}

	return result
}

// ConvertMarkdownToHTML returns an HTML for a given markdown.
func ConvertMarkdownToHTML(markdown Markdown) HTML {
	return blackfriday.MarkdownCommon(markdown)
}
