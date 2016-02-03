package converter

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/russross/blackfriday"
	"github.com/solarsailer/makeinvoice/common/entities"
	"github.com/solarsailer/makeinvoice/common/extensions"
)

const (
	pdfConverter = "wkhtmltopdf"
)

// -------------------------------------------------------
// Export.
// -------------------------------------------------------

// Export the data to path.
func Export(data []byte, path string) error {
	if filepath.Ext(path) == extensions.PDF {
		// TODO: Add "--user-style-sheet path/to/css" to wkhtmltopdf command.
		return ExportPDF(data, path)
	}

	if filepath.Ext(path) == extensions.HTML {
		return ExportHTML(data, path)
	}

	return ExportMarkdown(data, path)
}

// -------------------------------------------------------
// Conversion.
// -------------------------------------------------------

// ConvertMarkdownToHTML returns an HTML for a given markdown.
func ConvertMarkdownToHTML(markdown entities.Markdown) entities.HTML {
	return blackfriday.MarkdownCommon(markdown)
}

// -------------------------------------------------------
// Markdown & HTML.
// -------------------------------------------------------

// ExportMarkdown creates a Markdown file.
func ExportMarkdown(data []byte, filename string) error {
	filename = appendExtension(filename, extensions.Markdown)

	return exportFile(data, filename)
}

// ExportHTML creates an HTML file for a given markdown data.
func ExportHTML(data []byte, filename string) error {
	data = ConvertMarkdownToHTML(data)
	filename = appendExtension(filename, extensions.HTML)

	return exportFile(data, filename)
}

func exportFile(data []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.New("cannot create " + filename)
	}
	defer f.Close()

	// Write data to the file.
	_, err = f.Write(data)
	if err != nil {
		return errors.New("cannot write data to " + filename)
	}

	return nil
}

// -------------------------------------------------------
// PDF.
// -------------------------------------------------------

// ExportPDF creates a PDF file for a given markdown data.
func ExportPDF(data []byte, filename string) error {
	data = ConvertMarkdownToHTML(data)
	filename = appendExtension(filename, extensions.PDF)

	return createPDF(data, filename)
}

func createPDF(html entities.HTML, filename string) error {

	// Create a tmp file.
	tmpFile, err := ioutil.TempFile(os.TempDir(), "mkinv_")
	if err != nil {
		return errors.New("cannot create a temporary file")
	}
	defer os.Remove(tmpFile.Name())

	// Fill the temp file.
	_, err = tmpFile.Write(html)
	if err != nil {
		return errors.New("cannot write data to a temporary file")
	}

	// Rename to ".html".
	// `wkhtmltopdf` **needs** an extension to determine the type of the file.
	htmlPath := tmpFile.Name() + extensions.HTML

	if err := os.Rename(tmpFile.Name(), htmlPath); err != nil {
		return errors.New("cannot create a temporary html file")
	}
	defer os.Remove(htmlPath)

	// Check for the converter availability.
	if _, err := exec.LookPath(pdfConverter); err != nil {
		return errors.New("impossible to call `" + pdfConverter + "`: install it and run this command again")
	}

	return exec.Command(pdfConverter, htmlPath, filename).Run()
}

// -------------------------------------------------------
// Utils.
// -------------------------------------------------------

// Add the extension if there's none or if it's incorrect.
func appendExtension(filename, appendExtension string) string {
	ext := filepath.Ext(filename)

	if ext == "" || ext != appendExtension {
		return filename + appendExtension
	}

	return filename
}
