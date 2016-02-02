package converter

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/russross/blackfriday"
	"github.com/solarsailer/makeinvoice/extensions"
)

const (
	pdfConverter = "wkhtmltopdf"
)

// Markdown is a byte slice.
type Markdown []byte

// HTML is a byte slice.
type HTML []byte

// -------------------------------------------------------
// Conversion.
// -------------------------------------------------------

// ConvertMarkdownToHTML returns an HTML for a given markdown.
func ConvertMarkdownToHTML(markdown Markdown) HTML {
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
func ExportHTML(fromMarkdown bool, data []byte, filename string) error {
	filename = appendExtension(filename, extensions.HTML)

	if fromMarkdown {
		data = ConvertMarkdownToHTML(data)
	}

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
func ExportPDF(fromMarkdown bool, data []byte, filename string) error {
	filename = appendExtension(filename, extensions.PDF)

	if fromMarkdown {
		data = ConvertMarkdownToHTML(data)
	}

	return createPDF(data, filename)
}

func createPDF(html HTML, filename string) error {

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
