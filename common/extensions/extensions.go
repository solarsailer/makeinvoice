package extensions

import "path/filepath"

const (
	// Markdown extension.
	Markdown = ".md"

	// HTML extension.
	HTML = ".html"

	// PDF extension.
	PDF = ".pdf"

	// CSS extension.
	CSS = ".css"
)

// -------------------------------------------------------
// Functions.
// -------------------------------------------------------

// Force appends the extension if there's none or if it's incorrect.
func Force(filename, forceExtension string) string {
	ext := filepath.Ext(filename)

	if ext == "" || ext != forceExtension {
		return filename + forceExtension
	}

	return filename
}
